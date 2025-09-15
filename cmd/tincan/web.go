package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"tincan/pkg/s3client"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start web interface",
	Long:  `Start a web server providing a GUI interface for TinCan operations.`,
	Run:   runWebServer,
}

func runWebServer(cmd *cobra.Command, args []string) {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/download", handleDownload)
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/clean", handleClean)

	fmt.Printf("TinCan web interface starting on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>TinCan - File Transfer</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .section { margin: 20px 0; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
        .file-list { background: #f9f9f9; padding: 10px; margin: 10px 0; }
        button { padding: 10px 20px; margin: 5px; cursor: pointer; }
        input[type="file"] { margin: 10px 0; }
        .error { color: red; }
        .success { color: green; }
    </style>
</head>
<body>
    <h1>TinCan File Transfer</h1>

    <div class="section">
        <h2>Upload File</h2>
        <form id="uploadForm" enctype="multipart/form-data">
            <input type="file" id="fileInput" name="file" required>
            <button type="submit">Upload</button>
        </form>
        <div id="uploadResult"></div>
    </div>

    <div class="section">
        <h2>Files in Bucket</h2>
        <button onclick="listFiles()">Refresh List</button>
        <div id="fileList" class="file-list"></div>
    </div>

    <div class="section">
        <h2>Download File</h2>
        <input type="text" id="downloadKey" placeholder="Enter filename">
        <button onclick="downloadFile()">Download</button>
        <div id="downloadResult"></div>
    </div>

    <div class="section">
        <h2>Clean Up</h2>
        <button onclick="cleanFiles()" style="background-color: #ff4444; color: white;">Delete All Files</button>
        <div id="cleanResult"></div>
    </div>

    <script>
        document.getElementById('uploadForm').onsubmit = function(e) {
            e.preventDefault();
            const formData = new FormData();
            const fileInput = document.getElementById('fileInput');
            formData.append('file', fileInput.files[0]);

            fetch('/upload', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                document.getElementById('uploadResult').innerHTML =
                    data.success ? '<div class="success">' + data.message + '</div>' :
                                  '<div class="error">' + data.error + '</div>';
                if (data.success) listFiles();
            });
        };

        function listFiles() {
            fetch('/list')
            .then(response => response.json())
            .then(data => {
                const fileList = document.getElementById('fileList');
                if (data.success) {
                    if (data.files.length === 0) {
                        fileList.innerHTML = 'No files in bucket';
                    } else {
                        fileList.innerHTML = data.files.map(file =>
                            '<div>' + file + ' <button onclick="downloadFile(\'' + file + '\')">Download</button></div>'
                        ).join('');
                    }
                } else {
                    fileList.innerHTML = '<div class="error">' + data.error + '</div>';
                }
            });
        }

        function downloadFile(filename) {
            const key = filename || document.getElementById('downloadKey').value;
            if (!key) return;

            window.open('/download?key=' + encodeURIComponent(key));
        }

        function cleanFiles() {
            // First get the list of files to show what will be deleted
            fetch('/list')
            .then(response => response.json())
            .then(data => {
                if (!data.success) {
                    document.getElementById('cleanResult').innerHTML = '<div class="error">' + data.error + '</div>';
                    return;
                }

                if (data.files.length === 0) {
                    document.getElementById('cleanResult').innerHTML = '<div class="success">No files to delete</div>';
                    return;
                }

                const fileList = data.files.map(file => '- ' + file).join('\n');
                const confirmMessage = 'The following ' + data.files.length + ' files will be deleted:\n\n' + fileList + '\n\nAre you sure you want to delete these files?';

                if (!confirm(confirmMessage)) return;

                fetch('/clean', { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    document.getElementById('cleanResult').innerHTML =
                        data.success ? '<div class="success">' + data.message + '</div>' :
                                      '<div class="error">' + data.error + '</div>';
                    if (data.success) listFiles();
                });
            });
        }

        // Load file list on page load
        listFiles();
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	template.Must(template.New("home").Parse(tmpl)).Execute(w, nil)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Failed to read file"})
		return
	}
	defer file.Close()

	// Create temporary file
	tempFile, err := os.CreateTemp("", "tincan_upload_*_"+header.Filename)
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Failed to create temp file"})
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy uploaded file to temp file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Failed to save file"})
		return
	}

	// Upload to S3
	client, err := s3client.New()
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "S3 client error: " + err.Error()})
		return
	}

	err = client.Upload(tempFile.Name(), header.Filename)
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Upload failed: " + err.Error()})
		return
	}

	writeJSONResponse(w, map[string]interface{}{"success": true, "message": "File uploaded successfully"})
}

func handleList(w http.ResponseWriter, r *http.Request) {
	client, err := s3client.New()
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "S3 client error: " + err.Error()})
		return
	}

	files, err := client.List()
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Failed to list files: " + err.Error()})
		return
	}

	writeJSONResponse(w, map[string]interface{}{"success": true, "files": files})
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	client, err := s3client.New()
	if err != nil {
		http.Error(w, "S3 client error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create temporary file for download
	tempFile, err := os.CreateTemp("", "tincan_download_*_"+filepath.Base(key))
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	err = client.Download(key, tempFile.Name())
	if err != nil {
		http.Error(w, "Download failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Serve the file
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filepath.Base(key)+"\"")
	http.ServeFile(w, r, tempFile.Name())
}

func handleClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	client, err := s3client.New()
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "S3 client error: " + err.Error()})
		return
	}

	// Get list of files first
	files, err := client.List()
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Failed to list files: " + err.Error()})
		return
	}

	// Delete each file
	for _, file := range files {
		err = client.Delete(file)
		if err != nil {
			writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Failed to delete " + file + ": " + err.Error()})
			return
		}
	}

	writeJSONResponse(w, map[string]interface{}{"success": true, "message": fmt.Sprintf("Deleted %d files", len(files))})
}

func writeJSONResponse(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
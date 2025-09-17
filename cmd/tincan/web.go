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
	http.HandleFunc("/validate", handleValidate)
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/clean", handleClean)
	http.HandleFunc("/delete", handleDelete)

	fmt.Printf("TinCan web interface starting on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>TinCan - File Transfer</title>
    <style>
        * { box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', sans-serif;
            max-width: 900px;
            margin: 0 auto;
            padding: 20px;
            background: #f5f7fa;
            color: #2c3e50;
            line-height: 1.6;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            border-radius: 12px;
            margin-bottom: 30px;
            text-align: center;
            box-shadow: 0 8px 32px rgba(0,0,0,0.1);
        }
        .header h1 { margin: 0 0 10px 0; font-size: 2.5em; font-weight: 300; }
        .version { opacity: 0.9; font-size: 0.9em; margin: 0; }
        .section {
            background: white;
            margin: 20px 0;
            padding: 25px;
            border-radius: 12px;
            box-shadow: 0 4px 16px rgba(0,0,0,0.08);
            border: none;
        }
        .section h2 {
            margin-top: 0;
            color: #2c3e50;
            font-weight: 600;
            border-bottom: 2px solid #ecf0f1;
            padding-bottom: 10px;
        }
        .file-list {
            background: #f8fafc;
            padding: 15px;
            margin: 15px 0;
            border-radius: 8px;
            border: 1px solid #e2e8f0;
            min-height: 60px;
        }
        .file-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px;
            margin: 5px 0;
            background: white;
            border-radius: 6px;
            border: 1px solid #e2e8f0;
            transition: all 0.2s ease;
        }
        .file-item:hover {
            background: #f1f5f9;
            border-color: #cbd5e0;
            transform: translateY(-1px);
        }
        .file-name { font-weight: 500; color: #2d3748; }
        button {
            padding: 10px 20px;
            margin: 5px;
            cursor: pointer;
            border: none;
            border-radius: 6px;
            font-weight: 500;
            transition: all 0.2s ease;
            font-size: 14px;
        }
        .btn-primary {
            background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
            color: white;
        }
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(79, 70, 229, 0.3);
        }
        .btn-secondary {
            background: #6b7280;
            color: white;
        }
        .btn-secondary:hover {
            background: #4b5563;
            transform: translateY(-1px);
        }
        .btn-danger {
            background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
            color: white;
        }
        .btn-danger:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(239, 68, 68, 0.3);
        }
        .btn-download {
            background: #10b981;
            color: white;
            padding: 6px 12px;
            font-size: 12px;
        }
        .btn-download:hover {
            background: #059669;
            transform: translateY(-1px);
        }
        input[type="file"] {
            margin: 10px 0;
            padding: 10px;
            border: 2px dashed #cbd5e0;
            border-radius: 8px;
            background: #f8fafc;
            width: 100%;
            transition: all 0.2s ease;
        }
        input[type="file"]:hover {
            border-color: #4f46e5;
            background: #f0f4ff;
        }
        input[type="text"] {
            padding: 12px;
            border: 2px solid #e2e8f0;
            border-radius: 6px;
            font-size: 14px;
            width: 200px;
            transition: border-color 0.2s ease;
        }
        input[type="text"]:focus {
            outline: none;
            border-color: #4f46e5;
            box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
        }
        .auto-refresh-control {
            display: flex;
            align-items: center;
            gap: 10px;
            margin: 10px 0;
            padding: 12px;
            background: #f8fafc;
            border-radius: 6px;
            border: 1px solid #e2e8f0;
        }
        .auto-refresh-control label {
            display: flex;
            align-items: center;
            gap: 6px;
            font-size: 14px;
            color: #4b5563;
            cursor: pointer;
        }
        input[type="radio"] {
            width: 16px;
            height: 16px;
            cursor: pointer;
        }
        .alert {
            padding: 12px 16px;
            border-radius: 6px;
            margin: 10px 0;
            font-weight: 500;
        }
        .alert-error {
            background: #fef2f2;
            color: #991b1b;
            border: 1px solid #fecaca;
        }
        .alert-success {
            background: #f0fdf4;
            color: #166534;
            border: 1px solid #bbf7d0;
        }
        .loading {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 3px solid #f3f3f3;
            border-top: 3px solid #4f46e5;
            border-radius: 50%;
            animation: spin 1s linear infinite;
            margin-right: 10px;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        .empty-state {
            text-align: center;
            color: #6b7280;
            padding: 40px 20px;
            font-style: italic;
        }
        @media (max-width: 768px) {
            body { padding: 10px; }
            .header { padding: 20px; }
            .header h1 { font-size: 2em; }
            .section { padding: 20px; }
            .file-item { flex-direction: column; align-items: stretch; gap: 10px; }
            input[type="text"] { width: 100%; }
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>TinCan File Transfer</h1>
        <p class="version">Version {{.Version}} ({{.GitCommit}}) - Built {{.BuildDate}}</p>
    </div>

    <div class="section">
        <h2>&#128228; Upload File</h2>
        <form id="uploadForm" enctype="multipart/form-data">
            <input type="file" id="fileInput" name="file" required>
            <button type="submit" class="btn-primary" id="uploadBtn">
                <span id="uploadText">Upload File</span>
            </button>
        </form>
        <div id="uploadResult"></div>
    </div>

    <div class="section">
        <h2>&#128193; Files in Bucket</h2>
        <div style="display: flex; gap: 15px; align-items: center; flex-wrap: wrap;">
            <button onclick="listFiles()" class="btn-secondary" id="refreshBtn">
                <span id="refreshText">Refresh List</span>
            </button>
            <div class="auto-refresh-control">
                <label>
                    <input type="radio" name="autoRefresh" value="off" checked id="autoRefreshOff">
                    Auto-refresh: Off
                </label>
                <label>
                    <input type="radio" name="autoRefresh" value="on" id="autoRefreshOn">
                    <span id="autoRefreshLabel">Auto-refresh: On (30s)</span>
                </label>
            </div>
        </div>
        <div id="fileList" class="file-list"></div>
    </div>

    <div class="section">
        <h2>&#128229; Download File</h2>
        <div style="display: flex; gap: 10px; align-items: center; flex-wrap: wrap;">
            <input type="text" id="downloadKey" placeholder="Enter filename" list="fileNames">
            <datalist id="fileNames"></datalist>
            <button onclick="downloadFile()" class="btn-primary">Download</button>
        </div>
        <div id="downloadResult"></div>
    </div>

    <div class="section">
        <h2>&#128465;&#65039; Clean Up</h2>
        <p style="color: #6b7280; margin-bottom: 15px;">This will delete all files in the bucket. This action cannot be undone.</p>
        <button onclick="cleanFiles()" class="btn-danger" id="cleanBtn">
            <span id="cleanText">Delete All Files</span>
        </button>
        <div id="cleanResult"></div>
    </div>

    <script>
        function showLoading(buttonId, textId, originalText) {
            const button = document.getElementById(buttonId);
            const textSpan = document.getElementById(textId);
            button.disabled = true;
            textSpan.innerHTML = '<span class="loading"></span>' + 'Processing...';
        }

        function hideLoading(buttonId, textId, originalText) {
            const button = document.getElementById(buttonId);
            const textSpan = document.getElementById(textId);
            button.disabled = false;
            textSpan.innerHTML = originalText;
        }

        function showAlert(containerId, message, isSuccess) {
            const container = document.getElementById(containerId);
            container.innerHTML = '<div class="alert ' + (isSuccess ? 'alert-success' : 'alert-error') + '">' + message + '</div>';
            setTimeout(() => {
                if (container.innerHTML.includes(message)) {
                    container.innerHTML = '';
                }
            }, 5000);
        }

        function validateFileName(filename) {
            if (!filename || filename.trim() === '') {
                return { valid: false, error: 'Filename cannot be empty' };
            }

            // Check for invalid characters
            const invalidChars = /[<>:"/\\|?*\x00-\x1f]/;
            if (invalidChars.test(filename)) {
                return { valid: false, error: 'Filename contains invalid characters' };
            }

            // Check filename length
            if (filename.length > 255) {
                return { valid: false, error: 'Filename is too long (max 255 characters)' };
            }

            return { valid: true };
        }

        document.getElementById('uploadForm').onsubmit = function(e) {
            e.preventDefault();
            const fileInput = document.getElementById('fileInput');

            if (!fileInput.files[0]) {
                showAlert('uploadResult', 'Please select a file to upload.', false);
                return;
            }

            const file = fileInput.files[0];

            // Validate file size (50MB limit)
            const maxSize = 50 * 1024 * 1024; // 50MB
            if (file.size > maxSize) {
                showAlert('uploadResult', 'File is too large. Maximum size is 50MB.', false);
                return;
            }

            // Validate filename
            const validation = validateFileName(file.name);
            if (!validation.valid) {
                showAlert('uploadResult', validation.error, false);
                return;
            }

            const formData = new FormData();
            formData.append('file', file);

            showLoading('uploadBtn', 'uploadText', 'Upload File');

            fetch('/upload', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Server error: ' + response.status);
                }
                return response.json();
            })
            .then(data => {
                hideLoading('uploadBtn', 'uploadText', 'Upload File');
                if (data.success) {
                    showAlert('uploadResult', data.message + ' (' + formatFileSize(file.size) + ')', true);
                    fileInput.value = '';
                    listFiles();
                } else {
                    showAlert('uploadResult', data.error, false);
                }
            })
            .catch(error => {
                hideLoading('uploadBtn', 'uploadText', 'Upload File');
                showAlert('uploadResult', 'Upload failed: ' + error.message, false);
            });
        };

        function listFiles() {
            showLoading('refreshBtn', 'refreshText', 'Refresh List');

            fetch('/list')
            .then(response => response.json())
            .then(data => {
                hideLoading('refreshBtn', 'refreshText', 'Refresh List');
                const fileList = document.getElementById('fileList');

                if (data.success) {
                    if (data.files.length === 0) {
                        fileList.innerHTML = '<div class="empty-state">&#128237; No files in bucket<br><small>Upload a file to get started</small></div>';
                    } else {
                        fileList.innerHTML = data.files.map(function(file) {
                            var fileName = file.name || file;
                            var fileSize = formatFileSize(file.size || 0);
                            var uploadDate = 'Unknown';

                            if (file.lastModified) {
                                var date = new Date(file.lastModified);
                                uploadDate = date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
                            }

                            return '<div class="file-item">' +
                                '<div class="file-info">' +
                                    '<div class="file-name">&#128196; ' + fileName + '</div>' +
                                    '<div style="font-size: 0.8em; color: #6b7280;">' +
                                        fileSize + ' &bull; ' + uploadDate +
                                    '</div>' +
                                '</div>' +
                                '<div>' +
                                    '<button onclick="downloadFile(\'' + fileName + '\')" class="btn-download">' +
                                        '&#128229; Download' +
                                    '</button>' +
                                    '<button onclick="deleteFile(\'' + fileName + '\')" class="btn-danger" style="padding: 6px 12px; font-size: 12px; margin-left: 5px;">' +
                                        '&#128465;&#65039; Delete' +
                                    '</button>' +
                                '</div>' +
                            '</div>';
                        }).join('');
                    }

                    // Update autocomplete datalist
                    const datalist = document.getElementById('fileNames');
                    datalist.innerHTML = data.files.map(function(file) {
                        return '<option value="' + (file.name || file) + '">';
                    }).join('');
                } else {
                    fileList.innerHTML = '<div class="alert alert-error">' + data.error + '</div>';
                }
            })
            .catch(error => {
                hideLoading('refreshBtn', 'refreshText', 'Refresh List');
                document.getElementById('fileList').innerHTML = '<div class="alert alert-error">Failed to load files: ' + error.message + '</div>';
            });
        }

        function formatFileSize(bytes) {
            if (bytes === 0) return '0 Bytes';
            const k = 1024;
            const sizes = ['Bytes', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }

        function downloadFile(filename) {
            const key = filename || document.getElementById('downloadKey').value.trim();
            if (!key) {
                showAlert('downloadResult', 'Please enter a filename to download.', false);
                return;
            }

            // Validate filename format
            const validation = validateFileName(key);
            if (!validation.valid) {
                showAlert('downloadResult', validation.error, false);
                return;
            }

            showAlert('downloadResult', 'Validating file...', true);

            // First validate that the file exists
            fetch('/validate?key=' + encodeURIComponent(key))
            .then(response => {
                if (!response.ok) {
                    throw new Error('Server error: ' + response.status);
                }
                return response.json();
            })
            .then(data => {
                if (data.success) {
                    showAlert('downloadResult', 'Starting download...', true);
                    window.open('/download?key=' + encodeURIComponent(key));

                    // Clear the input field and show success message after a delay
                    setTimeout(() => {
                        if (!filename) {
                            document.getElementById('downloadKey').value = '';
                        }
                        showAlert('downloadResult', 'Download started successfully!', true);
                    }, 1000);
                } else {
                    showAlert('downloadResult', data.error, false);
                }
            })
            .catch(error => {
                showAlert('downloadResult', 'Validation failed: ' + error.message, false);
            });
        }

        function deleteFile(filename) {
            if (!confirm('Are you sure you want to delete "' + filename + '"?\n\nThis action cannot be undone.')) {
                return;
            }

            fetch('/delete?key=' + encodeURIComponent(filename), { method: 'DELETE' })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showAlert('fileList', 'File deleted successfully', true);
                    listFiles();
                } else {
                    showAlert('fileList', 'Failed to delete file: ' + data.error, false);
                }
            })
            .catch(error => {
                showAlert('fileList', 'Delete failed: ' + error.message, false);
            });
        }

        function cleanFiles() {
            showLoading('cleanBtn', 'cleanText', 'Delete All Files');

            fetch('/list')
            .then(response => response.json())
            .then(data => {
                hideLoading('cleanBtn', 'cleanText', 'Delete All Files');

                if (!data.success) {
                    showAlert('cleanResult', data.error, false);
                    return;
                }

                if (data.files.length === 0) {
                    showAlert('cleanResult', 'No files to delete', true);
                    return;
                }


				const fileList = data.files.map(file => '- ' + file).join('\n');
                const confirmMessage = 'WARNING: This will permanently delete ALL ' + data.files.length + ' files!\n\nFiles to be deleted:\n' + fileList + '\n\nType "DELETE" to confirm:';

                const userInput = prompt(confirmMessage);
                if (userInput !== 'DELETE') {
                    showAlert('cleanResult', 'Clean operation cancelled.', true);
                    return;
                }

                showLoading('cleanBtn', 'cleanText', 'Delete All Files');

                fetch('/clean', { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    hideLoading('cleanBtn', 'cleanText', 'Delete All Files');
                    showAlert('cleanResult', data.success ? data.message : data.error, data.success);
                    if (data.success) listFiles();
                })
                .catch(error => {
                    hideLoading('cleanBtn', 'cleanText', 'Delete All Files');
                    showAlert('cleanResult', 'Clean failed: ' + error.message, false);
                });
            })
            .catch(error => {
                hideLoading('cleanBtn', 'cleanText', 'Delete All Files');
                showAlert('cleanResult', 'Failed to get file list: ' + error.message, false);
            });
        }

        // Keyboard shortcuts
        document.addEventListener('keydown', function(e) {
            if (e.ctrlKey && e.key === 'u') {
                e.preventDefault();
                document.getElementById('fileInput').click();
            }
            if (e.ctrlKey && e.key === 'r') {
                e.preventDefault();
                listFiles();
            }
            if (e.key === 'Escape') {
                document.getElementById('downloadKey').value = '';
                document.querySelectorAll('.alert').forEach(alert => alert.remove());
            }
        });

        // Auto-refresh control
        let autoRefreshInterval = null;
        let countdownInterval = null;
        let countdownSeconds = 0;

        function updateCountdown() {
            const label = document.getElementById('autoRefreshLabel');
            if (countdownSeconds > 0) {
                label.textContent = 'Auto-refresh: On (' + countdownSeconds + 's)';
                countdownSeconds--;
            } else {
                label.textContent = 'Auto-refresh: On (refreshing...)';
            }
        }

        function toggleAutoRefresh() {
            const isAutoRefreshOn = document.getElementById('autoRefreshOn').checked;

            // Clear existing intervals
            if (autoRefreshInterval) {
                clearInterval(autoRefreshInterval);
                autoRefreshInterval = null;
            }
            if (countdownInterval) {
                clearInterval(countdownInterval);
                countdownInterval = null;
            }

            // Set new intervals if auto-refresh is on
            if (isAutoRefreshOn) {
                // Reset countdown
                countdownSeconds = 30;
                updateCountdown();

                // Start countdown
                countdownInterval = setInterval(updateCountdown, 1000);

                // Set refresh interval
                autoRefreshInterval = setInterval(function() {
                    listFiles();
                    // Reset countdown after refresh
                    countdownSeconds = 30;
                }, 30000);
            } else {
                // Reset label when turned off
                document.getElementById('autoRefreshLabel').textContent = 'Auto-refresh: On (30s)';
            }
        }

        // Add event listeners to radio buttons
        document.getElementById('autoRefreshOff').addEventListener('change', toggleAutoRefresh);
        document.getElementById('autoRefreshOn').addEventListener('change', toggleAutoRefresh);

        // Load file list on page load
        listFiles();

        // Add drag and drop support
        const fileInput = document.getElementById('fileInput');
        const uploadSection = document.querySelector('.section');

        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            uploadSection.addEventListener(eventName, preventDefaults, false);
        });

        function preventDefaults(e) {
            e.preventDefault();
            e.stopPropagation();
        }

        ['dragenter', 'dragover'].forEach(eventName => {
            uploadSection.addEventListener(eventName, highlight, false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            uploadSection.addEventListener(eventName, unhighlight, false);
        });

        function highlight(e) {
            uploadSection.style.background = '#f0f4ff';
            uploadSection.style.borderColor = '#4f46e5';
        }

        function unhighlight(e) {
            uploadSection.style.background = '';
            uploadSection.style.borderColor = '';
        }

        uploadSection.addEventListener('drop', handleDrop, false);

        function handleDrop(e) {
            const dt = e.dataTransfer;
            const files = dt.files;

            if (files.length > 0) {
                fileInput.files = files;
                document.getElementById('uploadForm').dispatchEvent(new Event('submit'));
            }
        }
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")

	data := struct {
		Version   string
		GitCommit string
		BuildDate string
	}{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
	}

	template.Must(template.New("home").Parse(tmpl)).Execute(w, data)
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
	fileNames, err := client.ListNames()
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Failed to list files: " + err.Error()})
		return
	}

	// Delete each file
	for _, fileName := range fileNames {
		err = client.Delete(fileName)
		if err != nil {
			writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Failed to delete " + fileName + ": " + err.Error()})
			return
		}
	}

	writeJSONResponse(w, map[string]interface{}{"success": true, "message": fmt.Sprintf("Deleted %d files", len(fileNames))})
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Missing key parameter"})
		return
	}

	client, err := s3client.New()
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "S3 client error: " + err.Error()})
		return
	}

	// Get list of files to check if the file exists
	fileNames, err := client.ListNames()
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Failed to validate file: " + err.Error()})
		return
	}

	// Check if the file exists in the list
	fileExists := false
	for _, fileName := range fileNames {
		if fileName == key {
			fileExists = true
			break
		}
	}

	if !fileExists {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "File '" + key + "' not found in bucket"})
		return
	}

	writeJSONResponse(w, map[string]interface{}{"success": true, "message": "File exists and is ready for download"})
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Missing key parameter"})
		return
	}

	client, err := s3client.New()
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "S3 client error: " + err.Error()})
		return
	}

	err = client.Delete(key)
	if err != nil {
		writeJSONResponse(w, map[string]interface{}{"success": false, "error": "Failed to delete file: " + err.Error()})
		return
	}

	writeJSONResponse(w, map[string]interface{}{"success": true, "message": "File deleted successfully"})
}

func writeJSONResponse(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
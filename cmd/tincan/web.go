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
        :root {
            /* Light theme colors */
            --bg-primary: #f5f7fa;
            --bg-secondary: white;
            --bg-tertiary: #f8fafc;
            --bg-accent: #f0f4ff;
            --text-primary: #2c3e50;
            --text-secondary: #6b7280;
            --text-tertiary: #4b5563;
            --border-primary: #e2e8f0;
            --border-secondary: #cbd5e0;
            --border-accent: #4f46e5;
            --header-gradient: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            --button-primary: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
            --button-danger: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
            --button-success: #10b981;
            --progress-gradient: linear-gradient(135deg, #10b981 0%, #059669 100%);
            --shadow-primary: rgba(0,0,0,0.08);
            --shadow-secondary: rgba(0,0,0,0.1);
            --alert-success-bg: #f0fdf4;
            --alert-success-text: #166534;
            --alert-success-border: #bbf7d0;
            --alert-error-bg: #fef2f2;
            --alert-error-text: #991b1b;
            --alert-error-border: #fecaca;
        }

        [data-theme="dark"] {
            /* Dark theme colors */
            --bg-primary: #0f172a;
            --bg-secondary: #1e293b;
            --bg-tertiary: #334155;
            --bg-accent: #1e293b;
            --text-primary: #f1f5f9;
            --text-secondary: #94a3b8;
            --text-tertiary: #cbd5e1;
            --border-primary: #334155;
            --border-secondary: #475569;
            --border-accent: #6366f1;
            --header-gradient: linear-gradient(135deg, #4338ca 0%, #5b21b6 100%);
            --button-primary: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
            --button-danger: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%);
            --button-success: #059669;
            --progress-gradient: linear-gradient(135deg, #059669 0%, #047857 100%);
            --shadow-primary: rgba(0,0,0,0.3);
            --shadow-secondary: rgba(0,0,0,0.4);
            --alert-success-bg: #064e3b;
            --alert-success-text: #6ee7b7;
            --alert-success-border: #047857;
            --alert-error-bg: #7f1d1d;
            --alert-error-text: #fca5a5;
            --alert-error-border: #dc2626;
        }

        * { box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', sans-serif;
            max-width: 900px;
            margin: 0 auto;
            padding: 20px;
            background: var(--bg-primary);
            color: var(--text-primary);
            line-height: 1.6;
            transition: background-color 0.3s ease, color 0.3s ease;
        }
        .header {
            background: var(--header-gradient);
            color: white;
            padding: 30px;
            border-radius: 12px;
            margin-bottom: 30px;
            text-align: center;
            box-shadow: 0 8px 32px var(--shadow-secondary);
        }
        .header h1 { margin: 0 0 5px 0; font-size: 3.5em; font-weight: 700; }
        .subtitle { margin: 0 0 10px 0; font-size: 1.3em; font-weight: 400; opacity: 0.95; }
        .version { opacity: 0.9; font-size: 0.9em; margin: 0; }
        .section {
            background: var(--bg-secondary);
            margin: 20px 0;
            padding: 25px;
            border-radius: 12px;
            box-shadow: 0 4px 16px var(--shadow-primary);
            border: none;
            transition: background-color 0.3s ease;
        }
        .section h2 {
            margin-top: 0;
            color: var(--text-primary);
            font-weight: 600;
            border-bottom: 2px solid var(--border-primary);
            padding-bottom: 10px;
        }
        .file-list {
            background: var(--bg-tertiary);
            padding: 15px;
            margin: 15px 0;
            border-radius: 8px;
            border: 1px solid var(--border-primary);
            min-height: 60px;
            transition: background-color 0.3s ease;
        }
        .file-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px;
            margin: 5px 0;
            background: var(--bg-secondary);
            border-radius: 6px;
            border: 1px solid var(--border-primary);
            transition: all 0.2s ease;
        }
        .file-item:hover {
            background: var(--bg-accent);
            border-color: var(--border-secondary);
            transform: translateY(-1px);
        }
        .file-name { font-weight: 500; color: var(--text-primary); }
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
            background: var(--button-primary);
            color: white;
        }
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(79, 70, 229, 0.3);
        }
        .btn-secondary {
            background: var(--text-secondary);
            color: white;
        }
        .btn-secondary:hover {
            background: var(--text-tertiary);
            transform: translateY(-1px);
        }
        .btn-danger {
            background: var(--button-danger);
            color: white;
        }
        .btn-danger:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(239, 68, 68, 0.3);
        }
        .btn-download {
            background: var(--button-success);
            color: white;
            padding: 6px 12px;
            font-size: 12px;
        }
        .btn-download:hover {
            background: #047857;
            transform: translateY(-1px);
        }
        input[type="file"] {
            margin: 10px 0;
            padding: 10px;
            border: 2px dashed var(--border-secondary);
            border-radius: 8px;
            background: var(--bg-tertiary);
            color: var(--text-primary);
            width: 100%;
            transition: all 0.2s ease;
        }
        input[type="file"]:hover {
            border-color: var(--border-accent);
            background: var(--bg-accent);
        }
        input[type="text"] {
            padding: 12px;
            border: 2px solid var(--border-primary);
            border-radius: 6px;
            font-size: 14px;
            width: 200px;
            background: var(--bg-secondary);
            color: var(--text-primary);
            transition: all 0.2s ease;
        }
        input[type="text"]:focus {
            outline: none;
            border-color: var(--border-accent);
            box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
        }
        .auto-refresh-control {
            display: flex;
            align-items: center;
            gap: 10px;
            margin: 10px 0;
            padding: 12px;
            background: var(--bg-tertiary);
            border-radius: 6px;
            border: 1px solid var(--border-primary);
            transition: background-color 0.3s ease;
        }
        .auto-refresh-control label {
            display: flex;
            align-items: center;
            gap: 6px;
            font-size: 14px;
            color: var(--text-secondary);
            cursor: pointer;
        }
        input[type="radio"] {
            width: 16px;
            height: 16px;
            cursor: pointer;
        }
        .progress-container {
            margin: 15px 0;
            padding: 0;
        }
        .progress-bar {
            width: 100%;
            height: 24px;
            background: var(--bg-tertiary);
            border-radius: 12px;
            overflow: hidden;
            border: 1px solid var(--border-primary);
            position: relative;
        }
        .progress-fill {
            height: 100%;
            background: var(--progress-gradient);
            border-radius: 12px;
            transition: width 0.3s ease;
            position: relative;
            width: 0%;
        }
        .progress-text {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            font-size: 12px;
            font-weight: 600;
            color: var(--text-primary);
            z-index: 10;
        }
        .progress-info {
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-size: 13px;
            color: var(--text-secondary);
            margin-top: 5px;
        }
        .hidden {
            display: none;
        }
        .alert {
            padding: 12px 16px;
            border-radius: 6px;
            margin: 10px 0;
            font-weight: 500;
        }
        .alert-error {
            background: var(--alert-error-bg);
            color: var(--alert-error-text);
            border: 1px solid var(--alert-error-border);
        }
        .alert-success {
            background: var(--alert-success-bg);
            color: var(--alert-success-text);
            border: 1px solid var(--alert-success-border);
        }
        .loading {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 3px solid var(--border-primary);
            border-top: 3px solid var(--border-accent);
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
            color: var(--text-secondary);
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
<body data-theme="dark">
    <div class="header">
        <h1>TINCAN</h1>
        <p class="subtitle">File Transfer</p>
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
        <div id="uploadProgress" class="progress-container hidden">
            <div class="progress-bar">
                <div id="uploadProgressFill" class="progress-fill"></div>
                <div id="uploadProgressText" class="progress-text">0%</div>
            </div>
            <div id="uploadProgressInfo" class="progress-info">
                <span id="uploadFileName"></span>
                <span id="uploadFileSize"></span>
            </div>
        </div>
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
        <div id="downloadProgress" class="progress-container hidden">
            <div class="progress-bar">
                <div id="downloadProgressFill" class="progress-fill"></div>
                <div id="downloadProgressText" class="progress-text">0%</div>
            </div>
            <div id="downloadProgressInfo" class="progress-info">
                <span id="downloadFileName"></span>
                <span id="downloadFileSize"></span>
            </div>
        </div>
        <div id="downloadResult"></div>
    </div>

    <div class="section">
        <h2>&#128465;&#65039; Clean Up</h2>
        <p style="color: #6b7280; margin-bottom: 15px;">This will delete all files in the bucket. This action cannot be undone.</p>
        <button onclick="cleanFiles()" class="btn-danger" id="cleanBtn">
            <span id="cleanText">Delete All Files</span>
        </button>
        <div id="cleanProgress" class="progress-container hidden">
            <div class="progress-bar">
                <div id="cleanProgressFill" class="progress-fill"></div>
                <div id="cleanProgressText" class="progress-text">0%</div>
            </div>
            <div id="cleanProgressInfo" class="progress-info">
                <span id="cleanFileName">Preparing...</span>
                <span id="cleanFileSize"></span>
            </div>
        </div>
        <div id="cleanResult"></div>
    </div>

    <div class="section">
        <h2>&#127912; Theme</h2>
        <div class="auto-refresh-control">
            <label>
                <input type="radio" name="theme" value="light" id="themeLight">
                Light Mode
            </label>
            <label>
                <input type="radio" name="theme" value="dark" id="themeDark" checked>
                Dark Mode
            </label>
        </div>
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

        function showProgress(containerId, fileName, fileSize) {
            const container = document.getElementById(containerId);
            const fileNameEl = document.getElementById(containerId.replace('Progress', 'FileName'));
            const fileSizeEl = document.getElementById(containerId.replace('Progress', 'FileSize'));

            container.classList.remove('hidden');
            if (fileNameEl) fileNameEl.textContent = fileName;
            if (fileSizeEl) fileSizeEl.textContent = formatFileSize(fileSize);

            updateProgress(containerId, 0);
        }

        function updateProgress(containerId, percentage) {
            const fillEl = document.getElementById(containerId.replace('Progress', 'ProgressFill'));
            const textEl = document.getElementById(containerId.replace('Progress', 'ProgressText'));

            if (fillEl) fillEl.style.width = percentage + '%';
            if (textEl) textEl.textContent = Math.round(percentage) + '%';
        }

        function hideProgress(containerId) {
            const container = document.getElementById(containerId);
            container.classList.add('hidden');
            updateProgress(containerId, 0);
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
            showProgress('uploadProgress', file.name, file.size);

            // Use XMLHttpRequest for progress tracking
            const xhr = new XMLHttpRequest();

            // Track upload progress
            xhr.upload.addEventListener('progress', function(e) {
                if (e.lengthComputable) {
                    const percentage = (e.loaded / e.total) * 100;
                    updateProgress('uploadProgress', percentage);
                }
            });

            // Handle completion
            xhr.addEventListener('load', function() {
                hideLoading('uploadBtn', 'uploadText', 'Upload File');
                hideProgress('uploadProgress');

                if (xhr.status === 200) {
                    try {
                        const data = JSON.parse(xhr.responseText);
                        if (data.success) {
                            showAlert('uploadResult', data.message + ' (' + formatFileSize(file.size) + ')', true);
                            fileInput.value = '';
                            listFiles();
                        } else {
                            showAlert('uploadResult', data.error, false);
                        }
                    } catch (e) {
                        showAlert('uploadResult', 'Upload completed but response was invalid', false);
                    }
                } else {
                    showAlert('uploadResult', 'Upload failed: Server error ' + xhr.status, false);
                }
            });

            // Handle errors
            xhr.addEventListener('error', function() {
                hideLoading('uploadBtn', 'uploadText', 'Upload File');
                hideProgress('uploadProgress');
                showAlert('uploadResult', 'Upload failed: Network error', false);
            });

            // Handle abort
            xhr.addEventListener('abort', function() {
                hideLoading('uploadBtn', 'uploadText', 'Upload File');
                hideProgress('uploadProgress');
                showAlert('uploadResult', 'Upload cancelled', false);
            });

            // Start upload
            xhr.open('POST', '/upload');
            xhr.send(formData);
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

            // Show progress for validation phase
            showProgress('downloadProgress', key, 0);
            updateProgress('downloadProgress', 20);

            // First validate that the file exists
            fetch('/validate?key=' + encodeURIComponent(key))
            .then(response => {
                updateProgress('downloadProgress', 50);
                if (!response.ok) {
                    throw new Error('Server error: ' + response.status);
                }
                return response.json();
            })
            .then(data => {
                if (data.success) {
                    updateProgress('downloadProgress', 80);

                    // Get file size for display
                    return fetch('/list').then(response => response.json()).then(listData => {
                        if (listData.success) {
                            const fileInfo = listData.files.find(file => (file.name || file) === key);
                            const fileSize = fileInfo ? fileInfo.size : 0;

                            // Update progress with file size
                            document.getElementById('downloadFileSize').textContent = formatFileSize(fileSize);
                            updateProgress('downloadProgress', 100);

                            // Start download
                            setTimeout(() => {
                                hideProgress('downloadProgress');
                                showAlert('downloadResult', 'Download started! Check your browser downloads.', true);
                                window.open('/download?key=' + encodeURIComponent(key));

                                // Clear the input field
                                if (!filename) {
                                    document.getElementById('downloadKey').value = '';
                                }
                            }, 500);
                        } else {
                            throw new Error('Could not get file information');
                        }
                    });
                } else {
                    hideProgress('downloadProgress');
                    showAlert('downloadResult', data.error, false);
                }
            })
            .catch(error => {
                hideProgress('downloadProgress');
                showAlert('downloadResult', 'Download failed: ' + error.message, false);
            });
        }

        function deleteFile(filename) {
            if (!confirm('Are you sure you want to delete "' + filename + '"?\n\nThis action cannot be undone.')) {
                return;
            }

            // Show mini progress in the file list
            showAlert('fileList', 'Deleting ' + filename + '...', true);

            fetch('/delete?key=' + encodeURIComponent(filename), { method: 'DELETE' })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showAlert('fileList', 'File "' + filename + '" deleted successfully', true);
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
                showProgress('cleanProgress', 'Deleting ' + data.files.length + ' files', 0);

                // Simulate progress during deletion
                let progressStep = 0;
                const totalSteps = 10;
                const progressInterval = setInterval(() => {
                    progressStep++;
                    const percentage = (progressStep / totalSteps) * 90; // Leave 10% for completion
                    updateProgress('cleanProgress', percentage);

                    if (progressStep < totalSteps) {
                        document.getElementById('cleanFileName').textContent = 'Deleting files... (' + Math.round(percentage) + '%)';
                    }
                }, 200);

                fetch('/clean', { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    clearInterval(progressInterval);
                    updateProgress('cleanProgress', 100);
                    document.getElementById('cleanFileName').textContent = 'Cleanup completed!';

                    setTimeout(() => {
                        hideLoading('cleanBtn', 'cleanText', 'Delete All Files');
                        hideProgress('cleanProgress');
                        showAlert('cleanResult', data.success ? data.message : data.error, data.success);
                        if (data.success) listFiles();
                    }, 800);
                })
                .catch(error => {
                    clearInterval(progressInterval);
                    hideLoading('cleanBtn', 'cleanText', 'Delete All Files');
                    hideProgress('cleanProgress');
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

        // Theme switching functionality
        function setTheme(theme) {
            document.body.setAttribute('data-theme', theme);
            localStorage.setItem('theme', theme);

            // Update radio buttons
            document.getElementById('themeLight').checked = theme === 'light';
            document.getElementById('themeDark').checked = theme === 'dark';
        }

        function toggleTheme() {
            const selectedTheme = document.querySelector('input[name="theme"]:checked').value;
            setTheme(selectedTheme);
        }

        // Add event listeners to theme radio buttons
        document.getElementById('themeLight').addEventListener('change', toggleTheme);
        document.getElementById('themeDark').addEventListener('change', toggleTheme);

        // Initialize theme from localStorage or default to dark
        document.addEventListener('DOMContentLoaded', function() {
            const savedTheme = localStorage.getItem('theme') || 'dark';
            setTheme(savedTheme);
        });
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
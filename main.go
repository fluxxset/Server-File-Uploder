package main

import (
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

// HTML templates
var tmpl = template.Must(template.New("upload").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>File Upload</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
        }
        .container {
            max-width: 500px;
            margin: auto;
        }
        h1 {
            text-align: center;
        }
        form {
            text-align: center;
            margin-top: 20px;
        }
        input[type="file"] {
            margin-bottom: 20px;
        }
        .progress {
            width: 100%;
            background-color: #f3f3f3;
            margin-top: 20px;
        }
        .progress-bar {
            width: 0%;
            height: 20px;
            background-color: #4CAF50;
        }
    </style>
    <script>
        function uploadFile() {
            const fileInput = document.querySelector('input[type="file"]');
            const formData = new FormData();
            formData.append("file", fileInput.files[0]);

            const xhr = new XMLHttpRequest();
            xhr.open("POST", "/upload", true);

            xhr.upload.onprogress = function(event) {
                const progressBar = document.getElementById("progress-bar");
                const percent = (event.loaded / event.total) * 100;
                progressBar.style.width = percent + "%";
            };

            xhr.onload = function() {
                if (xhr.status == 200) {
                    window.location.href = "/success";
                } else {
                    alert("Error uploading file!");
                }
            };

            xhr.send(formData);
        }
    </script>
</head>
<body>
    <div class="container">
        <h1>Upload File</h1>
        <form onsubmit="event.preventDefault(); uploadFile();">
            <input type="file" name="file" required><br>
            <input type="submit" value="Upload">
        </form>
        <div class="progress">
            <div id="progress-bar" class="progress-bar"></div>
        </div>
    </div>
</body>
</html>
`))

// Success page template
var successTmpl = template.Must(template.New("success").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Upload Success</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
        }
        .container {
            max-width: 500px;
            margin: auto;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>File Uploaded Successfully!</h1>
        <p>Your file has been uploaded.</p>
        <a href="/">Upload another file</a>
    </div>
</body>
</html>
`))

func GetServerIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %v", err)
	}

	for _, iface := range interfaces {
		// Skip down or loopback interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", fmt.Errorf("failed to get addresses for interface %v: %v", iface.Name, err)
		}

		for _, addr := range addrs {
			// Get IP address, ignore non-IP addresses
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Skip IPv6 addresses and loopback addresses
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("no active network interface found")
}

func main() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}
	uploadDir := filepath.Join(cwd, "web-uploads")

	// Ensure the upload directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		fmt.Println("Error creating upload directory:", err)
		return
	}

	// HTTP handlers
	http.HandleFunc("/", uploadForm)
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		uploadFile(w, r, uploadDir)
	})
	http.HandleFunc("/success", successPage)
	ip, err := GetServerIP()

	// Start server
	fmt.Printf("Server started on http://%s:8080", ip)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// Serve the HTML form
func uploadForm(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handle file upload
func uploadFile(w http.ResponseWriter, r *http.Request, uploadDir string) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to upload file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create destination file in the upload directory
	destPath := filepath.Join(uploadDir, header.Filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	// Copy uploaded file content to destination file
	if _, err := io.Copy(destFile, file); err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	// Redirect to success page
	http.Redirect(w, r, "/success", http.StatusSeeOther)
}

// Serve the success page after upload
func successPage(w http.ResponseWriter, r *http.Request) {
	if err := successTmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

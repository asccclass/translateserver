package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"translateserver/db"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit upload size (e.g., 500MB)
	r.Body = http.MaxBytesReader(w, r.Body, 500<<20)

	if err := r.ParseMultipartForm(500 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{
		".m4a": true, ".wav": true, ".mp3": true, ".mp4": true, ".webm": true,
	}
	if !allowed[ext] {
		http.Error(w, "Invalid file type. Allowed: m4a, wav, mp3, mp4, webm", http.StatusBadRequest)
		return
	}

	// Save file
	// Use a unique name or keep original? For simplicity, we use the original but handle duplicates if needed.
	// Or just append timestamp.

	// Actually, create a unique filename to avoid collision
	safeFilename := filepath.Base(header.Filename)
	dstPath := filepath.Join("data", safeFilename)

	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}

	// Save to DB
	stmt, err := db.DB.Prepare("INSERT INTO jobs(filename, filepath, email, status) VALUES(?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(safeFilename, dstPath, email, "pending")
	if err != nil {
		http.Error(w, "DB insert error", http.StatusInternalServerError)
		return
	}

	// HTMX response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<div class='alert alert-success'>Job uploaded successfully! We will email you at %s when done.</div>", email)
}

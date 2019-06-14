package test

import (
	"github.com/fulunyong/code/file"
	"log"
	"net/http"
	"testing"
)

const maxSize = 200 * 1024 * 1024   // 200 MB
const singleSize = 20 * 1024 * 1024 // 20 MB
const uploadPath = "D:/tmp"

func Test1(t *testing.T) {
	http.HandleFunc("/upload", file.UploadFileHandler(uploadPath, maxSize, singleSize))
	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))
	log.Print("Server started on localhost:2080, use /upload for uploading files and /files/{fileName} for downloading files.")
	serve := http.ListenAndServe(":2080", nil)
	if serve != nil {
		t.Fatal(serve)
	}
}

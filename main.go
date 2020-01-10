package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func isMediaFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasSuffix(part, ".mp3") {
			return true
		}
	}
	return false
}

type httpFile struct {
	http.File
}

// Readdir is a wrapper around the Readdir method of the embedded File
// that filters out all files that start with a period in their name.
func (f httpFile) Readdir(n int) (wantedFiles []os.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".mp3") {
			wantedFiles = append(wantedFiles, file)
		}
	}
	return
}

type httpFileSystem struct {
	http.FileSystem
}

func (fs httpFileSystem) Open(name string) (http.File, error) {
	if !isMediaFile(name) {
		log.Println(name)
		return nil, os.ErrPermission
	}

	file, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return httpFile{file}, err
}

func main() {
	fs := httpFileSystem{http.Dir(".")}
	http.Handle("/", http.FileServer(fs))
	log.Fatal(http.ListenAndServe(":8081", nil))
}

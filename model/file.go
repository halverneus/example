package model

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/halverneus/example/database"
	"github.com/halverneus/example/storage"
)

var (
	// File namespace contains all file-specific functions.
	File FileNamespace
)

// Start the model service.
func Start() (wg *sync.WaitGroup) {
	wg = &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			f, ok := <-database.FileDeletionChan
			if !ok {
				return
			}
			if err := storage.Delete(f.Path); nil != err {
				log.Printf("Received error while deleting %s: %v\n", f.Path, err)
			}
		}
	}()
	return
}

// FileMetadata contains general information about file contents.
type FileMetadata struct {
	Path        string
	ContentType string
	Uploader    string
}

// FileNamespace is used to organize the controller/model functions.
type FileNamespace struct{}

// Upload a new file.
func (fn FileNamespace) Upload(meta *FileMetadata, r io.Reader) (err error) {
	// Copy to database object to assure no race condition due to misuse.
	f := &database.File{
		Path:        meta.Path,
		ContentType: meta.ContentType,
		Uploader:    meta.Uploader,
		Created:     time.Now().Format("Jan 2, 2006 3:04 PM"),
	}

	// Upload to the file system.
	if err = storage.Upload(f.Path, r); nil != err {
		return
	}

	// Push metadata into database. On failure, attempt to delete uploaded file.
	if err = database.AddFile(f); nil != err {
		storage.Delete(f.Path)
	}
	return
}

// Metadata for a file.
func (fn FileNamespace) Metadata(filePath string) (meta *FileMetadata, err error) {
	var file *database.File
	if file, err = database.GetMetadata(filePath); nil != err {
		return
	}

	// Copy to model metadata.
	meta = &FileMetadata{
		Path:        file.Path,
		ContentType: file.ContentType,
		Uploader:    file.Uploader,
	}
	return
}

// Download an existing file.
func (fn FileNamespace) Download(filePath string, w io.Writer) (meta *FileMetadata, err error) {
	// Get a lock on the file to prevent deletion while downloading.
	var f *database.File
	if f, err = database.GetFileForDownload(filePath); nil != err {
		return
	}
	defer f.Done()

	// Create metadata.
	meta = &FileMetadata{
		Path:        f.Path,
		ContentType: f.ContentType,
		Uploader:    f.Uploader,
	}

	// Download file from storage.
	err = storage.Download(filePath, w)
	return
}

// Delete a file.
func (fn FileNamespace) Delete(filePath string) error {
	return database.RemoveFile(filePath)
}

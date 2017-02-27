package storage

import (
	"io"
	"os"
	"path"

	"github.com/halverneus/example/config"
)

// Delete file from storage.
func Delete(filePath string) (err error) {
	fullpath := path.Join(config.Get.Storage.Folder, filePath)

	// Remove file.
	if err = os.Remove(fullpath); nil != err {
		return
	}

	// Remove empty parent directories.
	dir := path.Dir(fullpath)
	for isDirEmpty(dir) {
		if err = os.Remove(dir); nil != err {
			return nil // Not error. Folder removal complete.
		}
		dir = path.Clean(dir + "/..") // Move up one directory.
	}
	return
}

// isDirEmpty returns true when a provided directory is empty.
func isDirEmpty(dirPath string) bool {
	// Open directory for evaluation.
	dir, err := os.Open(dirPath)
	if nil != err {
		return false
	}
	defer dir.Close()

	// Request to read one FileInfo from directory. If none found then directory
	// is empty.
	_, err = dir.Readdir(1)
	return io.EOF == err
}

// Download a file from storage to the client.
func Download(filePath string, w io.Writer) (err error) {
	fullpath := path.Join(config.Get.Storage.Folder, filePath)

	// Open file for reading.
	var file *os.File
	if file, err = os.OpenFile(fullpath, os.O_RDONLY, 0755); nil != err {
		return
	}
	defer file.Close()

	// Send contents to client.
	_, err = io.Copy(w, file)
	return
}

// Upload a file from the client to storage.
func Upload(filePath string, r io.Reader) (err error) {
	fullpath := path.Join(config.Get.Storage.Folder, filePath)
	dir := path.Dir(fullpath)

	// Create all parent directories.
	if err = os.MkdirAll(dir, 0755); nil != err {
		return
	}

	// Open file for writing.
	var file *os.File
	if file, err = os.OpenFile(
		fullpath,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0755,
	); nil != err {
		return
	}
	defer file.Close()

	// Write to file.
	_, err = io.Copy(file, r)
	return
}

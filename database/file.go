package database

import (
	"errors"
	"sync"
)

// File stored on the system.
type File struct {
	Path        string `json:"path"`
	ContentType string `json:"content-type"`
	Uploader    string `json:"uploader"`
	Created     string `json:"created"`
	mtx         sync.RWMutex
}

// Wait for all downloaders to complete before deleting from file system.
func (f *File) Wait() {
	f.mtx.Lock()
}

// Done is called when a downloader is done downloading.
func (f *File) Done() {
	f.mtx.RUnlock()
}

// Deleted is called to free the write mutex after file is deleted.
func (f *File) Deleted() {
	f.mtx.Unlock()
}

// addFile information to the database.
func addFile(meta *File) (err error) {
	getMtx.Lock()
	defer getMtx.Unlock()

	// If file exists, overwrite.
	var orig *File
	if orig, err = getFileFromIndex(meta.Path); nil == err {
		// File exists. Replace.
		if !replaceFile(orig, meta) {
			refreshIndex()
		}
	}

	// File doesn't exist. Add file.
	get.Files = append(get.Files, meta)
	addFileToIndex(meta)

	return save()
}

// replaceFile with a new file.
func replaceFile(origFile, newFile *File) bool {
	// Replace original file with new file.
	for i, f := range get.Files {
		if f == origFile {
			get.Files[i] = newFile
			return true
		}
	}
	return false
}

// getMetadata that is safe to return.
func getMetadata(filePath string) (file *File, err error) {
	getMtx.RLock()
	defer getMtx.RUnlock()

	// Acquire file object.
	var f *File
	if f, err = getFileFromIndex(filePath); nil != err {
		return
	}

	// Copy file object to prevent risk of race conditions.
	file = &File{
		Path:        f.Path,
		ContentType: f.ContentType,
		Uploader:    f.Uploader,
		Created:     f.Created,
	}
	return
}

// getFileForDownload so that a deletion needs to wait.
func getFileForDownload(filePath string) (file *File, err error) {
	getMtx.RLock()
	defer getMtx.RUnlock()

	// Acquire file object.
	if file, err = getFileFromIndex(filePath); nil != err {
		return
	}

	// Lock file for download.
	file.mtx.RLock()
	return
}

// removeFile from the database.
func removeFile(filePath string) (err error) {
	getMtx.Lock()
	defer getMtx.Unlock()

	// Acquire file object.
	var f *File
	if f, err = getFileFromIndex(filePath); nil != err {
		return
	}

	// Find index of file.
	index := 0
	found := false
	for i, fil := range get.Files {
		if f == fil {
			index = i
			found = true
			break
		}
	}

	// Found in index, but not in database.
	if !found {
		refreshIndex()
		return
	}

	// By locking here, assures deletion completes before channel closes.
	fileDeletionMtx.RLock()
	if fileDeletionClosed {
		// Channel is closed and program is exiting.
		fileDeletionMtx.RUnlock()
		err = errors.New("application is shutting down")
		return
	}

	// Memory-leak-safe implementation for removing an item from a list.
	copy(get.Files[index:], get.Files[index+1:]) // Shift left to remove file.
	get.Files[len(get.Files)-1] = nil            // Garbage collect the trailing item.
	get.Files = get.Files[:len(get.Files)-1]     // Slice the tail from the slice.

	// Update index.
	removeFileFromIndex(filePath)

	go func() {
		f.Wait()
		FileDeletionChan <- f
		fileDeletionMtx.RUnlock() // Free FileDeletionChan for closing.
	}()

	return save()
}

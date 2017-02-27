package database

import "sync"

var (
	// FileDeletionChan is a pipe for files awaiting deletion.
	FileDeletionChan chan *File

	// fileDeletionMtx protects the entrance of the deletion channel on shutdown.
	fileDeletionMtx sync.RWMutex

	// fileDeletionClosed notifies other not to proceed with a deletion. Too late.
	fileDeletionClosed bool
)

func init() {
	fileDeletionClosed = false
	FileDeletionChan = make(chan *File, 0)
}

// Shutdown the database which flushes pending file deletions.
func Shutdown() {
	fileDeletionMtx.Lock()
	close(FileDeletionChan)
	fileDeletionClosed = true
	fileDeletionMtx.Unlock()
}

// Load the database from disk.
func Load(filename string) (err error) {
	return load(filename)
}

// AddUser to the database.
func AddUser(username, password string) (err error) {
	var u *user
	if u, err = newUser(username, password); nil != err {
		return
	}
	return addUser(u)
}

// UpdateUserPassword in database.
func UpdateUserPassword(username, password string) error {
	return setPassword(username, password)
}

// RemoveUser from the database.
func RemoveUser(username string) error {
	return removeUser(username)
}

// AuthenticateUser to allow access to the application..
func AuthenticateUser(username, password string) bool {
	return checkPassword(username, password)
}

// AddFile to be tracked by the database.
func AddFile(file *File) error {
	return addFile(file)
}

// GetMetadata return a copy of the metadata.
func GetMetadata(filePath string) (file *File, err error) {
	return getMetadata(filePath)
}

// GetFileForDownload so that the contents are not deleted in progress.
func GetFileForDownload(filePath string) (file *File, err error) {
	return getFileForDownload(filePath)
}

// RemoveFile from the database.
func RemoveFile(filePath string) error {
	return removeFile(filePath)
}

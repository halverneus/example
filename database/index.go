package database

import (
	"fmt"
	"log"
)

var (
	users map[string]*user
	files map[string]*File
)

// refreshIndex used for quick access.
func refreshIndex() {
	log.Println("Refreshing index.")
	defer log.Println("Refreshing index complete.")

	// Refresh users index.
	users = map[string]*user{}
	for _, u := range get.Users {
		users[u.Username] = u
	}

	// Refresh files.
	files = map[string]*File{}
	for _, f := range get.Files {
		files[f.Path] = f
	}
}

// addUserToIndex for new users.
func addUserToIndex(u *user) {
	users[u.Username] = u
}

// getUserFromIndex for, likely, changing passwords.
func getUserFromIndex(name string) (u *user, err error) {
	found := false
	if u, found = users[name]; !found {
		err = fmt.Errorf("user named %s not found", name)
	}
	return
}

// removeUserFromIndex for removing users.
func removeUserFromIndex(name string) {
	delete(users, name)
}

// updateUserPasswordInIndex is a quick way to change a password.
func updateUserPasswordInIndex(name, password string) (err error) {
	user, found := users[name]
	if !found {
		err = fmt.Errorf("user named %s not found for password change", name)
		return
	}
	return user.setPassword(password)
}

// addFileToIndex for a new upload.
func addFileToIndex(f *File) {
	files[f.Path] = f
}

// removeFileFromIndex for a delete.
func removeFileFromIndex(filePath string) {
	delete(files, filePath)
}

// getFileFromIndex for metadata.
func getFileFromIndex(filePath string) (f *File, err error) {
	found := false
	if f, found = files[filePath]; !found {
		err = fmt.Errorf("file at %s not found", filePath)
	}
	return
}

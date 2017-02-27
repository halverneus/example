package model

import "github.com/halverneus/example/database"

var (
	// User namespace contains all user-specific functions.
	User UserNamespace
)

// UserNamespace is used to organize the controller/model functions.
type UserNamespace struct{}

// Add a new user to the database.
func (un UserNamespace) Add(username, password string) error {
	return database.AddUser(username, password)
}

// Remove a user from the database.
func (un UserNamespace) Remove(username string) error {
	return database.RemoveUser(username)
}

// UpdatePassword for a given user.
func (un UserNamespace) UpdatePassword(username, password string) error {
	return database.UpdateUserPassword(username, password)
}

// CheckPassword for a given users.
func (un *UserNamespace) CheckPassword(username, password string) bool {
	return database.AuthenticateUser(username, password)
}

package database

import (
	"errors"
	"fmt"

	"github.com/halverneus/example/lib/encrypt"
)

// user of the system.
type user struct {
	Username string `json:"username"`
	Salt     string `json:"salt"`
	Password string `json:"password"`
}

// newUser with random salt applied.
func newUser(username, password string) (u *user, err error) {
	// Generate a new salt for a new user.
	var salt string
	if salt, err = encrypt.NewBase64Salt(); nil != err {
		return
	}

	// Create new user.
	u = &user{
		Username: username,
		Salt:     salt,
	}
	if err = u.setPassword(password); nil != err {
		u = nil
	}
	return
}

// setPassword so it is encrypted.
func (u *user) setPassword(password string) (err error) {
	if 8 > len(password) {
		err = errors.New("`password` must be at least 8 characters in length")
		return
	}

	// Encrypt password.
	var encPassword string
	if encPassword, err = encrypt.Password(password, u.Salt); nil != err {
		return
	}

	u.Password = encPassword
	return
}

// checkPassword against current.
func (u *user) checkPassword(password string) (matches bool) {
	// Encrypt password.
	var encPassword string
	var err error
	if encPassword, err = encrypt.Password(password, u.Salt); nil != err {
		return false
	}

	// Compare results.
	return encPassword == u.Password
}

// addUser to database and index and save.
func addUser(u *user) (err error) {
	getMtx.Lock()
	defer getMtx.Unlock()

	// Verify user does not exist.
	if _, err = getUserFromIndex(u.Username); nil == err {
		err = fmt.Errorf("User %s already exists", u.Username)
		return
	}

	// Add user to database and index.
	get.Users = append(get.Users, u)
	addUserToIndex(u)

	return save()
}

// removeUser from database and index and save.
func removeUser(username string) (err error) {
	getMtx.Lock()
	defer getMtx.Unlock()

	// Last user cannot be deleted.
	if 2 > len(get.Users) {
		err = errors.New("last user cannot be removed")
		return
	}

	// Acquire user object.
	var u *user
	if u, err = getUserFromIndex(username); nil != err {
		return
	}

	// Find index of user.
	index := 0
	found := false
	for i, usr := range get.Users {
		if u == usr {
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

	// Memory-leak-safe implementation for removing an item from a list.
	copy(get.Users[index:], get.Users[index+1:]) // Shift left to remove user.
	get.Users[len(get.Users)-1] = nil            // Garbage collect the trailing item.
	get.Users = get.Users[:len(get.Users)-1]     // Slice the tail from the slice.

	// Update index.
	removeUserFromIndex(username)

	return save()
}

// setPassword for a given user.
func setPassword(username, password string) (err error) {
	getMtx.Lock()
	defer getMtx.Unlock()

	// Acquire user object.
	var u *user
	if u, err = getUserFromIndex(username); nil != err {
		return
	}

	// Update password for user.
	if err = u.setPassword(password); nil != err {
		return
	}

	return save()
}

// checkPassword to verify the user is authenticated.
func checkPassword(username, password string) bool {
	getMtx.RLock()
	defer getMtx.RUnlock()

	// Acquire user object.
	u, err := getUserFromIndex(username)
	if nil != err {
		return false
	}

	// Check password for user.
	return u.checkPassword(password)
}

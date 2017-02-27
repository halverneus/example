package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

var (
	getMtx sync.RWMutex
	// get or set any value.
	get struct {

		// Users of the system.
		Users []*user `json:"users"`

		// Files stored on the system.
		Files []*File `json:"files"`
	}

	// dbFilename is the last loaded database.
	dbFilename string
)

// load the database from disk.
func load(filename string) (err error) {
	getMtx.Lock()
	defer getMtx.Unlock()

	// Read contents from configuration file.
	var contents []byte
	if contents, err = ioutil.ReadFile(filename); nil != err {
		dbFilename = filename // Write to database failed. Try saving.
		if errX := save(); nil != errX {
			err = fmt.Errorf("Multiple errors: %v & %v", err, errX)
			return
		}
		// Was able to perform initial write to database file. Retry read.
		if contents, err = ioutil.ReadFile(filename); nil != err {
			return
		}
	}

	// Parse contents into 'Get' database.
	if err = json.Unmarshal(contents, &get); nil != err {
		return
	}

	// Assign successfully parsed filename.
	dbFilename = filename

	// Populate index.
	refreshIndex()
	return
}

// save the database to disk.
func save() (err error) {
	// Read contents from structure into slice.
	var contents []byte
	if contents, err = json.Marshal(&get); nil != err {
		return
	}

	// Save the contents to disk.
	return ioutil.WriteFile(dbFilename, contents, 0666)
}

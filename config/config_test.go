package config

import (
	"io/ioutil"
	"os"
	"testing"
)

// No known way to error yaml.Marshal(&Get).

const (
	testBadConfig = `
database:
	example.db
storage:
	storage
example:
	bind: 127.0.0.1:8080
`
)

const (
	filename    = "example.yaml"
	invalidPath = "/this/path/will/never/exist.yaml"

	// Expected values when assigning.
	initialDatabaseFilename = "random"
	initialStorageFolder    = "modnar"
	initialExampleBind      = "narmod"
)

func TestGoodSave(t *testing.T) {
	defer os.Remove(filename)

	if err := Save(filename); nil != err {
		t.Fatalf("Failed to save file with: %v\n", err)
	}
}

func TestBadSave(t *testing.T) {
	defer os.Remove(invalidPath)

	if err := Save(invalidPath); nil == err {
		t.Fatal("Successfully saved while expecting failure")
	}
}

func TestGoodLoad(t *testing.T) {
	defer os.Remove(filename)

	// Assign expected values.
	Get.Database.Filename = initialDatabaseFilename
	Get.Storage.Folder = initialStorageFolder
	Get.Example.Bind = initialExampleBind

	// Save and reload new values.
	if err := Save(filename); nil != err {
		t.Fatalf("Failed to save file with: %v\n", err)
	}
	if err := Load(filename); nil != err {
		t.Fatalf("Failed to load file with: %v\n", err)
	}

	// Evaluate if values match.
	evaluate(t)
}

func TestBadLoad(t *testing.T) {
	defer os.Remove(filename)

	if err := ioutil.WriteFile(filename, []byte(testBadConfig), 0666); nil != err {
		t.Fatalf("Failed to save 'bad' file with: %v\n", err)
	}
	if err := Load(filename); nil == err {
		t.Fatal("Expected parsing error")
	}
}

func TestReadMissingFile(t *testing.T) {
	// File doesn't exist due to previous tests removing file.
	if err := Load(filename); nil == err {
		t.Fatal("Expected file load error")
	}
}

func TestEnvVarWithFile(t *testing.T) {
	defer os.Remove(filename)

	// Save a default configuration file.
	if err := Save(filename); nil != err {
		t.Fatalf("Failed to save file with: %v\n", err)
	}

	// Load configuration file with envars
	os.Setenv("EXAMPLE_DATABASE_FILENAME", initialDatabaseFilename)
	os.Setenv("EXAMPLE_STORAGE_FOLDER", initialStorageFolder)
	os.Setenv("EXAMPLE_EXAMPLE_BIND", initialExampleBind)
	if err := Load(filename); nil != err {
		t.Fatalf("Failed to load file with: %v\n", err)
	}

	// Evaluate if values match.
	evaluate(t)
}

func TestEnvVarWithoutFile(t *testing.T) {
	// Load configuration file with envars
	os.Setenv("EXAMPLE_DATABASE_FILENAME", initialDatabaseFilename)
	os.Setenv("EXAMPLE_STORAGE_FOLDER", initialStorageFolder)
	os.Setenv("EXAMPLE_EXAMPLE_BIND", initialExampleBind)
	if err := Load(""); nil != err {
		t.Fatalf("No file to load and got: %v\n", err)
	}

	// Evaluate if values match.
	evaluate(t)
}

// Evaluate if values in configuration are expected.
func evaluate(t *testing.T) {
	// check if new values match initial values.
	const errStr = "%s does not match\nExpected: %s\nGot:%s\n"
	check := func(name, initial, result string) {
		if initial != result {
			t.Fatalf(errStr, name, initial, result)
		}
	}

	// Perform checks.
	check(
		"Database filename",
		initialDatabaseFilename,
		Get.Database.Filename,
	)
	check(
		"Storage folder",
		initialStorageFolder,
		Get.Storage.Folder,
	)
	check(
		"Example bind",
		initialExampleBind,
		Get.Example.Bind,
	)
}

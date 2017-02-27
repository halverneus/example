package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	// Args from the CLI.
	Args []string
	// ConfigPath is the path to the configuration file.
	ConfigPath string
	// Help flag is a request for help.
	Help bool
	// Version flag is used to request the application version.
	Version bool
	// Debug flag is used to enable chatty logging.
	Debug bool

	// Get and configuration value.
	Get struct {

		// Database connection parameters.
		Database struct {
			Filename string `yaml:"filename"`
		} `yaml:"database"`

		// Storage for files.
		Storage struct {
			Folder string `yaml:"folder"`
		} `yaml:"storage"`

		// Example configuration values.
		Example struct {
			Bind string `yaml:"bind"`
		} `yaml:"example"`
	}
)

func init() {
	// Pre-populate with default values.
	Get.Database.Filename = "example.db"
	Get.Storage.Folder = "storage"
	Get.Example.Bind = ":8080"
}

// Load the configuration file.
func Load(filename string) (err error) {
	// If no filename provided, assign envvars.
	if "" == filename {
		overrideWithEnvVars()
		return
	}

	// Read contents from configuration file.
	var contents []byte
	if contents, err = ioutil.ReadFile(filename); nil != err {
		return
	}

	// Parse contents into 'Get' configuration.
	if err = yaml.Unmarshal(contents, &Get); nil != err {
		return
	}
	overrideWithEnvVars()
	return
}

// Save the configuration file.
func Save(filename string) (err error) {
	// Read contents from structure into slice.
	var contents []byte
	if contents, err = yaml.Marshal(&Get); nil != err {
		return
	}

	// Save the contents to disk.
	return ioutil.WriteFile(filename, contents, 0666)
}

// overrideWithEnvVars the default values and the configuration file values.
func overrideWithEnvVars() {
	// resolve function returns value of envvar if set.
	resolve := func(key, original string) string {
		if value := os.Getenv(key); "" != value {
			return value
		}
		return original
	}

	// Assign envvars, if set.
	Get.Database.Filename = resolve("EXAMPLE_DATABASE_FILENAME", Get.Database.Filename)
	Get.Storage.Folder = resolve("EXAMPLE_STORAGE_FOLDER", Get.Storage.Folder)
	Get.Example.Bind = resolve("EXAMPLE_EXAMPLE_BIND", Get.Example.Bind)
}

// Package initialize handles the "example init" command by creating a new
// configuration file with defaults assigned, unless overridden by ENVIRONMENT
// variables. The long form of the package name is due to the reserved "init()"
// function in Go.
package initialize

import (
	"errors"

	"github.com/halverneus/example/config"
)

// Configuration file to be initialized.
func Configuration() (err error) {
	// Verify configuration path is set.
	if "" == config.ConfigPath {
		msg := `configuration file not assigned; try "example` +
			` --config /desired/path/to/config.yaml init"`
		err = errors.New(msg)
		return
	}

	// Load configuration to assign defaults. Save resulting file.
	if err = config.Load(""); nil != err {
		return
	}

	// Save the configuration to disk.
	err = config.Save(config.ConfigPath)
	return
}

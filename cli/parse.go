package cli

import (
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/halverneus/example/cli/help"
	"github.com/halverneus/example/cli/initialize"
	"github.com/halverneus/example/config"
	"github.com/halverneus/example/database"
	"github.com/halverneus/example/lib/exit"
	"github.com/halverneus/example/lib/web"
	"github.com/halverneus/example/model"
	"github.com/halverneus/example/router"
)

func init() {
	flag.StringVar(&config.ConfigPath, "config", "", "")
	flag.StringVar(&config.ConfigPath, "c", "", "")
	flag.BoolVar(&config.Help, "help", false, "")
	flag.BoolVar(&config.Help, "h", false, "")
	flag.BoolVar(&config.Version, "version", false, "")
	flag.BoolVar(&config.Version, "v", false, "")
	flag.BoolVar(&config.Debug, "debug", false, "")
	flag.BoolVar(&config.Debug, "d", false, "")
}

var (
	// Version is overridden by the build script to insert an actual value.
	Version = "0.0.0"
	// GoVersion is the version of Go used to build the application.
	GoVersion = runtime.Version()
)

// Parse CLI arguments.
func Parse() (err error) {
	// Parse flag options, then parse commands arguments.
	flag.Parse()
	args := Args{}.Parse(flag.Args())

	// Push settings.
	web.Debug = config.Debug

	// Evaluate CLI input. "checkForHelp" and "checkForVersion" exit on match.
	checkForHelp(args)
	checkForVersion(args)
	err = executeCommand(args)
	return
}

// checkForHelp and print help message if found.
func checkForHelp(args Args) {
	switch {

	// "example" help.
	case config.Help:
		fallthrough
	case args.Matches() && !config.Version:
		fallthrough
	case args.Matches("help"):
		exit.With(help.Example)

		// "example init" help.
	case args.Matches("init") && config.Help:
		fallthrough
	case args.Matches("init", "help"):
		exit.With(help.ExampleInit)

		// "example user" help.
	case args.Matches("user") && config.Help:
		fallthrough
	case args.Matches("user", "help"):
		exit.With(help.ExampleUser)

		// "example user add" help.
	case args.Matches("user", "add") && config.Help:
		fallthrough
	case args.Matches("user", "add", "help"):
		exit.With(help.ExampleUserAdd)

		// "example run" help
	case args.Matches("run") && config.Help:
		fallthrough
	case args.Matches("run", "help"):
		exit.With(help.ExampleRun)
	}
}

// checkForVersion request to print version and exit.
func checkForVersion(args Args) {
	if config.Version || args.Matches("version") {
		exit.Withf("Example v%s built with %s\n", Version, GoVersion)
	}
}

// executeCommand passed on the command line.
func executeCommand(args Args) (err error) {
	switch {

	case args.Matches("init"):
		err = initialize.Configuration()

	case args.Matches("user", "add", "*", "*"):
		const userIndex, passwordIndex = 2, 3
		err = withDB(
			func(a ...string) error { return model.User.Add(a[0], a[1]) },
			args[userIndex],
			args[passwordIndex],
		)

	case args.Matches("run"):
		// Start the server.
		err = withDB(
			func(a ...string) error { return router.Run() },
		)

	default:
		err = fmt.Errorf(
			"unknown command provided: %s\n\n%s\n",
			strings.Join(args, " "),
			help.Example,
		)
	}
	return
}

// withDB loads the configuration file and the database file before calling the
// handler with the desired variables and returning a consolidated error.
func withDB(handler func(a ...string) error, a ...string) (err error) {
	if err = config.Load(config.ConfigPath); nil != err {
		return
	}
	if err = database.Load(config.Get.Database.Filename); nil != err {
		return
	}
	return handler(a...)
}

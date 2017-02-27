package help

const (

	// Example help documentation.
	Example = `
NAME
    example

USAGE
    example [ OPTIONS ] COMMAND

DESCRIPTION
    Example provides a demo of some of the architectual and syntactic lessons
    learned while developing with Go. Example provides the ability to add an
    user, authenticate and perform file storage and retrieval operations.

DEPENDENCIES
    None

OPTIONS
    -c, --config     Location of configuration file (default: "").
    -d, --debug      Prints a chatty log to the screen.
    -h, --help       Print usage.
    -v, --version    Print version information.

ENVIRONMENT VARIABLES
    EXAMPLE_DATABASE_FILENAME    Location of database file.
    EXAMPLE_STORAGE_FOLDER       Folder for storing files.
    EXAMPLE_EXAMPLE_BIND         Bind to IP address (examples: ":8080",
                                 "10.0.5.6:8080").

MANAGEMENT COMMANDS
    init      Creates an empty configuration file. Server must be stopped!
    user      Modify users. Server must be stopped!

COMMANDS
    help      Print usage.
    run       Run the server.
    version   Print version information.

FILES
    configuration (ex: example.yaml)
        This file contains all configuration settings with defaults required by
        the Example application. File must be valid YAML.

        EXAMPLE
            database:
              filename: example.db    // Path to the JSON database file.
            storage:
              folder: ./storage       // Path to the file storage folder.
            example:
              bind: ":8080"           // Bind to IP address.

    database (ex: example.db)
        JSON formatted database file that stores users and file object metadata.
`

	// ExampleInit help documentation.
	ExampleInit = `
NAME
    example [ OPTIONS ] init

USAGE
    example --config /path/to/example.yaml init

DESCRIPTION
    Creates an empty configuration file. Requires configuration file location to
    be set. Server must be stopped before running this command.

OPTIONS
    -c, --config     Location of configuration file (default: "").
    -h, --help       Print usage.

COMMANDS
    help             Print usage.
`

	// ExampleUser help documentation.
	ExampleUser = `
NAME
    example [ OPTIONS ] user

USAGE
    example user COMMAND

DESCRIPTION
    Allows an administrator to modify users from the command line. Server must
    be stopped before running this command.

OPTIONS
    -c, --config     Location of configuration file (default: "").
    -h, --help       Print usage.

ENVIRONMENT VARIABLES
    EXAMPLE_DATABASE_FILENAME    Location of database file.

MANAGEMENT COMMANDS
    add      Add a new user to the system.

COMMANDS
    help      Print usage.
`

	// ExampleUserAdd help documentation.
	ExampleUserAdd = `
NAME
    example [ OPTIONS ] user add

USAGE
    example user add [ username ] [ password ]

DESCRIPTION
    Allows an administrator to add a new user from the command line. Server must
    be stopped before running this command.

OPTIONS
    -c, --config     Location of configuration file (default: "").
    -h, --help       Print usage.

ENVIRONMENT VARIABLES
    EXAMPLE_DATABASE_FILENAME    Location of database file.

COMMANDS
    help      Print usage.
`

	// ExampleRun help documentation.
	ExampleRun = `
NAME
    example [ OPTIONS ] run

USAGE
    example run

DESCRIPTION
    Runs the Example server.

OPTIONS
    -c, --config     Location of configuration file (default: "").
    -h, --help       Print usage.

ENVIRONMENT VARIABLES
    EXAMPLE_DATABASE_FILENAME    Location of database file.
    EXAMPLE_STORAGE_FOLDER       Folder for storing files.
    EXAMPLE_EXAMPLE_BIND         Bind to IP address (examples: ":8080",
                                 "10.0.5.6:8080").

COMMANDS
    help      Print usage.
`
)

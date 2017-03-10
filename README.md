# Example
An Object Storage implementation example written in Go.

## What it is?
Example is exactly what it sounds like: an example. After creating several other applications, I decided to throw together a little example application that others may enjoy using for learning or contributing. Anyone interested in contributing is welcome to place a pull request.

## Requirements
* Go 1.8

## How to get started
To help people get quickly spun up without invading their computer, there are two methods available to install. The first is a conventional method for those with a setup Go environment while the second is for those without any Go resources or knowledge.
### (Install) I have a Go development environment
Download and install using the conventional 'go get' method. The second line is for those that want to set a version number.
```bash
go get -u github.com/halverneus/example/bin/example
go install -ldflags "-X github.com/halverneus/example/cli.Version=$VERSION" github.com/halverneus/example/bin/example
```

### (Install) I do NOT have a Go development environment
On an Ubuntu 16.04 machine with 'wget' installed ('sudo apt-get -y install wget') Run the following to install Example:
```bash
wget https://github.com/halverneus/example/raw/master/up && chmod +x up && ./up setup
```

## Configuration
Example can use, either, environment variables or a YAML file configuration file to assign settings. This section is, again, divided between the Go developer and everyone else.

### (Configuration) The Go developer
With the binary now installed in you $PATH, run the following to create a configuration file:
```bash
example -c config.yaml init # This will create an empty configuration file.
```

Configuration file contents with defaults:
```yaml
database:
  filename: example.db
storage:
  folder: storage
example:
  bind: :8080
```
Description:
* database.filename -> Location to put the JSON database.
* storage.folder    -> Location of the folder to store files.
* example.bind      -> Network binding address. (":8080", "127.0.0.1:8888", "the.host:8088")

NOTE: Additionally, configuration can also be set with environment variables as follows (using defaults):
```bash
export EXAMPLE_DATABASE_FILENAME="example.db"
export EXAMPLE_STORAGE_FOLDER="storage"
export EXAMPLE_EXAMPLE_BIND=":8080"
```

After setting up and saving the configuration file, create your first user as follows (replacing "username" and "password" with your own):
```bash
example -c config.yaml user add username password
```

### (Configuration) NOT a Go developer
Configuration consists of adding the first user as follows:
```bash
./up init
```
The command will ask for the username and password of the first user.

## Running Example
Again, running example is slightly different depending on whether the 'up' script is used, so this section is again, divided.

### (Run) Go developer
Execute the following command to run as a Go developer:
```bash
example -c config.yaml run
```

Exit using ctrl+c.

### (Run) NOT a Go developer
Execute the following:
```bash
./up run
```

Exit using ctrl+c

## The 'up' script.
If the automatic script isn't desired, feel free to review the script and perform any necessary steps manually. Additionally, the 'up' script can:
```bash
up setup: Create an empty directory. It is okay to have this script in the directory. cd into the directory and execute './up setup'. This will install Go, all Example dependencies and Example.

up init: Performs initial setup routines, including creating a new user.

up update: This will pull down the latest Example code, update dependencies and reinstall Example.

up build: If any changes are made to the source code, this will build it.

up run: Starts the server.

up test: Performs all unit tests on the code.

up wipe: Deletes all files except for the 'up' file.
```

## NOT a Go developer development
If there is a desire to make changes to the code that was pulled using the 'up' script, simply source the source.me.sh file and all normal 'go' commands will be available
```bash
source source.me.sh
go run bin/example/main.go -c config.yaml run
```

## Additional help
For more assistance, see:
```bash
example help
example --help
example -help
example -h
example init help     # For command-specific help.
```

## Using the service
The following are a series of commands that can be executed to perform various functions.

Adding a first user:
```bash
example --config config.yaml user add yourname yourpassword
```

Adding more users:
```bash
curl --user yourname:yourpassword -X PUT \
    --data '{"username":"othername","password":"otherpassword"}' \
    http://127.0.0.1:8080/api/latest/user
# http://127.0.0.1:8080/api/v1/user is equally valid.
```

Update your password:
```bash
curl --user yourname:yourpassword -X POST \
    --data '{"password":"yournewpassword"}' \
    http://127.0.0.1:8080/api/latest/user
# http://127.0.0.1:8080/api/v1/user is equally valid.
```

Delete another user:
```bash
curl --user yourname:yourpassword -X DELETE \
    --data '{"username":"othername"}' \
    http://127.0.0.1:8080/api/latest/user
# http://127.0.0.1:8080/api/v1/user is equally valid.
```

Uploading a file:
```bash
curl --user yourname:yourpassword --upload-file my.pdf \
    -H "Content-Type: application/pdf" \
    http://127.0.0.1:8080/api/latest/file/random/folders/your.pdf
# http://127.0.0.1:8080/api/v1/file/random/folders/your.pdf is equally valid
```

Downloading a file:
```bash
curl --user yourname:yourpassword \
    http://127.0.0.1:8080/api/latest/file/random/folders/your.pdf > their.pdf
# http://127.0.0.1:8080/api/v1/file/random/folders/your.pdf > their.pdf is equally valid
```

Deleting a file:
```bash
curl --user yourname:yourpassword -X "DELETE" \
    http://127.0.0.1:8080/api/latest/file/random/folders/your.pdf
# http://127.0.0.1:8080/api/v1/file/random/folders/your.pdf is equally valid
```

## Code layout
Quick code layout explanation:
* api -> Everything in this folder relates to the URL address. For example, api/file/get.go refers to a HTTP GET request to http(s)://{host}/api/latest/file/* or http(s)://{host}/api/v1/file/*
* bin -> Contains the main.go executable.
* cli -> Command-line interface.
* config -> Configuration settings for running the application.
* database -> Cheesy JSON file database.
* lib/authenticate -> Authentication middleware for all requests.
* lib/encrypt -> Password encryption package.
* lib/exit -> Convenient exit handler.

* lib/web -> Convenience wrapper around http.Handler calls. Middleware that closes request bodies and more.
* model -> Simplified calls permitting reusable data manipulations.
* router -> Handles routing of API calls.
* storage -> File system interacting library for Object Storage.
* vendor -> Dependencies to ignore.

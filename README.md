# Example
An Object Storage implementation example written in Go.

## What it is?
Example is exactly what it sounds like: an example. After creating several other applications, I decided to throw together a little example application that others may enjoy using for learning or contributing. Anyone interested in contributing is welcome to place a pull request.

## Requirements
* Go 1.8

## How to get started
On an Ubuntu 16.04 machine, create an empty directory and copy the "up" file beside it. Review and update variables at the top of the file. Run the following:
```
../up setup
```
If the automatic script isn't desired, feel free to review the script and perform any necessary steps manually. Additionally, the 'up' script can:
```
up setup     # Installs everything from scratch in an empty directory.
up update    # Updates to the latest release of Example.
up install   # Builds and installs Example into GOBIN.
up test      # Runs all unit tests and reports output to the screen.
```

After installation, either setup your environment variables or everytime you open a new terminal run:
```
source source.me.sh
```
...to set your environment variables. This line is included in the 'up' script.

For more assistance, see:
```
example help
example --help
example -help
example -h
example --help init     # For command-specific help.
```

## Using the service
The following are a series of commands that can be executed to perform various functions.

Adding a first user:
```
example --config config.yaml user add yourname yourpassword
```

Adding more users:
```
curl --user yourname:yourpassword -X PUT \
    --data '{"username":"othername","password":"otherpassword"}' \
    http://127.0.0.1:8080/api/latest/user
# http://127.0.0.1:8080/api/v1/user is equally valid.
```

Update your password:
```
curl --user yourname:yourpassword -X POST \
    --data '{"password":"yournewpassword"}' \
    http://127.0.0.1:8080/api/latest/user
# http://127.0.0.1:8080/api/v1/user is equally valid.
```

Delete another user:
```
curl --user yourname:yourpassword -X DELETE \
    --data '{"username":"othername"}' \
    http://127.0.0.1:8080/api/latest/user
# http://127.0.0.1:8080/api/v1/user is equally valid.
```

Uploading a file:
```
curl --user yourname:yourpassword --upload-file my.pdf \
    -H "Content-Type: application/pdf" \
    http://127.0.0.1:8080/api/latest/file/random/folders/your.pdf
# http://127.0.0.1:8080/api/v1/file/random/folders/your.pdf is equally valid
```

Downloading a file:
```
curl --user yourname:yourpassword \
    http://127.0.0.1:8080/api/latest/file/random/folders/your.pdf > their.pdf
# http://127.0.0.1:8080/api/v1/file/random/folders/your.pdf > their.pdf is equally valid
```

Deleting a file:
```
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

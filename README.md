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
If the automatic script isn't desired, feel free to review the script and perform any necessary steps manually.

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
    http://127.0.0.1:8080/api/user
```
Update your password:
curl --user yourname:yourpassword -X POST \
    --data '{"password":"yournewpassword"}' \
    http://127.0.0.1:8080/api/user
```
Delete another user:
```
curl --user yourname:yourpassword -X DELETE \
    --data '{"username":"othername"}' \
    http://127.0.0.1:8080/api/user
```
Uploading a file:
```
curl --user yourname:yourpassword --upload-file my.pdf \
    -H "Content-Type: application/pdf" \
    http://127.0.0.1:8080/api/file/random/folders/your.pdf
```
Downloading a file:
```
curl --user yourname:yourpassword \
    http://127.0.0.1:8080/api/file/random/folders/your.pdf > their.pdf
```
Deleting a file:
```
curl --user yourname:yourpassword -X "DELETE" \
    http://127.0.0.1:8080/api/file/random/folders/your.pdf
```

## Code layout
Quick code layout explanation:
* api -> Everything in this folder relates to the URL address. For example, api/file/get.go refers to a HTTP GET request to http(s)://{host}/api/file/*
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

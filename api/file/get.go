package file

import (
	"net/http"

	"github.com/halverneus/example/lib/web"
	"github.com/halverneus/example/model"
)

// GetRequest is a file stream. File is loaded from a path specified in the URL.
// For example, to download a file called "my/folder/file.json", one would
// stream the file from the following endpoint: "/api/file/my/folder/file.json".
// Credentials required.

// GET file from storage.
func GET(ctx *web.Context) {
	filePath := ctx.PS.ByName("filepath")

	// Retrieve metadata for sending headers.
	metadata, err := model.File.Metadata(filePath)
	if nil != err {
		ctx.Respond().Status(http.StatusNotFound).With(err).Do()
		return
	}

	// Assign headers and retrieve writer.
	writer := ctx.Respond().Add(web.ContentType, metadata.ContentType).Stream()

	// The writer is passed in to prevent callers from risking a deadlock.
	if _, err = model.File.Download(filePath, writer); nil != err {
		ctx.Logf("Error while downloading: %v\n", err)
	}
	return
}

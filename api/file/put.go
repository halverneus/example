package file

import (
	"net/http"

	"github.com/halverneus/example/lib/web"
	"github.com/halverneus/example/model"
)

// PutRequest is a file stream. File is saved at a path specified in the URL.
// For example, to save a file as "my/folder/file.json", one would set the
// "Content-Type" header to "application/json" and stream the file to the
// following endpoint: "/api/file/my/folder/file.json". Credentials required.

// PutResponse returns nothing.
type PutResponse struct{}

// PUT file into storage.
func PUT(ctx *web.Context) {
	filePath := ctx.PS.ByName("filepath")

	// Collect metadata information about the file.
	meta := &model.FileMetadata{
		Path:        filePath,
		ContentType: ctx.R.Header.Get(web.ContentType),
		Uploader:    ctx.User,
	}

	// Upload the file with the metadata.
	if err := model.File.Upload(meta, ctx.Reader()); nil != err {
		ctx.Respond().Status(http.StatusTeapot).With(err).Do()
	}

	// Reply with success.
	resp := &PutResponse{}
	ctx.Respond().With(resp).Do()
}

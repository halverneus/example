package file

import (
	"net/http"

	"github.com/halverneus/example/lib/web"
	"github.com/halverneus/example/model"
)

// DeleteRequest is just a URL call. File exists at the path specified in the
// URL. For example, to delete a file called "my/folder/file.json", one would
// delete the file from the following endpoint: "/api/file/my/folder/file.json".
// Credentials required.

// DeleteResponse returns nothing.
type DeleteResponse struct{}

// DELETE file from storage.
func DELETE(ctx *web.Context) {
	filePath := ctx.PS.ByName("filepath")

	// Delete file.
	if err := model.File.Delete(filePath); nil != err {
		ctx.Respond().Status(http.StatusNotFound).With(err).Do()
	}

	// Reply with success.
	resp := &DeleteResponse{}
	ctx.Respond().With(resp).Do()
}

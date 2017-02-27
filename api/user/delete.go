package user

import (
	"errors"
	"net/http"

	"github.com/halverneus/example/lib/web"
	"github.com/halverneus/example/model"
)

// DeleteRequest is the expected format of the client request.
type DeleteRequest struct {
	Username string `json:"username"`
}

// Err is a validation check on the request message.
func (req *DeleteRequest) Err() error {
	if "" == req.Username {
		return errors.New("'username' was not supplied")
	}
	return nil
}

// DeleteResponse returns nothing.
type DeleteResponse struct{}

// DELETE user from the database
func DELETE(ctx *web.Context) {
	// Deserialize request into DeleteRequest.
	req := &DeleteRequest{}
	var err error
	if err = ctx.Decode(req); nil != err {
		ctx.Respond().Status(http.StatusBadRequest).With(err).Do()
		return
	}

	// Check that request is valid.
	if err = req.Err(); nil != err {
		ctx.Respond().Status(http.StatusBadRequest).With(err).Do()
		return
	}

	// Verify requesting user isn't deleting themselves.
	if req.Username == ctx.User {
		msg := "User cannot remove themselves from the dataabase."
		ctx.Respond().Status(http.StatusConflict).With(msg).Do()
		return
	}

	// Remove user.
	if err = model.User.Remove(req.Username); nil != err {
		ctx.Respond().Status(http.StatusNotFound).With(err).Do()
		return
	}

	// Reply with success.
	resp := &DeleteResponse{}
	ctx.Respond().With(resp).Do()
}

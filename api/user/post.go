package user

import (
	"errors"
	"net/http"

	"github.com/halverneus/example/lib/web"
	"github.com/halverneus/example/model"
)

// PostRequest is the expected format of the client request.
type PostRequest struct {
	Password string `json:"password"`
}

// Err is a validation check on the request message.
func (req *PostRequest) Err() error {
	if 8 > len(req.Password) {
		return errors.New("'password' must be at least 8 characters long")
	}
	return nil
}

// PostResponse returns nothing.
type PostResponse struct{}

// POST updates user password in the database
func POST(ctx *web.Context) {
	// Deserialize request into DeleteRequest.
	req := &PostRequest{}
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

	// Update password.
	if err = model.User.UpdatePassword(ctx.User, req.Password); nil != err {
		ctx.Respond().Status(http.StatusInternalServerError).With(err).Do()
		return
	}

	// Reply with success.
	resp := &PostResponse{}
	ctx.Respond().With(resp).Do()
}

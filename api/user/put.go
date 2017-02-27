package user

import (
	"errors"
	"net/http"

	"github.com/halverneus/example/lib/web"
	"github.com/halverneus/example/model"
)

// PutRequest is the expected format of the client request.
type PutRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Err is a validation check on the request message.
func (req *PutRequest) Err() error {
	if "" == req.Username {
		return errors.New("'username' was not supplied")
	}
	if 8 > len(req.Password) {
		return errors.New("'password' must be at least 8 characters long")
	}
	return nil
}

// PutResponse returns nothing.
type PutResponse struct{}

// PUT new user into the database
func PUT(ctx *web.Context) {
	// Deserialize request into PutRequest.
	req := &PutRequest{}
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

	// Add new user.
	if err = model.User.Add(req.Username, req.Password); nil != err {
		ctx.Respond().Status(http.StatusConflict).With(err).Do()
		return
	}

	// Reply with success.
	resp := &PutResponse{}
	ctx.Respond().With(resp).Do()
}

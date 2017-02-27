package authenticate

import (
	"net/http"

	"github.com/halverneus/example/lib/web"
	"github.com/halverneus/example/model"
)

// User request in need of authentication.
func User(handler func(*web.Context)) func(*web.Context) {
	return func(ctx *web.Context) {

		// Check password and return on failure.
		if !model.User.CheckPassword(ctx.User, ctx.Password) {
			msg := "Invalid username/password combination."
			ctx.Respond().Status(http.StatusForbidden).With(msg).Do()
			return
		}

		// Call handler.
		handler(ctx)
	}
}

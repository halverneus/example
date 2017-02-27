package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	// Debug gets set to 'true' if chatty logging is desired.
	Debug bool
)

const (
	// Unknown content type.
	Unknown = iota
	// JSON content type.
	JSON
	// Gobs content type.
	Gobs
	// Stream content type.
	Stream
)

const (
	// ContentType is used for retrieving and setting the Content-Type header.
	ContentType = "Content-Type"
	// JSONContent is used for setting the JSON Content-Type.
	JSONContent = "application/json"
)

// Context provides a simplified interface for handling responding and logging.
type Context struct {
	W        http.ResponseWriter
	R        *http.Request
	PS       httprouter.Params
	User     string
	Password string
	stream   bool
}

// New Context constructor.
func New(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	user string,
	password string,
) *Context {
	ctx := &Context{W: w, R: r, PS: ps, User: user, Password: password}
	ctx.Debugln("New incoming request.")
	return ctx
}

// Decode request into passed interface.
func (ctx *Context) Decode(v interface{}) error {
	decoder := json.NewDecoder(ctx.R.Body)
	return decoder.Decode(v)
}

// Reader returns the request reader.
func (ctx *Context) Reader() io.Reader {
	return ctx.R.Body
}

// Writer returns the response writer.
func (ctx *Context) Writer() io.Writer {
	ctx.stream = true
	return ctx.W
}

// Debug a message without a newline.
func (ctx *Context) Debug(v ...interface{}) {
	if Debug {
		ctx.Log(v...)
	}
}

// Debugf a formatted message.
func (ctx *Context) Debugf(format string, v ...interface{}) {
	if Debug {
		ctx.Logf(format, v...)
	}
}

// Debugln a message and a newline.
func (ctx *Context) Debugln(v ...interface{}) {
	if Debug {
		ctx.Logln(v...)
	}
}

// Log a message without a newline.
func (ctx *Context) Log(v ...interface{}) {
	format := "(%s) %s [%s] %s"
	log.Printf(
		format,
		ctx.User,
		ctx.R.URL.String(),
		ctx.R.Method,
		fmt.Sprint(v...),
	)
}

// Logf a formatted message.
func (ctx *Context) Logf(format string, v ...interface{}) {
	ctx.Log(fmt.Sprintf(format, v...))
}

// Logln a message and add a newline.
func (ctx *Context) Logln(v ...interface{}) {
	ctx.Log(fmt.Sprintln(v...))
}

// Respond to the client. This is the sole constructor for a Response.
func (ctx *Context) Respond() *Response {
	return &Response{
		ctx:     ctx,
		status:  http.StatusOK,
		headers: map[string]string{},
	}
}

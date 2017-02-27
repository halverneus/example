package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Response builder for HTTP requests.
type Response struct {
	ctx     *Context
	status  int
	headers map[string]string
	msg     interface{}
}

// Status code to be returned for the request.
func (resp *Response) Status(status int) *Response {
	resp.status = status
	return resp
}

// Add header to response.
func (resp *Response) Add(key, value string) *Response {
	resp.headers[key] = value
	return resp
}

// With the passed in message, encode and create response.
func (resp *Response) With(msg interface{}) *Response {
	resp.msg = msg
	return resp
}

// Stream returns a writer for downloading. Do not call 'Do'.
func (resp *Response) Stream() io.Writer {
	// Set all headers.
	for k, v := range resp.headers {
		resp.ctx.W.Header().Add(k, v)
	}
	resp.ctx.W.WriteHeader(resp.status)

	// Return writer for streaming.
	return resp.ctx.W
}

// Do finalizes the response.
func (resp *Response) Do() {
	// This defer function runs at any return and logs the result of the exchange.
	var err error
	defer func() {
		// First priority: log any errors that occurred.
		if nil != err {
			resp.ctx.Logf("Received error while responding: %v\n", err)
			return
		}

		// Second priority: log any unsuccessful requests that were made.
		if http.StatusOK != resp.status {
			resp.ctx.Logf(
				"Response code %s: %s\n",
				http.StatusText(resp.status),
				fmt.Sprint(resp.msg),
			)
			return
		}

		// Default priority: log successful completion.
		resp.ctx.Debugln("Request completed successfully")
	}()

	// MESSAGE RESPONSE HERE
	// Set all headers.
	for k, v := range resp.headers {
		resp.ctx.W.Header().Add(k, v)
	}
	resp.ctx.W.WriteHeader(resp.status)

	// Convert error messages to strings to avoid 'err.Error()' throughout app.
	if errX, ok := resp.msg.(error); ok {
		resp.msg = errX.Error()
	}

	// Encode and write JSON.
	var raw []byte
	if raw, err = json.Marshal(resp.msg); nil != err {
		return
	}
	_, err = resp.ctx.W.Write(raw)
}

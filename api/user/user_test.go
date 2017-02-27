package user

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/halverneus/example/database"
	"github.com/halverneus/example/lib/web"
	"github.com/julienschmidt/httprouter"
)

// testRequest is a stand-in for [Post|Put|Get]Request. Also can make bad JSON.
type testRequest struct {
	UN  string `json:"username"`
	PW  string `json:"password"`
	bad bool
}

// reader returns either, a bad JSON reader, or a valid JSON reader depending on
// the value of testRequest.bad.
func (tr *testRequest) reader() io.Reader {
	// Return corrupt JSON.
	if tr.bad {
		return bytes.NewReader([]byte(`{"username": "12345678}`))
	}

	// Return valid JSON.
	raw, err := json.Marshal(tr)
	if nil != err {
		log.Fatalf("While parsing JSON: %v\n", err)
	}
	return bytes.NewReader(raw)
}

// TestAll methods for /api/user.
func TestAll(t *testing.T) {

	// All test cases to be performed.
	testCases := []struct {
		name    string
		method  string
		body    *testRequest
		status  int
		result  string
		ctxUser string
	}{
		{"Add user successfully", "PUT", &testRequest{UN: "admin", PW: "12345678"}, http.StatusOK, "{}", "admin"},
		{"Add user bad JSON", "PUT", &testRequest{bad: true}, http.StatusBadRequest, "*", "admin"},
		{"Add invalid user", "PUT", &testRequest{UN: "", PW: "12345678"}, http.StatusBadRequest, "*", "admin"},
		{"Add invalid password", "PUT", &testRequest{UN: "john", PW: "1234567"}, http.StatusBadRequest, "*", "admin"},
		{"Add second user successfully", "PUT", &testRequest{UN: "john", PW: "12345678"}, http.StatusOK, "{}", "admin"},
		{"Add existing user", "PUT", &testRequest{UN: "john", PW: "87654321"}, http.StatusConflict, "*", "admin"},

		{"Change password successfully", "POST", &testRequest{PW: "87654321"}, http.StatusOK, "{}", "admin"},
		{"Change password bad JSON", "POST", &testRequest{bad: true}, http.StatusBadRequest, "*", "admin"},
		{"Change invalid password", "POST", &testRequest{PW: "1234567"}, http.StatusBadRequest, "*", "admin"},
		{"Change password bad user", "POST", &testRequest{PW: "12345678"}, http.StatusInternalServerError, "*", "alex"},

		{"Delete user bad JSON", "DELETE", &testRequest{bad: true}, http.StatusBadRequest, "*", "admin"},
		{"Delete user invalid user", "DELETE", &testRequest{UN: ""}, http.StatusBadRequest, "*", "admin"},
		{"Delete user delete self", "DELETE", &testRequest{UN: "admin"}, http.StatusConflict, "*", "admin"},
		{"Delete user bad user", "DELETE", &testRequest{UN: "alex"}, http.StatusNotFound, "*", "admin"},
		{"Delete user successfully", "DELETE", &testRequest{UN: "john"}, http.StatusOK, "{}", "admin"},
	}

	// Load the database. Use default configuration values. Delete database and
	// storage folder on completion.
	if err := database.Load("example.db"); nil != err {
		t.Errorf("While loading database: %v\n", err)
		return
	}
	defer os.Remove("example.db")
	defer os.RemoveAll("storage")

	// Setup routes to API calls.
	router := httprouter.New()
	router.DELETE("/api/user", web.Wrap(DELETE))
	router.POST("/api/user", web.Wrap(POST))
	router.PUT("/api/user", web.Wrap(PUT))

	// Start server.
	server := httptest.NewServer(router)
	defer server.Close()

	// client is used for making requests.
	client := &http.Client{}

	// Run all test cases.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare a request and setup basic authentication with arbitrary values.
			req, err := http.NewRequest(tc.method, server.URL+"/api/user", tc.body.reader())
			if nil != err {
				t.Errorf("Failed to create request with: %v\n", err)
				return
			}
			req.SetBasicAuth(tc.ctxUser, "12345678")

			// Perform request, not forgetting to close the response body.
			var resp *http.Response
			if resp, err = client.Do(req); nil != err {
				t.Error("Failed to receive response with: %v\n", err)
				return
			}
			defer resp.Body.Close()

			// Pull all contents from body of response and...
			var rawResp []byte
			if rawResp, err = ioutil.ReadAll(resp.Body); nil != err {
				t.Error("Failed to read response with: %v\n", err)
				return
			}

			// ... compare those contents against expected contents. '*' is wild.
			if "*" != tc.result && string(rawResp) != tc.result {
				t.Errorf("Response mismatch. Expected %s and got %s\n", tc.result, string(rawResp))
				t.FailNow()
			}

			// Compare error codes.
			if tc.status != resp.StatusCode {
				t.Errorf(
					"Status code mismatch. Expected %s and got %s\n",
					http.StatusText(tc.status),
					http.StatusText(resp.StatusCode),
				)
				t.FailNow()
			}
		})
	}
}

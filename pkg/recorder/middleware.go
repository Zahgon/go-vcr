package recorder

import (
	"net/http"
	"net/http/httptest"
)

// HTTPMiddleware intercepts and records all incoming requests and the server's response
func (rec *Recorder) HTTPMiddleware(next http.Handler) http.Handler {
	_ = "STUB: not implemented"
	return *new(http.Handler)
}

// Tee the body so it can be read by the next handler and by the recorder

// On the server side, requests do not have Host and Scheme so it must be set

// copy headers from real response

var _ http.ResponseWriter = &passthroughWriter{}

// passthroughWriter uses the original ResponseWriter and an httptest.ResponseRecorder
// so the middleware can capture response details and passthrough to the client
type passthroughWriter struct {
	recorder *httptest.ResponseRecorder
	real     http.ResponseWriter
}

func newPassthrough(real http.ResponseWriter) passthroughWriter {
	_ = "STUB: not implemented"
	return *new(passthroughWriter)
}

func (p passthroughWriter) Header() http.Header {
	_ = "STUB: not implemented"
	return *new(http.Header)
}

func (p passthroughWriter) Write(in []byte) (int, error) { _ = "STUB: not implemented"; return 0, nil }

func (p passthroughWriter) WriteHeader(statusCode int) { _ = "STUB: not implemented"; return }

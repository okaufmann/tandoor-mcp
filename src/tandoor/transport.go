package tandoor

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// authTransport is an http.RoundTripper middleware that attaches the Authorization header to all requests
type authTransport struct {
	token string
	next  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqCopy := req.Clone(req.Context())
	if t.token != "" {
		reqCopy.Header.Set("Authorization", t.token)
	}
	return t.next.RoundTrip(reqCopy)
}

// loggingTransport is an http.RoundTripper middleware that logs request and response bodies
type loggingTransport struct {
	next http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var reqBody []byte
	if req.Body != nil {
		var err error
		reqBody, err = io.ReadAll(req.Body)
		if err == nil {
			req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}
	}

	log.Printf("HTTP Request: %s %s | Body: %s", req.Method, req.URL, string(reqBody))

	resp, err := t.next.RoundTrip(req)
	if err != nil {
		log.Printf("HTTP Request failed: %v", err)
		return nil, err
	}

	var respBody []byte
	if resp.Body != nil {
		var readErr error
		respBody, readErr = io.ReadAll(resp.Body)
		if readErr == nil {
			resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
		}
	}

	log.Printf("HTTP Response: %d | Body: %s", resp.StatusCode, string(respBody))

	return resp, nil
}


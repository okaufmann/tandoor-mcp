package tandoor

import "net/http"

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

package jira

import (
	"net/http"
)

type (
	BasicAuthRoundTripper struct {
		username string
		password string
		next     http.RoundTripper
	}
)

func NewBasicAuthRoundTripper(username, password string, next http.RoundTripper) http.RoundTripper {
	return &BasicAuthRoundTripper{
		username: username,
		password: password,
		next:     next,
	}
}

func (b BasicAuthRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	request.SetBasicAuth(b.username, b.password)
	return b.next.RoundTrip(request)
}

package internal

import (
	"net/http"
	"time"
)

type (
	BasicAuthRoundTripper struct {
		username string
		password string
		next     http.RoundTripper
	}

	LoggerRoundTripper struct {
		next   http.RoundTripper
		before RequestLogger
		after  ResponseLogger
	}

	RequestLogger  func(req *http.Request)
	ResponseLogger func(res *http.Response, duration time.Duration)
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

func (l LoggerRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	start := time.Now()
	if l.before != nil {
		l.before(request)
	}
	response, err := l.next.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if l.after != nil {
		l.after(response, time.Since(start))
	}
	return response, nil
}

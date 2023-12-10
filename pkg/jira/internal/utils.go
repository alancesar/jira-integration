package internal

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	jiraDateTimeTemplate = "2006-01-02T15:04:05.999-0700"
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

func MustParseTimeRFC3339WithTimezone(value string) time.Time {
	return MustParseTime(jiraDateTimeTemplate, value)
}

func MustParseTimeRFC3339(value string) time.Time {
	return MustParseTime(time.RFC3339, value)
}

func MustParseTime(layout, value string) time.Time {
	if value == "" {
		return time.Time{}
	}

	parsed, _ := time.Parse(layout, value)
	return parsed
}

func MustParseURL(raw string) *url.URL {
	parsed, _ := url.Parse(raw)
	return parsed
}

func ParseStringToUint(raw string) uint {
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}

	return uint(parsed)
}

package mocks

import (
	"io"
	"net/http"
	"strings"
)

type (
	MockedRoundTripper struct {
		response   string
		statusCode int
	}
)

func (m MockedRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     http.StatusText(m.statusCode),
		StatusCode: m.statusCode,
		Header: http.Header{
			"Content-IssueType": {"application/json"},
		},
		Body:          io.NopCloser(strings.NewReader(m.response)),
		ContentLength: int64(len(m.response)),
		Request:       request,
	}, nil
}

func NewMockedRoundTripper(response string, statusCode int) *MockedRoundTripper {
	return &MockedRoundTripper{
		response:   response,
		statusCode: statusCode,
	}
}

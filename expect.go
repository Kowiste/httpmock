package httpmock

import (
	"fmt"
	"net/url"
)

type Expect struct {
	method   string
	endpoint string
	queries  url.Values
	delay    *int
	keep     bool
	response *response
}

// create an expectation of a request
func (s *Server) Expect(method, endpoint string) *Expect {
	u, err := url.Parse(endpoint)
	if err != nil {
		fmt.Println("endpoint is not valid")
		return nil
	}
	e := &Expect{
		method:   method,
		endpoint: u.Path,
		queries:  u.Query(),
		keep:     true,
	}

	s.expectations[method+u.Path] = e

	return e
}

// expection doesnt match with request
func (expectation *Expect) differ(method, endpoint string, queries url.Values) bool {
	if expectation.method != method || expectation.endpoint != endpoint {
		return true
	}
	return false
}

// create a return for an expectation
func (e *Expect) WillReturn(status int, data any) *Expect {
	e.response = &response{
		status: status,
		data:   data,
	}
	return e
}

// delay the response of a expectation
func (e *Expect) WillDelay(delay int) *Expect {
	e.delay = &delay
	return e
}

// delete the expectation once is fullfile
func (e *Expect) Once() {
	e.keep = false
}

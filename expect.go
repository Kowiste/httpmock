package httpmock

import "net/url"

type Expect struct {
	method   string
	endpoint string
	delay    *int
	keep     bool
	response *response
}

// create an expectation of a request
func (s *Server) Expect(method, endpoint string) *Expect {
	e := &Expect{
		method:   method,
		endpoint: endpoint,
		keep:     true,
	}
	s.expectations = append(s.expectations, e)
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

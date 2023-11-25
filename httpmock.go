package httpmock

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

type Server struct {
	conn         *httptest.Server
	expectations map[string]*Expect
}

// create a mock server that allow insecure connection
func New() *Server {
	s := new(Server)
	s.conn = httptest.NewUnstartedServer(s)
	s.conn.TLS = &tls.Config{
		InsecureSkipVerify: true,
	}
	s.expectations=make(map[string]*Expect)
	s.conn.StartTLS()
	return s

}

// close the mock http server
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path
	queries := r.URL.Query()
	method := r.Method
	for _, expectation := range s.expectations {
		if expectation.differ(method, endpoint, queries) {
			continue
		}
		if expectation.delay != nil {
			time.Sleep(time.Duration(*expectation.delay * int(time.Millisecond)))
		}

		expectation.response.setResponse(w)
		if !expectation.keep {
			delete(s.expectations, expectation.method+expectation.endpoint)
		}

		return
	}
	w.WriteHeader(http.StatusInternalServerError) // You can replace http.StatusOK with any other HTTP status code
	w.Write([]byte("expectation not found"))
	fmt.Println("Expectation not found for method: ", method, " endpoint: ", endpoint)

}

// close the mock http server
func (s *Server) GetURL() string {
	return s.conn.URL
}
func (s *Server) GetShortURL() string {
	return s.conn.URL[8:]
}

// close the mock http server
func (s *Server) Close() {
	s.conn.Close()
}

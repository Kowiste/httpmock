package httpmock

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Test struct {
	Name  string
	Other int
}

func TestSomething(t *testing.T) {
	s := New()
	defer s.Close()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	data := Test{
		Name:  "john",
		Other: 40,
	}
	url:=s.GetURL()
	assert.Equal(t,url[8:],s.GetShortURL())

	client := &http.Client{Transport: tr}
	s.Expect(http.MethodGet, "/user").WillReturn(http.StatusBadRequest, data)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	resp, err = client.Get(url + "/user")
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	s.Expect(http.MethodGet, "/other").WillReturn(http.StatusOK, data).Once()
	resp, err = client.Get(url + "/other")
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp, err = client.Get(url + "/other")
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

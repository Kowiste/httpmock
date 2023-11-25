package httpmock

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
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
	url := s.GetURL()
	assert.Equal(t, url[8:], s.GetShortURL())

	client := &http.Client{Transport: tr}
	//Test expectation not found
	s.Expect(http.MethodGet, "/user").WillReturn(http.StatusBadRequest, data)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	//test expectation return
	resp, err = client.Get(url + "/user")
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	respByte, _ := io.ReadAll(resp.Body)
	respData := new(Test)
	json.Unmarshal(respByte, &respData)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, data, *respData)

	//test query
	resp, err = client.Get(url + "/user?id=abcdefghi")
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	json.Unmarshal(respByte, &respData)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	//Test Once method
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

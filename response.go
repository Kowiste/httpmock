package httpmock

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	status int
	data   any
}

// create the response for the request
func (resp *response) setResponse(w http.ResponseWriter) {
	b, err := json.Marshal(resp.data)
	if err != nil {
		fmt.Println("error marshal response ")
		return
	}
	// Write the JSON response with a specific status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.status) // You can replace http.StatusOK with any other HTTP status code
	w.Write(b)
}

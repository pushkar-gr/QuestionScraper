package hackerearth

import (
	"fmt"
	"io"
	"net/http"
)

const ENDPOINT = "https://www.hackerearth.com" //api endpoint

// send GET request to ENDPOINT
// input: request URL
// output: request body or error if any
func sendRequest(url string) ([]byte, error) {
	//send request
	resp, err := http.Get(ENDPOINT + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//return if any error in https response
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("received error status code: %d", resp.StatusCode)
	}

	//read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return respBody, nil
}

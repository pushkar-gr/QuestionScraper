package leetcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const ENDPOINT = "https://leetcode.com/graphql/" //api endpoint

type graphQLRequest struct {
	Query         string         `json:"query"`
	Variables     map[string]any `json:"variables"`
	OperationName string         `json:"operationName"`
}

type graphQLError struct {
	Message   string `json:"message"`
	Locations []struct {
		Line   int `json:"line"`
		Column int `json:"column"`
	} `json:"locations"`
	Path []string `json:"path"`
}

// send POST request to ENDPOINT
// input: GraphQLRequest
// output: response body or error if any
func sendRequest(body *graphQLRequest) ([]byte, error) {
	//encode requestBody to json format
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewBuffer(jsonBody)

	//create a new post request to api endpoint with json body
	req, err := http.NewRequest("POST", ENDPOINT, bodyReader)
	if err != nil {
		return nil, err
	}

	//set headers for request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://leetcode.com")
	req.Header.Set("x-csrftoken", "x-csrftoken")
	req.Header.Set("User-Agent", "go-script")

	client := &http.Client{}

	//send request
	resp, err := client.Do(req)
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

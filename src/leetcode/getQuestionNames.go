package leetcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const LIMIT = 100                                //number of question names to get in 1 request
const ENDPOINT = "https://leetcode.com/graphql/" //api endpoint

type GraphQLRequest struct {
	Query         string         `json:"query"`
	Variables     map[string]any `json:"variables"`
	OperationName string         `json:"operationName"`
}

type RespFormatSuccess struct {
	Errors []struct {
		Message   string `json:"message"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
		Path []string `json:"path"`
	} `json:"errors"`

	Data struct {
		ProblemsetQuestionList struct {
			Questions []struct {
				PaidOnly  bool   `json:"paidOnly"`
				TitleSlug string `json:"titleSlug"`
			} `json:"questions"`
		} `json:"problemsetQuestionList"`
	} `json:"data"`
}

// get LIMIT question names related to topic
// input: topic name, number of questions to skip from front, list of questions to append
// output: end of question list, error if any
func getQuestionNames(topic string, skip int, questions *[]string) (bool, error) {
	//generate a request body
	requestBody := GraphQLRequest{
		Query: `query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {
	           problemsetQuestionList: questionList(
							categorySlug: $categorySlug
							limit: $limit
							skip: $skip
							filters: $filters
	           ) {
							questions: data {
								paidOnly: isPaidOnly
								titleSlug
							}
	           }
	       }`,

		Variables: map[string]any{
			"categorySlug": "all-code-essentials",
			"limit":        LIMIT,
			"skip":         skip,
			"filters": map[string]any{
				"tags": []string{topic},
			},
		},

		OperationName: "problemsetQuestionList",
	}

	//encode requestBody to json format
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return false, err
	}

	bodyReader := bytes.NewBuffer(jsonBody)

	//create a new post request to api endpoint with json body
	req, err := http.NewRequest("POST", ENDPOINT, bodyReader)
	if err != nil {
		return false, err
	}

	//set headers for request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://leetcode.com")
	req.Header.Set("x-csrftoken", "x-csrftoken")
	req.Header.Set("User-Agent", "go-script")

	//create new client to send request
	client := &http.Client{}

	//send request
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	//return if any error in https response
	if resp.StatusCode >= 400 {
		return false, fmt.Errorf("received error status code: %d", resp.StatusCode)
	}

	//read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %v", err)
	}

	//structure to hold json
	var respStruct RespFormatSuccess

	//convert response to json
	err = json.Unmarshal(respBody, &respStruct)
	if err != nil {
		return false, err
	}

	//return if any error in response json
	if respStruct.Errors != nil {
		return false, fmt.Errorf("received error response: %v", string(respBody))
	}

	lenPrev := len(*questions)

	//iterate throught question names and append to list if question is not paidOnly
	for _, v := range respStruct.Data.ProblemsetQuestionList.Questions {
		if !v.PaidOnly {
			*questions = append(*questions, v.TitleSlug)
		}
	}

	return lenPrev == len(*questions), nil
}

// get question names related to given topics
// input: topic name
// output: list of question names and error if any
func GetQuestionNames(topic string) ([]string, error) {
	//create list to hold question names
	questions := make([]string, 0, LIMIT)

	//skip 0 questions
	skip := 0
	for {
		//get LIMIT question names
		completed, err := getQuestionNames(topic, skip, &questions)
		if err != nil {
			return nil, err
		}
		//break if reached to last question
		if completed {
			break
		}
		//add LIMIT to skip to skip already read questions
		skip += LIMIT
	}

	return questions, nil
}

package geeksforgeeks

import (
	"encoding/json"
	"fmt"
)

type respFormatTitleSlug struct {
	Error respError `json:"error"`

	Next    int             `json:"next"`
	Count   int             `json:"count"`
	Results []questionTitle `json:"results"`
}

type questionTitle struct {
	Slug       string `json:"slug"`
	ProblemUrl string `json:"problem_url"`
}

// get LIMIT question names related to topic
// input: topic name, number of questions to skip from front, list of questions to append
// output: end of question list, error if any
func getTitleSlugs(topic string, page int, questions *[]questionTitle) (int, error) {
	url := fmt.Sprintf("?pageMode=explore&page=%d&category=%s&sortBy=submissions", page, topic)
	respBody, err := sendRequest(url)
	if err != nil {
		return 0, err
	}

	//structure to hold json
	var respStruct respFormatTitleSlug

	//convert response to json
	err = json.Unmarshal(respBody, &respStruct)
	if err != nil {
		return 0, err
	}

	//return if any error in response json
	if respStruct.Error.Code != 0 {
		return 0, fmt.Errorf("received error response: %v", string(respBody))
	}

	//iterate throught question names and append to list if question is not paidOnly
	for _, v := range respStruct.Results {
		*questions = append(*questions, v)
	}

	return respStruct.Next, nil
}

// get question names related to given topics
// input: topic name
// output: list of question names and error if any
func GetTitleSlugs(topic string) ([]questionTitle, error) {
	//create list to hold question names
	questions := make([]questionTitle, 0, 20)

	page := 1

	for {
		page, err := getTitleSlugs(topic, page, &questions)
		if err != nil {
			return nil, err
		}
		//break if reached to last page
		if page == 0 {
			break
		}
	}

	return questions, nil
}

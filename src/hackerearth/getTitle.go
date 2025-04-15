package hackerearth

import (
	"encoding/json"
	"fmt"
)

const LIMIT = 100 //number of question names to get in 1 request

type respFormatTitleSlug struct {
	Problems struct {
		Algorithm []questionTitle `json:"algorithm"`
	} `json:"problems"`
	TotalProblemCount int `json:"total_problem_count"`
}

type questionTitle struct {
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
	Url        string `json:"url"`
}

// get LIMIT question names related to topic
// input: topic name, number of questions to skip from front, list of questions to append
// output: end of question list, error if any
func getTitles(topic string, offset int, questions *[]questionTitle) (bool, error) {
	url := fmt.Sprintf("/practice/api/problems/?limit=%d&offset=%d&tag=%s", LIMIT, offset, topic)
	respBody, err := sendRequest(url)
	if err != nil {
		return false, err
	}

	//structure to hold json
	var respStruct respFormatTitleSlug

	//convert response to json
	err = json.Unmarshal(respBody, &respStruct)
	if err != nil {
		return false, err
	}

	//return if any error in response json
	if respStruct.Problems.Algorithm == nil {
		return false, fmt.Errorf("received error response: %v", string(respBody))
	}

	//iterate throught question names and append to list if question is not paidOnly
	for _, v := range respStruct.Problems.Algorithm {
		*questions = append(*questions, v)
	}

	return respStruct.TotalProblemCount == len(*questions), nil
}

// get question names related to given topics
// input: topic name
// output: list of question names and error if any
func GetTitles(topic string) ([]questionTitle, error) {
	//create list to hold question names
	questions := make([]questionTitle, 0, 20)

	offset := 0

	for {
		page, err := getTitles(topic, offset, &questions)
		if err != nil {
			return nil, err
		}
		//break if reached to last page
		if page {
			break
		}

		offset += LIMIT
	}

	return questions, nil
}

package leetcode

import (
	"encoding/json"
	"fmt"
)

const LIMIT = 100 //number of question names to get in 1 request

type respFormatTitleSlug struct {
	Errors []GraphQLError `json:"errors"`

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
func getTitleSlugs(topic string, skip int, questions *[]string) (bool, error) {
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

	respBody, err := sendRequest(requestBody)
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
func GetTitleSlugs(topic string) ([]string, error) {
	//create list to hold question names
	questions := make([]string, 0, LIMIT)

	//skip 0 questions
	skip := 0
	for {
		//get LIMIT question names
		completed, err := getTitleSlugs(topic, skip, &questions)
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

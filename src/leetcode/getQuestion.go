package leetcode

import (
	"encoding/json"
	"fmt"

	"github.com/pushkar-gr/QuestionScraper/src/types"
)

type respFormatQuestion struct {
	Errors []graphQLError `json:"errors"`

	Data struct {
		Question struct {
			Title      string `json:"title"`
			QuestionId string `json:"questionId"`
			Difficulty string `json:"difficulty"`
			Content    string `json:"content"`
			TopicTags  []struct {
				Slug string `json:"slug"`
			} `json:"topicTags"`
			Solution struct {
				PaidOnly bool   `json:"paidOnly"`
				Content  string `json:"content"`
			} `json:"solution"`
		} `json:"question"`
	} `json:"data"`
}

func GetQuestion(titleSlug string) (*types.Question, error) {
	//get question details
	requestBody := graphQLRequest{
		Query: `query GetProblemData($titleSlug: String) {
							question(titleSlug: $titleSlug) {
								title
								questionId
								difficulty
								content
								solution {
									paidOnly
									content
								}
								topicTags {
									slug
								}
							}
						}`,

		Variables: map[string]any{
			"titleSlug": titleSlug,
		},

		OperationName: "GetProblemData",
	}

	respBody, err := sendRequest(&requestBody)
	if err != nil {
		return nil, err
	}

	//structure to hold json
	var respStruct respFormatQuestion

	//convert response to json
	err = json.Unmarshal(respBody, &respStruct)
	if err != nil {
		return nil, err
	}

	//return if any error in response json
	if respStruct.Errors != nil {
		return nil, fmt.Errorf("received error response: %v", string(respBody))
	}

	question := new(types.Question)

	//fill data
	question.Title = respStruct.Data.Question.Title
	question.Platform = "LeetCode"
	question.ExternalID = respStruct.Data.Question.QuestionId
	question.Link = fmt.Sprintf("www.leetcode.com/problems/%s/description", titleSlug)
	question.Difficulty = respStruct.Data.Question.Difficulty
	question.Question = respStruct.Data.Question.Content
	question.Solution = respStruct.Data.Question.Solution.Content

	for _, topic := range respStruct.Data.Question.TopicTags {
		question.Topics = append(question.Topics, topic.Slug)
	}

	return question, nil
}

package geeksforgeeks

import (
	"encoding/json"
	"fmt"

	"github.com/pushkar-gr/QuestionScraper/src/types"
)

type respFormatQuestion struct {
	Error respError `json:"error"`

	Results struct {
		ProblemName     string   `json:"problem_name"`
		Id              string   `json:"id"`
		Difficulty      string   `json:"difficulty"`
		ProblemQuestion string   `json:"problem_question"`
		TopicTags       []string `json:"topic_tags"`
	} `json:"results"`
}

func GetQuestion(questionTitle questionTitle) (*types.Question, error) {
	respBody, err := sendRequest(questionTitle.Slug)
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
	if respStruct.Error.Code != 0 {
		return nil, fmt.Errorf("received error response: %v", string(respBody))
	}

	question := new(types.Question)

	//fill data
	question.Title = respStruct.Results.ProblemName
	question.Platform = "GeeksForGeeks"
	question.ExternalID = respStruct.Results.Id
	question.Link = questionTitle.ProblemUrl
	question.Difficulty = respStruct.Results.Difficulty
	question.Question = respStruct.Results.ProblemQuestion

	for _, topic := range respStruct.Results.TopicTags {
		question.Topics = append(question.Topics, topic)
	}

	return question, nil

}

package geeksforgeeks

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pushkar-gr/QuestionScraper/src/types"
)

type respFormatQuestion struct {
	Error respError `json:"error"`

	Results struct {
		ProblemName     string   `json:"problem_name"`
		Id              int      `json:"id"`
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

	var difficulty types.DifficultyLevel
	respStruct.Results.Difficulty = strings.ToLower(respStruct.Results.Difficulty)
	if respStruct.Results.Difficulty == "medium" {
		difficulty = types.Medium
	} else if respStruct.Results.Difficulty == "hard" {
		difficulty = types.Hard
	} else {
		difficulty = types.Easy
	}

	//fill data
	question.Title = respStruct.Results.ProblemName
	question.Platform = "GeeksForGeeks"
	question.ExternalID = fmt.Sprintf("%d", respStruct.Results.Id)
	question.Link = questionTitle.ProblemUrl
	question.Difficulty = difficulty
	if question.Difficulty == "Basic" {
		question.Difficulty = "Easy"
	}
	question.Question = respStruct.Results.ProblemQuestion

	for _, topic := range respStruct.Results.TopicTags {
		question.Topics = append(question.Topics, topic)
	}

	return question, nil

}

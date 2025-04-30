package hackerearth

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/pushkar-gr/QuestionScraper/src/types"
)

type respFormatQuestion struct {
	Id                int    `json:"id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	SampleExplanation string `json:"sample_explanation"`
	Tags              string `json:"tags"`
	Editorial         struct {
		State string `json:"state"`
	} `json:"editorial"`
}

func GetQuestion(questionTitle questionTitle) (*types.Question, error) {
	respBody, err := sendRequest(questionTitle.Url)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`problemData:\s*({(?:[^{}]|{(?:[^{}]|{[^{}]*})*})*})`)
	matches := re.FindStringSubmatch(string(respBody))
	if len(matches) < 2 {
		return nil, fmt.Errorf("Could not extract problemData from the HTML")
	}

	problemDataStr := matches[1]

	//structure to hold json
	var respStruct respFormatQuestion

	//convert response to json
	err = json.Unmarshal([]byte(problemDataStr), &respStruct)
	if err != nil {
		return nil, err
	}

	// return if any error in response json
	if respStruct.Id == 0 {
		return nil, fmt.Errorf("received error response: %v", string(respBody))
	}

	question := new(types.Question)

	var difficulty types.DifficultyLevel
	questionTitle.Difficulty = strings.ToLower(questionTitle.Difficulty)
	if questionTitle.Difficulty == "medium" {
		difficulty = types.Medium
	} else if questionTitle.Difficulty == "hard" {
		difficulty = types.Hard
	} else {
		difficulty = types.Easy
	}

	//fill data
	question.Title = respStruct.Title
	question.Platform = "HackerEarth"
	question.ExternalID = fmt.Sprintf("%d", respStruct.Id)
	question.Link = "https://www.hackerearth.com" + questionTitle.Url
	question.Difficulty = difficulty
	question.Question = respStruct.Description
	question.Explanation = respStruct.SampleExplanation

	for topic := range strings.SplitSeq(respStruct.Tags, ",") {
		question.Topics = append(question.Topics, topic)
	}

	return question, nil
}

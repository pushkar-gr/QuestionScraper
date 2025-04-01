// src/hacker_rank/service/question_service.go
package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"src/hacker_rank/config"
	"src/hacker_rank/util"
)

type Question struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
	Body       string `json:"body"`
}

type QuestionService struct {
	config     *config.Config
	httpClient *util.HTTPClient
}

func NewQuestionService(config *config.Config, httpClient *util.HTTPClient) *QuestionService {
	return &QuestionService{config: config, httpClient: httpClient}
}

func (s *QuestionService) GetQuestions(arrayID int) ([]Question, error) {
	url := fmt.Sprintf("https://api.hackerrank.com/api/v3/challenges/arrays/%d/questions", arrayID)

	req, err := s.httpClient.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("x-api-key", s.config.HackerRankAPIKey)

	resp, body, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HackerRank API returned status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Data []Question `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return response.Data, nil
}

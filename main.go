package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pushkar-gr/QuestionScraper/src/hacker_rank/config"
	"github.com/pushkar-gr/QuestionScraper/src/hacker_rank/handler"
	"github.com/pushkar-gr/QuestionScraper/src/hacker_rank/service"
	"github.com/pushkar-gr/QuestionScraper/src/hacker_rank/util"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	httpClient := util.NewHTTPClient()
	questionService := service.NewQuestionService(cfg, httpClient)
	questionHandler := handler.NewQuestionHandler(questionService)

	router := gin.Default()
	router.GET("/getquestions/:arrayID", questionHandler.GetQuestions)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

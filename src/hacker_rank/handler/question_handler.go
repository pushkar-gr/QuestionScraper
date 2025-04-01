// src/hacker_rank/handler/question_handler.go
package handler

import (
	"net/http"
	"strconv"

	"src/hacker_rank/service"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	service *service.QuestionService
}

func NewQuestionHandler(service *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{service: service}
}

func (h *QuestionHandler) GetQuestions(c *gin.Context) {
	arrayIDStr := c.Param("arrayID")
	arrayID, err := strconv.Atoi(arrayIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid array ID"})
		return
	}

	questions, err := h.service.GetQuestions(arrayID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

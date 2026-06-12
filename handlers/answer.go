package handlers

import (
	"net/http"
	"quiz/entities/dto"
	"quiz/services"

	"github.com/gin-gonic/gin"
)

type AnswerHandler struct {
	AnswerService services.AnswerServiceInterface
}

func NewAnswerHandler(service services.AnswerServiceInterface) *AnswerHandler {
	return &AnswerHandler{
		AnswerService: service,
	}
}

func (h *AnswerHandler) CheckAnswer(c *gin.Context) {
	ctx := c.Request.Context()

	var body dto.CheckAnswerDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}

	questionID := c.Param("question_id")

	correct, err := h.AnswerService.CheckAnswer(ctx, questionID, body.AnswerID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"correct": correct,
	})
}

func (h *AnswerHandler) ListAnswers(c *gin.Context) {
	ctx := c.Request.Context()
	questionID := c.Param("question_id")

	answers, err := h.AnswerService.ListAnswers(ctx, questionID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, answers)
}

func (h *AnswerHandler) CreateAnswer(c *gin.Context) {
	ctx := c.Request.Context()
	questionID := c.Param("question_id")

	var body dto.CreateAnswerDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json format",
		})
		return
	}

	userID := c.MustGet("userID").(int)

	newAnswer, err := h.AnswerService.CreateAnswer(ctx, questionID, body, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, newAnswer)
}

func (h *AnswerHandler) GetAnswer(c *gin.Context) {
	ctx := c.Request.Context()
	answerID := c.Param("answer_id")

	answer, err := h.AnswerService.GetAnswer(ctx, answerID)

	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, answer)
}

func (h *AnswerHandler) DeleteAnswer(c *gin.Context) {
	ctx := c.Request.Context()
	answerID := c.Param("answer_id")

	userID := c.MustGet("userID").(int)

	err := h.AnswerService.DeleteAnswer(ctx, answerID, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AnswerHandler) UpdateAnswer(c *gin.Context) {
	ctx := c.Request.Context()
	answerID := c.Param("answer_id")

	var body dto.UpdateAnswerDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json format",
		})
		return
	}

	userID := c.MustGet("userID").(int)

	updatedAnswer, err := h.AnswerService.UpdateAnswer(ctx, answerID, body, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedAnswer)
}

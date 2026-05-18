package handlers

import (
	"net/http"
	"quiz/entities/dto"
	"quiz/services"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	QuestionService *services.QuestionService
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
		ctx := c.Request.Context()
		quizID := c.Param("quiz_id")

		var body dto.CreateQuestionDTO
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json",
			})
			return
		}
		newQuestion, err := h.QuestionService.CreateQuestion(ctx, quizID, body)
		if err != nil {
			HandleError(c, err)
			return
		}

		c.JSON(http.StatusCreated, newQuestion)
}

func (h *QuestionHandler) ListQuestions(c *gin.Context){
		ctx := c.Request.Context()
		id := c.Param("quiz_id")
		questions, err := h.QuestionService.ListQuestions(ctx, id)
		if err != nil{
			HandleError(c, err)
			return;
		}
		c.JSON(http.StatusOK, questions);
	}

func (h *QuestionHandler) GetQuestion(c *gin.Context){
		ctx := c.Request.Context()
		questionID := c.Param("question_id");

		question, err := h.QuestionService.GetQuestion(ctx, questionID);
		if err != nil{
			HandleError(c, err)
			return;
		}
		c.JSON(http.StatusOK, question);
	}

func (h *QuestionHandler) UpdateQuestion(c *gin.Context){
		ctx := c.Request.Context()
		questionID := c.Param("question_id");

		var body dto.UpdateQuestionDTO;
		if err := c.ShouldBindJSON(&body); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json",
			})
			return;
		}

		question, err := h.QuestionService.UpdateQuestion(ctx, questionID, body)
		if err != nil{
			HandleError(c, err)
			return;
		}

		c.JSON(http.StatusOK, question);
	}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context){
		ctx := c.Request.Context()
		questionID := c.Param("question_id");

		err := h.QuestionService.DeleteQuestion(ctx, questionID)
		if err != nil{
			HandleError(c, err)
			return;
		}

		c.Status(http.StatusNoContent);
	}
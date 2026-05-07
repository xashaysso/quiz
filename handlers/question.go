package handlers

import (
	"net/http"
	"quiz/db/repositories"
	"quiz/entities/dto"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	Repo *repositories.PgQuestionRepo
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

		if len(body.Answers) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no 'answers' provided in json",
			})
			return
		}

		qID, err := strconv.Atoi(quizID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid 'ID' field format",
			})
			return
		}

		newQuestion, err := h.Repo.CreateQuestion(ctx, qID, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, newQuestion)
	}

func (h *QuestionHandler) ListQuestions(c *gin.Context){
		ctx := c.Request.Context()
		id := c.Param("quiz_id")
		questions, err := h.Repo.GetQuizQuestions(ctx, id);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()});
			return;
		}
		c.JSON(http.StatusOK, questions);
	}

func (h *QuestionHandler) GetQuestion(c *gin.Context){
		ctx := c.Request.Context()
		questionID := c.Param("question_id");

		question, err := h.Repo.GetQuestion(ctx, questionID);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
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

		if body.Text == nil && body.NewCorrectID == nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "fields 'text' and/or 'correct_id' are required in json body",
			})
			return;
		}

		question, err := h.Repo.UpdateQuestion(ctx, questionID, body);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return;
		}

		c.JSON(http.StatusOK, question);
	}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context){
		ctx := c.Request.Context()
		questionID := c.Param("question_id");

		err := h.Repo.DeleteQuestion(ctx, questionID);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return;
		}

		c.Status(http.StatusNoContent);
	}
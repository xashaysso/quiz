package handlers

import (
	"net/http"
	"quiz/db/repositories"
	"quiz/entities/dto"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	Repo repositories.QuizRepository
}

func (h *QuizHandler) ListQuizzes(c *gin.Context) {
	ctx := c.Request.Context()
	quizzes, err := h.Repo.GetQuiz(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}
	c.JSON(http.StatusOK, quizzes);
}

func (h *QuizHandler) DeleteQuiz(c *gin.Context){
	ctx := c.Request.Context()
	quizID := c.Param("quiz_id")
	userID := c.MustGet("userID").(int)

	err := h.Repo.DeleteQuiz(ctx, quizID, userID)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return;
	}
	c.Status(http.StatusNoContent);
}

func (h *QuizHandler) CreateQuiz(c *gin.Context){
	ctx := c.Request.Context()
	var body dto.CreateQuizDTO;

	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return;
	}
	if body.Name == nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "field 'name' is required in JSON body",
		})
		return;
	}
	
	// middleware guarantees userID of type int there
	userID := c.MustGet("userID").(int)

	newQuiz, err := h.Repo.CreateQuiz(ctx, *body.Name, body.Description, userID);
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return;
	}
	c.JSON(http.StatusCreated, newQuiz);
}

func (h *QuizHandler) UpdateQuiz(c *gin.Context){
	ctx := c.Request.Context()
	quizID := c.Param("quiz_id");
	userID := c.MustGet("userID").(int)

	var body dto.UpdateQuizDTO;

	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return;
	}
	if body.Name == nil && body.Description == nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"field 'name' and/or 'description' required in json body",
		})
		return;
	}

	newQuiz, err := h.Repo.UpdateQuiz(ctx, quizID, body.Name, body.Description, userID);
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return;
	}

	c.JSON(http.StatusOK, newQuiz);
}
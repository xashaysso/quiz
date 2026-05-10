package handlers

import (
	"errors"
	"net/http"
	"quiz/entities/dto"
	"quiz/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	QuizService *services.QuizService
}

func (h *QuizHandler) ListQuizzes(c *gin.Context) {
	ctx := c.Request.Context()
	quizzes, err := h.QuizService.ListQuizzes(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return;
	}
	c.JSON(http.StatusOK, quizzes);
}

func (h *QuizHandler) DeleteQuiz(c *gin.Context){
	ctx := c.Request.Context()
	idStr := c.Param("quiz_id")

	quizID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid quiz id format",
		})
		return
	}
	userID := c.MustGet("userID").(int)

	err = h.QuizService.DeleteQuiz(ctx, quizID, userID)
	if err != nil{
		if errors.Is(err, services.ErrQuizNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		} else if errors.Is(err, services.ErrNotAnAuthor) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
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
	
	// middleware guarantees userID of type int there
	userID := c.MustGet("userID").(int)

	newQuiz, err := h.QuizService.CreateQuiz(ctx, body.Name, body.Description, userID)
	if err != nil{
		if errors.Is(err, services.ErrInvalidName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return;
	}
	c.JSON(http.StatusCreated, newQuiz);
}

func (h *QuizHandler) UpdateQuiz(c *gin.Context){
	ctx := c.Request.Context()
	idStr := c.Param("quiz_id");
	quizID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid quiz id format",
		})
		return
	}
	userID := c.MustGet("userID").(int)

	var body dto.UpdateQuizDTO;

	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return;
	}

	newQuiz, err := h.QuizService.UpdateQuiz(ctx, quizID, body.Name, body.Description, userID);
	if err != nil{
		if errors.Is(err, services.ErrQuizNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		} else if errors.Is(err, services.ErrNotAnAuthor) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return;
	}

	c.JSON(http.StatusOK, newQuiz);
}
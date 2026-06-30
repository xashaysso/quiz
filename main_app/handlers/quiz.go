package handlers

import (
	"net/http"
	"quiz/entities/dto"
	"quiz/services"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	QuizService services.QuizServiceInterface
}

func NewQuizHandler(service services.QuizServiceInterface) *QuizHandler {
	return &QuizHandler{
		QuizService: service,
	}
}

func (h *QuizHandler) ListQuizzes(c *gin.Context) {
	ctx := c.Request.Context()
	quizzes, err := h.QuizService.ListQuizzes(ctx)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, quizzes)
}

func (h *QuizHandler) DeleteQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	quizID := c.Param("quiz_id")
	userID := c.MustGet("userID").(int)

	err := h.QuizService.DeleteQuiz(ctx, quizID, userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	var body dto.CreateQuizDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}

	// middleware guarantees userID of type int there
	userID := c.MustGet("userID").(int)

	newQuiz, err := h.QuizService.CreateQuiz(ctx, body.Name, body.Description, userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newQuiz)
}

func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	quizID := c.Param("quiz_id")

	userID := c.MustGet("userID").(int)

	var body dto.UpdateQuizDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}

	newQuiz, err := h.QuizService.UpdateQuiz(ctx, quizID, body.Name, body.Description, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, newQuiz)
}

func (h *QuizHandler) StartQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	quizID := c.Param("quiz_id")

	userIDInt := c.MustGet("userID").(int64)
	userID := int64(userIDInt)

	session, err := h.QuizService.StartQuiz(ctx, userID, quizID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.SetCookie("quiz_session", session, 3600, "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"quiz_session": session,
	})
}

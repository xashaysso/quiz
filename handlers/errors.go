package handlers

import (
	"errors"
	"net/http"
	"quiz/services"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	switch {
	// 400 Bad Request
	case errors.Is(err, services.ErrInvalidIDFormat),
		errors.Is(err, services.ErrInvalidPassword),
		errors.Is(err, services.ErrInvalidUsername),
		errors.Is(err, services.ErrNoQuestionAnswers),
		errors.Is(err, services.ErrWrongCredentials),
		errors.Is(err, services.ErrInvalidName),
		errors.Is(err, services.ErrInvalidAnswerText),
		errors.Is(err, services.ErrNoFieldsToUpdate),
		errors.Is(err, services.ErrInvalidCorrectID),
		errors.Is(err, services.ErrNoRequiredFields):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	// 403 Not Found
	case errors.Is(err, services.ErrQuizNotFound),
		errors.Is(err, services.ErrQuestionNotFound),
		errors.Is(err, services.ErrAnswerNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	// 404 Forbidden
	case errors.Is(err, services.ErrNotAnAuthor):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

	// 409 Conflict
	case errors.Is(err, services.ErrUserAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	// 500 Internal Server Err
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
}

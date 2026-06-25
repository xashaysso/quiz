package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"quiz/services"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	logAttrs := []any{
		slog.String("path", c.Request.URL.Path),
		slog.String("method", c.Request.Method),
		slog.String("err_message", err.Error()),
	}

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
		errors.Is(err, services.ErrQuestionAlreadyAnswered),
		errors.Is(err, services.ErrQuestionDoesNotBelongToQuiz),
		errors.Is(err, services.ErrNoRequiredFields):
		slog.Warn("client bad request", logAttrs...)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	// 401 Unauthorized
	case errors.Is(err, services.ErrWrongCredentials),
		errors.Is(err, services.ErrSessionExpired):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

	// 403 Not Found
	case errors.Is(err, services.ErrQuizNotFound),
		errors.Is(err, services.ErrQuestionNotFound),
		errors.Is(err, services.ErrAnswerNotFound):
		slog.Warn("client request failed", logAttrs...)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	// 404 Forbidden
	case errors.Is(err, services.ErrNotAnAuthor):
		slog.Warn("client request failed", logAttrs...)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

	// 409 Conflict
	case errors.Is(err, services.ErrUserAlreadyExists):
		slog.Warn("client request failed", logAttrs...)
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	// 500 Internal Server Err
	default:
		slog.Warn("internal server error detected", logAttrs...)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
}

package middleware

import (
	"net/http"
	"quiz/db/repositories"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(repo repositories.SessionRepository) gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		token, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		userID, err := repo.Get(ctx, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid session",
			})
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
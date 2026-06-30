package middleware

import (
	"net/http"
	"quiz/pkg/authv1"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

func AuthMiddleware(authClient authv1.AuthServiceClient, localCache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		token, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		if cachedUserID, found := localCache.Get(token); found {
			c.Set("userID", cachedUserID.(int64))
			c.Next()
			return
		}

		res, err := authClient.CheckSession(ctx, &authv1.CheckSessionRequest{
			SessionId: token,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid session",
			})
			return
		}

		userID := res.GetUserId()
		localCache.Set(token, userID, 1*time.Minute)

		c.Set("userID", userID)
		c.Next()
	}
}

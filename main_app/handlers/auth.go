package handlers

import (
	"net/http"
	"quiz/entities/dto"
	"quiz/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService services.AuthServiceInterface
}

func NewAuthHandler(service services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		AuthService: service,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var body dto.AuthRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}

	user, token, err := h.AuthService.Register(ctx, body.Username, body.Password)
	if err != nil {
		HandleError(c, err)
		return
	}

	h.setAuthCookie(c, token)

	newUser := dto.RegisterResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(http.StatusCreated, gin.H{
		"user":  newUser,
		"token": token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var body dto.AuthRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}

	token, err := h.AuthService.Login(ctx, body.Username, body.Password)
	if err != nil {
		HandleError(c, err)
		return
	}

	h.setAuthCookie(c, token)

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "successfully logged in",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "already logged out",
		})
		return
	}
	err = h.AuthService.Logout(ctx, token)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.SetCookie("token", token, -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "logged out",
	})
}

// helper method
func (h *AuthHandler) setAuthCookie(c *gin.Context, token string) {
	c.SetCookie("token", token, 3600*24, "/", "", false, true)
}

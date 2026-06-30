package handlers

import (
	"net/http"
	"quiz/entities/dto"
	"quiz/pkg/authv1"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authClient authv1.AuthServiceClient
}

func NewAuthHandler(client authv1.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		authClient: client,
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

	res, err := h.authClient.Register(ctx, &authv1.RegisterRequest{
		Username: body.Username,
		Password: body.Password,
	})
	if err != nil {
		HandleGrpcError(c, err)
		return
	}

	h.setAuthCookie(c, res.GetToken())

	newUser := dto.RegisterResponse{
		ID:        int(res.GetUserId()),
		Username:  body.Username,
		CreatedAt: time.Now(),
	}
	c.JSON(http.StatusCreated, gin.H{
		"user":  newUser,
		"token": res.GetToken(),
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

	res, err := h.authClient.Login(ctx, &authv1.LoginRequest{
		Username: body.Username,
		Password: body.Password,
	})
	if err != nil {
		HandleGrpcError(c, err)
		return
	}

	h.setAuthCookie(c, res.GetToken())

	c.JSON(http.StatusOK, gin.H{
		"token":   res.GetToken(),
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
	_, err = h.authClient.DeleteSession(ctx, &authv1.DeleteSessionRequest{
		SessionId: token,
	})
	if err != nil {
		HandleGrpcError(c, err)
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

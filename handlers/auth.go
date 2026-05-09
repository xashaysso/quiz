package handlers

import (
	"net/http"
	"quiz/db/repositories"
	"quiz/entities/dto"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Repo repositories.UserRepository
	SessionRepo repositories.SessionRepository
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
	if len(body.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "minimum username length is 3",
		})
		return
	}
	if len(body.Password) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "minimum password length is 5",
		})
		return
	}
	password := []byte(body.Password)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "couldn't hash the password",
		})
		return
	}

	user, err := h.Repo.CreateUser(ctx, body.Username, string(hash))
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
            if pgErr.Code == "23505" { 
                c.JSON(http.StatusBadRequest, gin.H{
					"error": "username already taken",
				})
				return
            }
        }
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	res := dto.RegisterResponse{
		ID: user.ID,
		Username: user.Username,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(http.StatusCreated, res)
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
	user, err := h.Repo.GetByUsername(ctx, body.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid username or password",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password) )
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid username or password",
		})
		return
	}

	token := uuid.New().String()
	ttl := 24 * time.Hour
	err = h.SessionRepo.Set(ctx, token, user.ID, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user session",
		})
		return
	}

	c.SetCookie("token", token, int(ttl.Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
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
	err = h.SessionRepo.Delete(ctx, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete session",
		})
	}
	c.SetCookie("token", token, -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "logged out",
	})
}
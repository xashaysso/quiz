package dto

import "time"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
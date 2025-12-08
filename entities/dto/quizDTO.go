package dto

type CreateQuizDTO struct {
	Name        *string `json:"name"`
	Description string  `json:"description"`
}

type UpdateQuizDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
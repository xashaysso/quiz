package dto

type CreateQuizDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateQuizDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type QuizResponse struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Questions   []QuestionResponse `json:"questions"`
}

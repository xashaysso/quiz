package entities

type Quiz struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatorID   int     `json:"creator_id"`
}

type Question struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	QuizID int    `json:"quiz_id"`
}

type Answer struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}
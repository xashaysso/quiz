package entities

type Quiz struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Question struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	QuizID    int    `json:"quiz_id"`
	CorrectID int    `json:"correct_id"`
}

type Answer struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	Text       string `json:"text"`
}
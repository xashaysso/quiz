package entities

type QuizAPI struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Questions   []QuestionAPI `json:"questions"`
}

type QuestionAPI struct {
	ID        int         `json:"id"`
	Text      string      `json:"text"`
	CorrectID int         `json:"correct_id"`
	Answers   []AnswerAPI `json:"answers"`
}

type AnswerAPI struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}
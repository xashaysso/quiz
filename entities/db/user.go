package entities

type User struct {
	ID        int            `json:"id"`
	Username  string         `json:"username"`
	QuizScore map[string]int `json:"quiz_score"`
}
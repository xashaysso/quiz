package entities

type UserAPI struct {
	ID        int
	Username  string
	QuizScore map[int]int
}
package entities

type Quiz struct {
	ID        int
	Name      string
	Questions []Question
}

type Question struct {
	ID        int
	Text      string
	CorrectID int
	Answers   []Answer
}

type Answer struct {
	ID   int
	Text string
}
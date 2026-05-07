package main

import (
	"log"
	"os"

	"quiz/db"
	"quiz/db/repositories"
	"quiz/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	PORT := os.Getenv("PORT")

	globalPool := db.Serve()
	defer globalPool.Close()

	quizRepo := repositories.NewQuizRepo(globalPool)
	questionRepo := repositories.NewQuestionRepo(globalPool)
	answerRepo := repositories.NewAnswerRepo(globalPool)

	quizH := &handlers.QuizHandler{Repo: quizRepo}
	questionH := &handlers.QuestionHandler{Repo: questionRepo}
	answerH := &handlers.AnswerHandler{Repo: answerRepo}

	quiz := router.Group("/quizzes")
	{
		// quiz handlers
		quiz.GET("/", quizH.ListQuizzes)
		quiz.POST("/", quizH.CreateQuiz)
		quiz.PATCH("/:quiz_id", quizH.UpdateQuiz)
		quiz.DELETE("/:quiz_id", quizH.DeleteQuiz)

		// question handlers
		quiz.GET("/:quiz_id/questions", questionH.ListQuestions)
		quiz.POST("/:quiz_id/questions", questionH.CreateQuestion)
	}

	question := router.Group("/questions")
	{
		// question handlers
		question.GET("/:question_id", questionH.GetQuestion)
		question.PATCH("/:question_id", questionH.UpdateQuestion)
		question.DELETE("/:question_id", questionH.DeleteQuestion)

		// answer handlers
		question.GET("/:question_id/answers", answerH.ListAnswers)
		question.POST("/:question_id/answers", answerH.CreateAnswer)

		question.POST("/:question_id/check", answerH.CheckAnswer)
	}

	answer := router.Group("/answers")
	{
		answer.GET("/:answer_id", answerH.GetAnswer)
		answer.PATCH("/:answer_id", answerH.UpdateAnswer)
		answer.DELETE("/:answer_id", answerH.DeleteAnswer)
	}

	router.Run(PORT)
}


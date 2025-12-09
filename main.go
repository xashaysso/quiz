package main

import (
	"context"
	"log"
	"os"
	"quiz/db"
	"quiz/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	router := gin.Default();

	err := godotenv.Load()
	if err != nil{
		log.Println("No .env file found");
	}
	PORT := os.Getenv("PORT");

	conn := db.Serve();
	defer conn.Close(context.Background());

	quiz := router.Group("/quizzes")
	{
		quiz.GET("/", handlers.ListQuizzes(conn));
		quiz.POST("/", handlers.CreateQuiz(conn));
		quiz.PATCH("/:quiz_id", handlers.UpdateQuiz(conn));
		quiz.DELETE("/:quiz_id", handlers.DeleteQuiz(conn))

		quiz.GET("/:quiz_id/questions", handlers.ListQuestions(conn));
		quiz.POST("/:quiz_id/questions", handlers.CreateQuestion(conn));
	}

	question := router.Group("/questions")
	{
		question.GET("/:question_id", handlers.GetQuestion(conn));
		question.PATCH("/:question_id", handlers.UpdateQuestion(conn));
		question.DELETE("/:question_id", handlers.DeleteQuestion(conn));

		question.GET("/:question_id/answers", handlers.ListAnswers(conn));
		question.POST("/:question_id/answers", handlers.CreateAnswer(conn));

		question.POST("/:question_id/check", handlers.CheckAnswer(conn));
	}

	answer := router.Group("/answers")
	{
		answer.GET("/:answer_id", handlers.GetAnswer(conn));
	}

	router.Run(PORT);
}
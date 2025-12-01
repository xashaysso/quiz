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

		quiz.GET("/:quiz_id", handlers.ListQuestions(conn));
	}

	question := router.Group("/questions")
	{
		question.GET("/:question_id/answers", handlers.ListAnswers(conn));
		question.POST("/:question_id/check", handlers.CheckAnswer(conn));
	}

	_, _ = quiz, question // ignore

	router.Run(PORT);
}
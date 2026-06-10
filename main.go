package main

import (
	"log"
	"os"

	"quiz/db"
	"quiz/db/repositories/pg"
	"quiz/db/repositories/redis"
	"quiz/handlers"
	"quiz/middleware"
	"quiz/services"

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
	REDIS_ADDR := os.Getenv("REDIS_ADDR")


	// postgre init
	globalPool := db.Serve()
	defer globalPool.Close()

	// repositories
	quizRepo := pg.NewQuizRepo(globalPool)
	questionRepo := pg.NewQuestionRepo(globalPool)
	answerRepo := pg.NewAnswerRepo(globalPool)
	userRepo := pg.NewUserRepo(globalPool)
	txManager := pg.NewPgTxManager(globalPool)

	// services
	quizService := &services.QuizService{QuizRepo: quizRepo}
	questionService := &services.QuestionService{QuestionRepo: questionRepo, AnswerRepo: answerRepo, TxManager: txManager}
	answerService := &services.AnswerService{AnswerRepo: answerRepo}

	// handlers
	quizH := &handlers.QuizHandler{QuizService: quizService}
	questionH := &handlers.QuestionHandler{QuestionService: questionService}
	answerH := &handlers.AnswerHandler{AnswerService: answerService}

	// redis init
	rdb := db.NewRedisClient(REDIS_ADDR);

	sessionRepo := redis.NewSessionRepository(rdb)

	authService := &services.AuthService{
		UserRepo: userRepo,
		SessionRepo: sessionRepo,
	}

	authH := &handlers.AuthHandler{AuthService: authService}
	// routes
	auth := router.Group("/auth")
	{
		// auth handlers
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/logout", authH.Logout)
	}

	quiz := router.Group("/quizzes")
	quiz.Use(middleware.AuthMiddleware(sessionRepo))
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
	question.Use(middleware.AuthMiddleware(sessionRepo))
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
	answer.Use(middleware.AuthMiddleware(sessionRepo))
	{
		// answer handlers
		answer.GET("/:answer_id", answerH.GetAnswer)
		answer.PATCH("/:answer_id", answerH.UpdateAnswer)
		answer.DELETE("/:answer_id", answerH.DeleteAnswer)
	}

	router.Run(PORT)
}


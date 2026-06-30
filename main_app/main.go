package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"quiz/db"
	kafka_producers "quiz/db/repositories/kafka"
	"quiz/db/repositories/pg"
	"quiz/db/repositories/redis"
	"quiz/handlers"
	"quiz/middleware"
	"quiz/pkg/authv1"
	"quiz/pkg/kafka"
	"quiz/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	router := gin.Default()

	err := godotenv.Load()
	if err != nil {
		slog.Error("No .env file found")
	}

	// logger
	var logger *slog.Logger

	if os.Getenv("APP_ENV") == "prod" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	slog.SetDefault(logger)

	slog.Info("Logger service has started", slog.String("env", os.Getenv("APP_ENV")))

	PORT := os.Getenv("PORT")
	REDIS_ADDR := os.Getenv("REDIS_ADDR")

	// postgre init
	globalPool := db.Serve()
	defer globalPool.Close()

	// kafka producer
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "127:0.0.1:9092"
	}
	brokers := strings.Split(kafkaBrokers, ",")

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "quiz-results"
	}

	kafkaProducer := kafka.NewProducer(brokers, kafkaTopic)
	defer kafkaProducer.Close()

	quizEventProducer := kafka_producers.NewQuizKafkaProducer(kafkaProducer)

	// repositories
	quizRepo := pg.NewQuizRepo(globalPool)
	questionRepo := pg.NewQuestionRepo(globalPool)
	answerRepo := pg.NewAnswerRepo(globalPool)
	txManager := pg.NewPgTxManager(globalPool)

	sessionCache := cache.New(1*time.Minute, 5*time.Minute)

	authServiceAddr := os.Getenv("AUTH_SERVICE_ADDR")

	conn, err := grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("cannot connect to auth service", slog.Any("err", err))
		os.Exit(1)
	}
	defer conn.Close()

	authClient := authv1.NewAuthServiceClient(conn)

	// redis init
	rdb := db.NewRedisClient(REDIS_ADDR)

	sessionRepo := redis.NewSessionRepository(rdb)

	// services
	quizService := services.NewQuizService(quizRepo, questionRepo, sessionRepo)
	questionService := services.NewQuestionService(questionRepo, answerRepo, txManager)
	answerService := services.NewAnswerService(answerRepo, quizRepo, sessionRepo, txManager, quizEventProducer)

	// handlers
	quizH := handlers.NewQuizHandler(quizService)
	questionH := handlers.NewQuestionHandler(questionService)
	answerH := handlers.NewAnswerHandler(answerService)

	authH := handlers.NewAuthHandler(authClient)

	authMiddleware := middleware.AuthMiddleware(authClient, sessionCache)
	// routes
	auth := router.Group("/auth")
	{
		// auth handlers
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/logout", authH.Logout)
	}

	quiz := router.Group("/quizzes")
	quiz.Use(authMiddleware)
	{
		// quiz handlers
		quiz.GET("/", quizH.ListQuizzes)
		quiz.POST("/", quizH.CreateQuiz)
		quiz.PATCH("/:quiz_id", quizH.UpdateQuiz)
		quiz.DELETE("/:quiz_id", quizH.DeleteQuiz)
		quiz.POST("/:quiz_id/start", quizH.StartQuiz)

		// question handlers
		quiz.GET("/:quiz_id/questions", questionH.ListQuestions)
		quiz.POST("/:quiz_id/questions", questionH.CreateQuestion)
	}

	question := router.Group("/questions")
	question.Use(authMiddleware)
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
	answer.Use(authMiddleware)
	{
		// answer handlers
		answer.GET("/:answer_id", answerH.GetAnswer)
		answer.PATCH("/:answer_id", answerH.UpdateAnswer)
		answer.DELETE("/:answer_id", answerH.DeleteAnswer)
	}

	// server start
	srv := &http.Server{
		Addr:    PORT,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen and serve failed", slog.Any("err", err))
			os.Exit(1)
		}
	}()
	slog.Info("server started succesfully", slog.String("port", PORT))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("shutting down server...")

	shutdownTime := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.Any("err", err))
		os.Exit(1)
	}
	slog.Info("server exited cleanly")
}

package handlers

import (
	"net/http"
	"quiz/db/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func ListQuizzes(conn *pgx.Conn) gin.HandlerFunc{
	return func (c *gin.Context){
		quizzes, err := repositories.GetQuiz(conn);
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return;
		}
		c.JSON(http.StatusOK, quizzes);
	}
}

func ListQuestions(conn *pgx.Conn) gin.HandlerFunc{
	return func (c *gin.Context){
		questions, err := repositories.GetQuizQuestions(c, conn);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()});
			return;
		}
		c.JSON(http.StatusOK, questions);
	}
}

func ListAnswers(conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){
		question_id := c.Param("question_id")

		answers, err := repositories.GetQuizAnswers(conn, question_id);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()});
			return;
		}
		c.JSON(http.StatusOK, answers);
	}
}

func CheckAnswer(conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){

		type RequestBody struct {
			AnswerID *int `json:"answer_id"`
		}

		var body RequestBody;
		if err := c.ShouldBindJSON(&body); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json",
			})
			return;
		}

		if body.AnswerID == nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "field 'answer_id' is required",
			});
			return;
		}

		questionID := c.Param("question_id");

		correct, err := repositories.CheckAnswer(conn, questionID, *body.AnswerID);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return;
		}
		c.JSON(http.StatusOK, gin.H{
			"correct": correct,
		})
	}
}

func DeleteQuiz(conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){

		quizID := c.Param("quiz_id")

		err := repositories.DeleteQuiz(conn, quizID)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return;
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "deleted succesfully",
		})
	}
}

func CreateQuiz(conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){

		type RequestBody struct{
			Name *string `json:"name"`
			Description string `json:"description"`
		}
		var body RequestBody;

		if err := c.ShouldBindJSON(&body); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json",
			})
			return;
		}
		if body.Name == nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "field 'name' is required",
			})
			return;
		}

		newQuiz, err := repositories.CreateQuiz(conn, *body.Name, body.Description);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return;
		}
		c.JSON(http.StatusOK, newQuiz);
	}
}

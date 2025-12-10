package handlers

import (
	"net/http"
	"quiz/db/repositories"
	"quiz/entities/dto"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func CreateQuestion(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		quizID := c.Param("quiz_id")

		var body dto.CreateQuestionDTO

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json",
			})
			return
		}

		if len(body.Answers) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no 'answers' provided in json",
			})
			return
		}

		qID, err := strconv.Atoi(quizID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid 'ID' field format",
			})
			return
		}

		newQuestion, err := repositories.CreateQuestion(conn, qID, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, newQuestion)
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

func GetQuestion(conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){
		questionID := c.Param("question_id");

		question, err := repositories.GetQuestion(conn, questionID);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return;
		}
		c.JSON(http.StatusOK, question);
	}
}

func UpdateQuestion(conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){
		questionID := c.Param("question_id");

		var body dto.UpdateQuestionDTO;
		if err := c.ShouldBindJSON(&body); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json",
			})
			return;
		}

		if body.Text == nil && body.NewCorrectID == nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "fields 'text' and/or 'correct_id' are required in json body",
			})
			return;
		}

		question, err := repositories.UpdateQuestion(conn, questionID, body);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return;
		}

		c.JSON(http.StatusOK, question);
	}
}

func DeleteQuestion (conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){
		questionID := c.Param("question_id");

		err := repositories.DeleteQuestion(conn, questionID);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return;
		}

		c.Status(http.StatusNoContent);
	}
}
package handlers

import (
	"net/http"
	"quiz/db/repositories"
	"quiz/entities/dto"

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
		
		c.Status(http.StatusNoContent);
	}
}

func CreateQuiz(conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){

		var body dto.CreateQuizDTO;

		if err := c.ShouldBindJSON(&body); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json",
			})
			return;
		}
		if body.Name == nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "field 'name' is required in JSON body",
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
		c.JSON(http.StatusCreated, newQuiz);
	}
}

func UpdateQuiz(conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){
		quizID := c.Param("quiz_id");

		var body dto.UpdateQuizDTO;

		if err := c.ShouldBindJSON(&body); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json",
			})
			return;
		}
		if body.Name == nil && body.Description == nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error":"field 'name' and/or 'description' required in json body",
			})
			return;
		}

		newQuiz, err := repositories.UpdateQuiz(conn, quizID, body.Name, body.Description);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return;
		}

		c.JSON(http.StatusOK, newQuiz);
	}
}
package handlers

import (
	"net/http"
	"quiz/db/repositories"
	"quiz/entities/dto"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func CheckAnswer(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {

		type RequestBody struct {
			AnswerID *int `json:"answer_id"`
		}

		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json",
			})
			return
		}

		if body.AnswerID == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "field 'answer_id' is required in JSON body",
			})
			return
		}

		questionID := c.Param("question_id")

		correct, err := repositories.CheckAnswer(conn, questionID, *body.AnswerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"correct": correct,
		})
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

func CreateAnswer(conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){
		questionID, err := strconv.Atoi(c.Param("question_id"));
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid question_id format",
			})
			return;
		}

		var body dto.CreateAnswerDTO;

		if err := c.ShouldBindJSON(&body); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json format",
			})
			return;
		}

		if body.Text == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "field 'text' is required in json body",
			})
			return;
		}

		newAnswer, err := repositories.CreateAnswer(conn, questionID, body);
		if err != nil{
			errorMsg := err.Error();

			if strings.Contains(errorMsg, "not found"){
				c.JSON(http.StatusNotFound, gin.H{
					"error": errorMsg,
				})
				return;
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errorMsg,
			})
			return;
		}

		c.JSON(http.StatusCreated, newAnswer);
	}
}

func GetAnswer(conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){
		answerID, err := strconv.Atoi(c.Param("answer_id"));

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid answer_id value",
			})
			return;
		}

		answer, err := repositories.GetAnswer(conn, answerID);
		
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return;
		}
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return;
		}

		c.JSON(http.StatusOK, answer);
	}
}

func DeleteAnswer(conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){
		answerID, err := strconv.Atoi(c.Param("answer_id"));
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid answer_id format",
			})
			return;
		}

		err = repositories.DeleteAnswer(conn, answerID);
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return;
		}
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return;
		}

		c.Status(http.StatusNoContent);
	}
}

func UpdateAnswer(conn *pgx.Conn) gin.HandlerFunc{
	return func(c *gin.Context){
		answerID, err := strconv.Atoi(c.Param("answer_id"));
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid answer_id format",
			})
			return;
		}

		var body dto.UpdateAnswerDTO

		if err := c.ShouldBindJSON(&body); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json format",
			})
			return;
		}

		if body.Text == nil && body.NewCorrectID == nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "fields text and/or correct_id are required in json body",
			})
			return;
		}

		updatedAnswer, err := repositories.UpdateAnswer(conn, answerID, body);
		if err != nil{
			errMsg := err.Error();
			if strings.Contains(errMsg, "not found"){
				c.JSON(http.StatusNotFound, gin.H{
					"error": errMsg,
				})
				return;
			}
			if strings.Contains(errMsg, "new correct answer id"){
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errMsg,
				})
				return;
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errMsg,
			})
			return;
		}

		c.JSON(http.StatusOK, updatedAnswer);
	}
}
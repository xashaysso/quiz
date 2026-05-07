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

type AnswerHandler struct {
	Repo *repositories.PgAnswerRepo
}

func (h *AnswerHandler) CheckAnswer(c *gin.Context) {
		ctx := c.Request.Context()

		var body dto.CheckAnswerDTO;
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

		correct, err := h.Repo.CheckAnswer(ctx, questionID, *body.AnswerID)
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


func (h *AnswerHandler) ListAnswers(c *gin.Context){
		ctx := c.Request.Context()
		question_id := c.Param("question_id")

		answers, err := h.Repo.GetQuizAnswers(ctx, question_id);
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()});
			return;
		}
		c.JSON(http.StatusOK, answers);
	}

func (h *AnswerHandler) CreateAnswer(c *gin.Context){
		ctx := c.Request.Context()
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

		newAnswer, err := h.Repo.CreateAnswer(ctx, questionID, body);
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

func (h *AnswerHandler) GetAnswer(c *gin.Context){
		ctx := c.Request.Context()
		answerID, err := strconv.Atoi(c.Param("answer_id"));

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid answer_id value",
			})
			return;
		}

		answer, err := h.Repo.GetAnswer(ctx, answerID);
		
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

func (h *AnswerHandler) DeleteAnswer(c *gin.Context){
		ctx := c.Request.Context()
		answerID, err := strconv.Atoi(c.Param("answer_id"));
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid answer_id format",
			})
			return;
		}

		err = h.Repo.DeleteAnswer(ctx, answerID);
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

func (h *AnswerHandler) UpdateAnswer(c *gin.Context){
		ctx := c.Request.Context()
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

		updatedAnswer, err := h.Repo.UpdateAnswer(ctx, answerID, body);
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
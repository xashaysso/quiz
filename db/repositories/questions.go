package repositories

import (
	"context"
	entities "quiz/entities/api"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetQuizQuestions(c *gin.Context, conn *pgx.Conn) ([]entities.Question, error) {
	ctx := context.Background()
	id := c.Param("quiz_id")

	var questionList []entities.Question
	rows, err := conn.Query(ctx, `SELECT id, text FROM questions WHERE quiz_id = $1`, id)
	if rows.Err() != nil {
		return nil, err
	}

	for rows.Next() {
		var question entities.Question
		err := rows.Scan(&question.ID, &question.Text)
		if err != nil {
			return nil, err
		}
		questionList = append(questionList, question)
	}
	return questionList, nil
}
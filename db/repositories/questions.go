package repositories

import (
	"context"
	"fmt"
	APIentities "quiz/entities/api"
	entities "quiz/entities/db"
	"quiz/entities/dto"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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


func CreateQuestion(conn *pgx.Conn, quizID int, data dto.CreateQuestionDTO)(APIentities.QuestionAPI, error){
	ctx := context.Background();

	tx, err := conn.Begin(ctx);
	if err != nil{
		return APIentities.QuestionAPI{}, err;
	}
	defer tx.Rollback(ctx);

	var questionID int;

	err = tx.QueryRow(ctx, `INSERT INTO questions (quiz_id, text) VALUES ($1, $2) RETURNING id`, quizID, data.Text).Scan(&questionID);
	if err != nil{
		if pgErr, ok := err.(*pgconn.PgError); ok {
            if pgErr.Code == "23503" { 
                return APIentities.QuestionAPI{}, fmt.Errorf("quiz with id = %d not found", quizID)
            }
        }
		return APIentities.QuestionAPI{}, err;
	}

	answerIDs := make([]int, len(data.Answers));

	for i, answer := range data.Answers {
		var answerID int;
		
		err = tx.QueryRow(ctx, `INSERT INTO answers (question_id, text) VALUES ($1, $2) RETURNING id`, questionID, answer.Text).Scan(&answerID);
		if err != nil{
			return APIentities.QuestionAPI{}, err;
		}

		answerIDs[i] = answerID;
	}

	correctAnswerID := answerIDs[data.CorrectID];

	_, err = tx.Exec(ctx, `UPDATE questions SET correct_answer_id = $1 WHERE id = $2`, correctAnswerID, questionID);
	if err != nil{
		return APIentities.QuestionAPI{}, err;
	}

	err = tx.Commit(ctx);
	if err != nil{
		return APIentities.QuestionAPI{}, err;
	}

	var resData APIentities.QuestionAPI;

	resData.ID = questionID;
	resData.CorrectID = correctAnswerID;
	resData.Text = data.Text;
	resData.Answers = make([]APIentities.AnswerAPI, len(data.Answers));

	for i := range data.Answers{
		resData.Answers[i].ID = answerIDs[i];
		resData.Answers[i].Text = data.Answers[i].Text;
	}

	return resData, nil;
}
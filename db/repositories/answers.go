package repositories

import (
	"context"
	entities "quiz/entities/db"

	"github.com/jackc/pgx/v5"
)

func GetQuizAnswers(conn *pgx.Conn, question_id string) ([]entities.Answer, error) {
	ctx := context.Background()

	var answerList []entities.Answer
	rows, err := conn.Query(ctx, `SELECT id, text FROM answers WHERE question_id = $1`, question_id)
	if rows.Err() != nil {
		return nil, err
	}

	for rows.Next() {
		var answer entities.Answer
		err := rows.Scan(&answer.ID, &answer.Text)
		if err != nil {
			return nil, err
		}
		answerList = append(answerList, answer)
	}
	return answerList, nil
}

func CheckAnswer(conn *pgx.Conn, questionID string, answerID int) (bool, error) {
	var isCorrect bool

	ctx := context.Background()

	err := conn.QueryRow(ctx, `SELECT CASE WHEN correct_answer_id = $2 THEN true ELSE false END
						FROM questions
						WHERE id = $1`, questionID, answerID).Scan(&isCorrect)
	if err != nil {
		return false, err
	}
	return isCorrect, nil
}
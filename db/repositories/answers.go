package repositories

import (
	"context"
	"fmt"
	APIentities "quiz/entities/api"
	entities "quiz/entities/db"
	"quiz/entities/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

	err := conn.QueryRow(ctx, `SELECT correct FROM answers WHERE id = $1`, answerID).Scan(&isCorrect)
	if err != nil {
		return false, err
	}
	return isCorrect, nil
}

func CreateAnswer(conn *pgx.Conn, questionID int, data dto.CreateAnswerDTO)(APIentities.AnswerAPI, error){
	ctx := context.Background();

	var newAnswer APIentities.AnswerAPI;

	tx, err := conn.Begin(ctx);
	if err != nil{
		return APIentities.AnswerAPI{}, err;
	}
	defer tx.Rollback(ctx);

	if data.IsCorrect {
		_, err = tx.Exec(ctx, `UPDATE answers SET correct = FALSE WHERE question_id = $1`, questionID);
		if err != nil{
			return APIentities.AnswerAPI{}, fmt.Errorf("failed to reset correct flag: %w", err);
		}
	}

	err = tx.QueryRow(ctx, `INSERT INTO answers (text, correct, question_id) 
						VALUES ($1, $2, $3) returning id, text, correct`, data.Text, data.IsCorrect, questionID).Scan(&newAnswer.ID, &newAnswer.Text, &newAnswer.IsCorrect);
	if err != nil{
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23503" {
			return APIentities.AnswerAPI{}, fmt.Errorf("question with id %d not found", questionID);
		}
		return APIentities.AnswerAPI{}, err;
	}

	if err := tx.Commit(ctx); err != nil{
		return APIentities.AnswerAPI{}, err;
	}

	return newAnswer, nil
}

func GetAnswer(conn *pgx.Conn, answerID int)(APIentities.AnswerAPI, error){
	ctx := context.Background();

	var answer APIentities.AnswerAPI;
	err := conn.QueryRow(ctx, `SELECT id, text, correct FROM answers WHERE id = $1 ORDER BY id`, answerID).Scan(&answer.ID, &answer.Text, &answer.IsCorrect);
	
	if err == pgx.ErrNoRows {
		return APIentities.AnswerAPI{}, fmt.Errorf("answer with id %d not found", answerID);
	}
	if err != nil{
		return APIentities.AnswerAPI{}, err;
	}

	return answer, nil;
}

func DeleteAnswer(conn *pgx.Conn, answerID int)(error){
	ctx := context.Background();

	cmdTag, err := conn.Exec(ctx, `DELETE FROM answers WHERE id = $1`, answerID);
	if err != nil{
		return err;
	}
	if cmdTag.RowsAffected() == 0{
		return fmt.Errorf("no answer with id %d found", answerID);
	}

	return nil;
}

func UpdateAnswer(conn *pgx.Conn, answerID int, data dto.UpdateAnswerDTO)(APIentities.AnswerAPI, error){
	ctx := context.Background();

	tx, err := conn.Begin(ctx);
	if err != nil{
		return APIentities.AnswerAPI{}, err;
	}
	defer tx.Rollback(ctx);

	var questionID int;
	err = tx.QueryRow(ctx, `SELECT question_id FROM answers WHERE id = $1`, answerID).Scan(&questionID);
	if err == pgx.ErrNoRows{
		return APIentities.AnswerAPI{}, fmt.Errorf("answer with id %d not found", answerID);
	}
	if err != nil{
		return APIentities.AnswerAPI{}, err;
	}

	if data.Text != nil{
		tag, err := tx.Exec(ctx, `UPDATE answers SET text = $1 WHERE id = $2`, *data.Text, answerID);
		if err != nil{
			return APIentities.AnswerAPI{}, err;
		}
		if tag.RowsAffected() == 0{
			return APIentities.AnswerAPI{}, fmt.Errorf("answer with id %d not found", answerID);
		}
	}

	if data.NewCorrectID != nil {
		_, err := tx.Exec(ctx, `UPDATE answers SET correct = FALSE WHERE question_id = $1`, questionID);
		if err != nil{
			return APIentities.AnswerAPI{}, fmt.Errorf("failed to reset correct flag: %w", err);
		}

		tag, err := tx.Exec(ctx, `UPDATE answers SET correct = TRUE WHERE id = $1 AND question_id = $2`, *data.NewCorrectID, questionID);
		if err != nil{
			return APIentities.AnswerAPI{}, fmt.Errorf("failed to set new correct flag: %w", err);
		}
		if tag.RowsAffected() == 0{
			return APIentities.AnswerAPI{}, fmt.Errorf("new correct answer id %d belongs to another question", *data.NewCorrectID);
		}
	}

	if err := tx.Commit(ctx); err != nil{
		return APIentities.AnswerAPI{}, err;
	}

	return GetAnswer(conn, answerID);
}
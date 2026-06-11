package pg

import (
	"context"
	"errors"
	"fmt"
	"quiz/db/repositories"
	entities "quiz/entities/db"
	"quiz/entities/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// create db pool

type PgAnswerRepo struct {
	Pool *pgxpool.Pool
}

func NewAnswerRepo(p *pgxpool.Pool) *PgAnswerRepo {
	return &PgAnswerRepo{Pool: p}
}

// repo methods

func (r *PgAnswerRepo) GetQuizAnswers(ctx context.Context, questionID int) ([]entities.Answer, error) {
	var answerList []entities.Answer
	rows, err := r.Pool.Query(ctx, `SELECT id, text, correct, question_id FROM answers WHERE question_id = $1`, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var answer entities.Answer
		err := rows.Scan(&answer.ID, &answer.Text, &answer.IsCorrect, &answer.QuestionID)
		if err != nil {
			return nil, err
		}
		answerList = append(answerList, answer)
	}
	return answerList, nil
}

func (r *PgAnswerRepo) CheckAnswer(ctx context.Context, questionID int, answerID int) (bool, error) {
	var isCorrect bool

	err := r.Pool.QueryRow(ctx, `SELECT correct FROM answers WHERE id = $1 AND question_id = $2`, answerID, questionID).Scan(&isCorrect)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, repositories.ErrRecordNotFound
		}
		return false, err
	}
	return isCorrect, nil
}

func (r *PgAnswerRepo) CreateAnswer(ctx context.Context, txAny any, questionID int, text string, isCorrect bool) (int, error) {
	tx, ok := txAny.(pgx.Tx)
	if !ok {
		return -1, fmt.Errorf("invalid transaction type: expected pgx.Tx, got %T", txAny)
	}

	if isCorrect {
		_, err := tx.Exec(ctx, `UPDATE answers SET correct = FALSE WHERE question_id = $1`, questionID)
		if err != nil {
			return -1, err
		}
	}

	var id int
	err := tx.QueryRow(ctx, `INSERT INTO answers (text, correct, question_id) 
						VALUES ($1, $2, $3) returning id`, text, isCorrect, questionID).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r *PgAnswerRepo) GetAnswer(ctx context.Context, answerID int) (entities.Answer, error) {
	var answer entities.Answer
	err := r.Pool.QueryRow(ctx, `SELECT id, text, correct, question_id FROM answers WHERE id = $1 ORDER BY id`, answerID).Scan(&answer.ID, &answer.Text, &answer.IsCorrect, &answer.QuestionID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Answer{}, repositories.ErrRecordNotFound
		}
		return entities.Answer{}, err
	}

	return answer, nil
}

func (r *PgAnswerRepo) DeleteAnswer(ctx context.Context, answerID int) error {
	cmdTag, err := r.Pool.Exec(ctx, `DELETE FROM answers WHERE id = $1`, answerID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return repositories.ErrRecordNotFound
	}

	return nil
}

func (r *PgAnswerRepo) UpdateAnswer(ctx context.Context, answerID int, data dto.UpdateAnswerDTO) (entities.Answer, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return entities.Answer{}, err
	}
	defer tx.Rollback(ctx)

	var questionID int
	err = tx.QueryRow(ctx, `SELECT question_id FROM answers WHERE id = $1`, answerID).Scan(&questionID)
	if err != nil {
		return entities.Answer{}, err
	}

	if data.Text != nil {
		_, err := tx.Exec(ctx, `UPDATE answers SET text = $1 WHERE id = $2`, *data.Text, answerID)
		if err != nil {
			return entities.Answer{}, err
		}
	}

	if data.NewCorrectID != nil {
		_, err := tx.Exec(ctx, `UPDATE answers SET correct = FALSE WHERE question_id = $1 AND correct = TRUE`, questionID)
		if err != nil {
			return entities.Answer{}, err
		}

		tag, err := tx.Exec(ctx, `UPDATE answers SET correct = TRUE WHERE id = $1 AND question_id = $2`, *data.NewCorrectID, questionID)
		if err != nil {
			return entities.Answer{}, err
		}
		if tag.RowsAffected() == 0 {
			return entities.Answer{}, repositories.ErrInvalidCorrectID
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return entities.Answer{}, err
	}

	return r.GetAnswer(ctx, answerID)
}

func (r *PgAnswerRepo) CheckIfAnswerOwner(ctx context.Context, answerID int, userID int) (bool, error) {
	var isOwner bool

	err := r.Pool.QueryRow(ctx, `SELECT EXISTS (
								SELECT 1 FROM answers a 
								JOIN questions q ON a.question_id = q.id
								JOIN quiz qz ON q.quiz_id = qz.id
								WHERE a.id = $1 AND qz.creator_id = $2
								);`, answerID, userID).Scan(&isOwner)
	if err != nil {
		return false, err
	}
	return isOwner, nil
}

func (r *PgAnswerRepo) CheckIfQuestionOwner(ctx context.Context, questionID int, userID int) (bool, error) {
	var isOwner bool

	err := r.Pool.QueryRow(ctx, `SELECT EXISTS (
								SELECT 1 FROM questions q
								JOIN quiz qz ON q.quiz_id = qz.id
								WHERE q.id = $1 AND qz.creator_id = $2
								);`, questionID, userID).Scan(&isOwner)
	if err != nil {
		return false, err
	}
	return isOwner, nil
}

func (r *PgAnswerRepo) GetAnswersByQuestionIDs(ctx context.Context, questionIDs []int) ([]entities.Answer, error) {
	if len(questionIDs) == 0 {
		return []entities.Answer{}, nil
	}

	rows, err := r.Pool.Query(ctx, `SELECT id, question_id, text, correct FROM answers WHERE question_id = ANY($1) ORDER BY id`, questionIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []entities.Answer
	for rows.Next() {
		var a entities.Answer
		err := rows.Scan(&a.ID, &a.QuestionID, &a.Text, &a.IsCorrect)
		if err != nil {
			return nil, err
		}

		answers = append(answers, a)
	}

	return answers, nil
}

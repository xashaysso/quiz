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

type PgQuestionRepo struct {
	Pool *pgxpool.Pool
}

func NewQuestionRepo(p *pgxpool.Pool) *PgQuestionRepo {
	return &PgQuestionRepo{Pool: p}
}

// repo methods

func (r *PgQuestionRepo) GetQuizQuestions(ctx context.Context, id int) ([]entities.Question, error) {
	rows, err := r.Pool.Query(ctx, `SELECT id, text, quiz_id FROM questions WHERE quiz_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []entities.Question

	for rows.Next() {
		var q entities.Question
		if err := rows.Scan(&q.ID, &q.Text, &q.QuizID); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, nil
}

func (r *PgQuestionRepo) CreateQuestion(ctx context.Context, txAny any, quizID int, data dto.CreateQuestionDTO) (int, error) {
	var questionID int

	tx, ok := txAny.(pgx.Tx)
	if !ok {
		return -1, fmt.Errorf("invalid transaction type: expected pgx.Tx, got %T", tx)
	}

	err := tx.QueryRow(ctx, `INSERT INTO questions (quiz_id, text) VALUES ($1, $2) RETURNING id`, quizID, data.Text).Scan(&questionID)
	if err != nil {
		return -1, err
	}

	return questionID, nil
}

func (r *PgQuestionRepo) GetQuestion(ctx context.Context, questionID int) (entities.Question, error) {
	var question entities.Question

	err := r.Pool.QueryRow(ctx, `SELECT id, text, quiz_id FROM questions WHERE id = $1 ORDER BY id`, questionID).Scan(&question.ID, &question.Text, &question.QuizID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Question{}, repositories.ErrRecordNotFound
		}
		return entities.Question{}, err
	}

	return question, nil
}

func (r *PgQuestionRepo) UpdateQuestion(ctx context.Context, questionID int, data dto.UpdateQuestionDTO) (entities.Question, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return entities.Question{}, err
	}

	defer tx.Rollback(ctx)

	if data.Text != nil {
		tag, err := tx.Exec(ctx, `UPDATE questions SET text = $1 WHERE id = $2`, *data.Text, questionID)
		if err != nil {
			return entities.Question{}, err
		}
		if tag.RowsAffected() == 0 {
			return entities.Question{}, repositories.ErrRecordNotFound
		}
	}

	if data.NewCorrectID != nil {
		newCorrectID := data.NewCorrectID

		_, err = tx.Exec(ctx, `UPDATE answers SET correct = FALSE WHERE question_id = $1`, questionID)
		if err != nil {
			return entities.Question{}, err
		}

		tag, err := tx.Exec(ctx, `UPDATE answers SET correct = TRUE WHERE id = $1 AND question_id = $2`, newCorrectID, questionID)
		if err != nil {
			return entities.Question{}, err
		}
		if tag.RowsAffected() == 0 {
			return entities.Question{}, repositories.ErrInvalidCorrectID
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return entities.Question{}, err
	}

	return r.GetQuestion(ctx, questionID)
}

func (r *PgQuestionRepo) DeleteQuestion(ctx context.Context, questionID int) error {
	cmdTag, err := r.Pool.Exec(ctx, `DELETE FROM questions WHERE id = $1`, questionID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return repositories.ErrRecordNotFound
	}
	return nil
}

func (r *PgQuestionRepo) CheckIfQuizOwner(ctx context.Context, quizID int, userID int) (bool, error) {
	var isOwner bool

	err := r.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM quiz WHERE id = $1 AND creator_id = $2);`, quizID, userID).Scan(&isOwner)
	if err != nil {
		return false, err
	}

	return isOwner, nil
}

func (r *PgQuestionRepo) CheckIfQuestionOwner(ctx context.Context, questionID int, userID int) (bool, error) {
	var isOwner bool
	err := r.Pool.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM questions q
            JOIN quiz qz ON q.quiz_id = qz.id
            WHERE q.id = $1 AND qz.creator_id = $2
        );`, questionID, userID).Scan(&isOwner)

	return isOwner, err
}

func (r *PgQuestionRepo) GetQuestionIDsByQuizID(ctx context.Context, quizID int64) ([]int64, error) {
	rows, err := r.Pool.Query(ctx, `SELECT id FROM questions WHERE quiz_id = $1 ORDER BY id ASC`, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

package repositories

import (
	"context"
	"fmt"
	APIentities "quiz/entities/api"
	"quiz/entities/dto"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func GetQuizQuestions(c *gin.Context, conn *pgx.Conn) ([]APIentities.QuestionAPI, error) {
	ctx := context.Background()
	id := c.Param("quiz_id")

	rows, err := conn.Query(ctx, `SELECT q.id AS question_id, q.text AS question_text, a.id AS answer_id, a.text AS answer_id, a.correct
								FROM questions q JOIN answers a	ON q.id = a.question_id WHERE q.quiz_id = $1 ORDER BY q.id, a.id`, id)
	if rows.Err() != nil {
		return nil, err
	}

	var questionMap = make(map[int]APIentities.QuestionAPI);
	var orderedQuestionIDs []int;

	for rows.Next() {
		var qID, aID int;
		var correct bool;
		var qText, aText string;


		err := rows.Scan(&qID, &qText, &aID, &aText, &correct);
		if err != nil {
			return nil, err
		}
		defer rows.Close();

		question, exists := questionMap[qID];
		if !exists{
			question = APIentities.QuestionAPI{
				ID: qID,
				Text: qText,
				Answers: []APIentities.AnswerAPI{},
			}
			orderedQuestionIDs = append(orderedQuestionIDs, qID);
		}

		answer := APIentities.AnswerAPI{
			ID: aID,
			Text: aText,
			IsCorrect: correct,
		}
		question.Answers = append(question.Answers, answer);

		questionMap[qID] = question;
	}

	result := make([]APIentities.QuestionAPI, 0, len(questionMap));
	for _, id := range orderedQuestionIDs {
		result = append(result, questionMap[id]);
	}

	return result, nil
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
		
		err = tx.QueryRow(ctx, `INSERT INTO answers (question_id, text, correct) VALUES ($1, $2, $3) RETURNING id`, questionID, answer.Text, answer.IsCorrect).Scan(&answerID);
		if err != nil{
			return APIentities.QuestionAPI{}, err;
		}

		answerIDs[i] = answerID;
	}

	err = tx.Commit(ctx);
	if err != nil{
		return APIentities.QuestionAPI{}, err;
	}

	var resData APIentities.QuestionAPI;

	resData.ID = questionID;
	resData.Text = data.Text;
	resData.Answers = make([]APIentities.AnswerAPI, len(data.Answers));

	for i := range data.Answers{
		resData.Answers[i].ID = answerIDs[i];
		resData.Answers[i].IsCorrect = data.Answers[i].IsCorrect;
		resData.Answers[i].Text = data.Answers[i].Text;
	}

	return resData, nil;
}

func GetQuestion(conn *pgx.Conn, questionID string)(APIentities.QuestionAPI, error){
	ctx := context.Background();

	var question APIentities.QuestionAPI;
	questionAnswers := []APIentities.AnswerAPI{};

	err := conn.QueryRow(ctx, `SELECT id, text FROM questions WHERE id = $1 ORDER BY id`, questionID).Scan(&question.ID, &question.Text);
	if err != nil{
		return APIentities.QuestionAPI{}, err;
	}

	rows, err := conn.Query(ctx, `SELECT id, text, correct FROM answers WHERE question_id = $1 ORDER BY id`, questionID);
	if err != nil{
		return APIentities.QuestionAPI{}, nil;
	}
	defer rows.Close();

	for rows.Next(){
		var answer APIentities.AnswerAPI;
		err = rows.Scan(&answer.ID, &answer.Text, &answer.IsCorrect);
		if err != nil{
			return APIentities.QuestionAPI{}, nil;
		}

		questionAnswers = append(questionAnswers, answer);
	}
	question.Answers = questionAnswers;

	return question, nil;
}

func UpdateQuestion(conn *pgx.Conn, questionID string, data dto.UpdateQuestionDTO)(APIentities.QuestionAPI, error){
	ctx := context.Background();

	tx, err := conn.Begin(ctx);
	if err != nil{
		return APIentities.QuestionAPI{}, err;
	}

	defer tx.Rollback(ctx);

	if data.Text != nil{
		tag, err := conn.Exec(ctx, `UPDATE questions SET text = $1 WHERE id = $2`, *data.Text, questionID)
		if err != nil{
			return APIentities.QuestionAPI{}, err;
		}
		if tag.RowsAffected() == 0{
			return APIentities.QuestionAPI{}, fmt.Errorf("question ID %s not found", questionID);
		}
	}

	if data.NewCorrectID != nil{
		newCorrectID := data.NewCorrectID;

		_, err = tx.Exec(ctx, `UPDATE answers SET correct = FALSE WHERE question_id = $1`, questionID);
		if err != nil{
			return APIentities.QuestionAPI{}, fmt.Errorf("failed to reset correct flag: %w", err);
		}

		tag, err := tx.Exec(ctx, `UPDATE answers SET correct = TRUE WHERE id = $1`, newCorrectID);
		if err != nil{
			return APIentities.QuestionAPI{}, fmt.Errorf("failed to set new correct flag: %w", err);
		}
		if tag.RowsAffected() == 0{
			return APIentities.QuestionAPI{}, fmt.Errorf("answer ID %d not found", newCorrectID);
		}
	}

	if err = tx.Commit(ctx); err != nil{
		return APIentities.QuestionAPI{}, err;
	}

	return GetQuestion(conn, questionID);
}

func DeleteQuestion(conn *pgx.Conn, questionID string)(error){
	ctx := context.Background();

	cmdTag, err := conn.Exec(ctx, `DELETE FROM questions WHERE id = $1`, questionID);
	if err != nil{
		return err;
	}
	if cmdTag.RowsAffected() == 0{
		return fmt.Errorf("question not found");
	}
	return nil;
}
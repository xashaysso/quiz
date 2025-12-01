package repositories

import (
	"context"
	entities "quiz/entities/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetQuiz(conn *pgx.Conn)([]entities.Quiz, error){
	ctx := context.Background();

	var quizList []entities.Quiz;
	rows, err := conn.Query(ctx, `SELECT id, name FROM quiz`)
	if rows.Err() != nil {
		return nil, err;
	}
	for rows.Next(){
		var quiz entities.Quiz;
		err := rows.Scan(&quiz.ID, &quiz.Name);
		if err != nil{
			return nil, err;
		}
		quizList = append(quizList, quiz);
	}
	return quizList, nil;
}

func GetQuizQuestions(c *gin.Context, conn *pgx.Conn)([]entities.Question, error){
	ctx := context.Background();
	id := c.Param("quiz_id");

	var questionList []entities.Question;
	rows, err := conn.Query(ctx, `SELECT id, text FROM questions WHERE quiz_id = $1`, id);
	if rows.Err() != nil{
		return nil, err;
	}

	for rows.Next(){
		var question entities.Question;
		err := rows.Scan(&question.ID, &question.Text);
		if err != nil{
			return nil, err;
		}
		questionList = append(questionList, question);
	}
	return questionList, nil;
}

func GetQuizAnswers(c *gin.Context, conn *pgx.Conn)([]entities.Answer, error){
	ctx := context.Background();
	question_id := c.Param("question_id");

	var answerList []entities.Answer;
	rows, err := conn.Query(ctx, `SELECT id, text FROM answers WHERE question_id = $1`, question_id);
	if rows.Err() != nil{
		return nil, err;
	}

	for rows.Next() {
		var answer entities.Answer;
		err := rows.Scan(&answer.ID, &answer.Text);
		if err != nil{
			return nil, err;
		}
		answerList = append(answerList, answer);
	}
	return answerList, nil;
}

func CheckAnswer(conn *pgx.Conn, questionID string, answerID int)(bool, error){
	var isCorrect bool;

	ctx := context.Background();

	err := conn.QueryRow(ctx, `SELECT CASE WHEN correct_answer_id = $2 THEN true ELSE false END
						FROM questions
						WHERE id = $1`, questionID, answerID).Scan(&isCorrect);
	if err != nil{
		return false, err;
	}
	return isCorrect, nil;
}


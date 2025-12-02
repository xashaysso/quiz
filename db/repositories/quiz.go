package repositories

import (
	"context"
	"fmt"
	entities "quiz/entities/db"

	"github.com/jackc/pgx/v5"
)

func GetQuiz(conn *pgx.Conn)([]entities.Quiz, error){
	ctx := context.Background();

	var quizList []entities.Quiz;
	rows, err := conn.Query(ctx, `SELECT id, name, description FROM quiz`)
	if rows.Err() != nil {
		return nil, err;
	}
	for rows.Next(){
		var quiz entities.Quiz;
		err := rows.Scan(&quiz.ID, &quiz.Name, &quiz.Description);
		if err != nil{
			return nil, err;
		}
		quizList = append(quizList, quiz);
	}
	return quizList, nil;
}

func DeleteQuiz(conn *pgx.Conn, quizID string)(error){
	ctx := context.Background();

	cmdTag, err := conn.Exec(ctx, `DELETE FROM quiz WHERE id = $1`, quizID);
	if err != nil{
		return err;
	}

	if cmdTag.RowsAffected() == 0{
		return fmt.Errorf("quiz not found");
	}
	return nil;
}

func CreateQuiz(conn *pgx.Conn, quiz_name string, quiz_description string)(entities.Quiz, error){
	ctx := context.Background();
	var newQuiz entities.Quiz;

	err := conn.QueryRow(ctx, `INSERT INTO quiz (name, description) VALUES ($1, $2) RETURNING id, name, description`, quiz_name, quiz_description).Scan(&newQuiz.ID, &newQuiz.Name, &newQuiz.Description);
	if err != nil{
		return entities.Quiz{}, err;
	}
	return newQuiz, nil;
}


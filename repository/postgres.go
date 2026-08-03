package repository

import (
	"TeacherBot/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type postgresRepository struct {
	conn *pgx.Conn
	ctx  context.Context
}

func (r *postgresRepository) InsertTecher(telegram_id int, telegram_name *string, teacher_name string) error {
	SQLQuery := `
INSERT INTO bot.teacher (telegram_id, telegram_name, teacher_name)
VALUES ($1,$2,$3);
`
	queryCtx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	_, err := r.conn.Exec(queryCtx, SQLQuery, telegram_id, telegram_name, teacher_name)
	return err
}

func (r *postgresRepository) InsertUser(
	telegram_id int,
	telegram_name *string,
	full_name string,
	teacher_ID int) error {
	SQLQuery := `
INSERT INTO bot.users (telegram_id, telegram_name, full_name, teacher_ID)
VALUES ($1,$2,$3,$4);
`
	queryCtx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	_, err := r.conn.Exec(queryCtx, SQLQuery, telegram_id, telegram_name, full_name, teacher_ID)
	return err
}

func (r *postgresRepository) InsertTests(
	subject string,
	level string,
	topic string,
	test string,
	result int,
	time_finish time.Time,
	user_id int) error {
	SQLQuery := `
INSERT INTO bot.tests (subject,level, topic, test,result,time_finish,user_id)
VALUES ($1,$2,$3,$4,$5,$6,$7);
`
	queryCtx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	_, err := r.conn.Exec(queryCtx, SQLQuery, subject, level, topic, test, result, time_finish, user_id)
	return err
}

func (r *postgresRepository) UpdateRow(
	table string,
	column string,
	value any,
	indexColumn string,
	index int) error {
	queryCtx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	SQLQuery := fmt.Sprintf("UPDATE %s SET %s = $1 WHERE %s=$2", "bot."+table, column, indexColumn)
	_, err := r.conn.Exec(queryCtx, SQLQuery, value, index)
	return err
}

func (r *postgresRepository) DeleteRow(table string, column string, index int) error {
	queryCtx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	SQLQuery := fmt.Sprintf("DELETE FROM %s WHERE %s=$1", "bot."+table, column)
	_, err := r.conn.Exec(queryCtx, SQLQuery, index)
	return err
}

func (r *postgresRepository) GetStudentsByTeacher(teacher_ID int, limit int) ([]models.User, error) {
	SQLQuery := `
SELECT telegram_id,telegram_name ,full_name
FROM bot.users
WHERE teacher_ID = $1
LIMIT $2
`
	queryCtx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	rows, err := r.conn.Query(queryCtx, SQLQuery, teacher_ID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	students := []models.User{}
	for rows.Next() {
		var telegram_id int
		var telegram_name string
		var full_name string
		if err := rows.Scan(&telegram_id, &telegram_name, &full_name); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		students = append(students, models.User{Telegram_id: telegram_id, Telegram_name: telegram_name, Full_name: full_name})
	}
	return students, rows.Err()
}

func (r *postgresRepository) GetResultByUser(user_id int, limit int) ([]models.UserResult, error) {
	SQLQuery := `
SELECT id,result,time_finish
FROM bot.tests
WHERE user_id = $1
ORDER BY id ASC LIMIT $2
`
	queryCtx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	rows, err := r.conn.Query(queryCtx, SQLQuery, user_id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := []models.UserResult{}
	for rows.Next() {
		var id int
		var result int
		var time_finish time.Time
		if err := rows.Scan(&id, &result, &time_finish); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, models.UserResult{Id: id, Result: result, Time_finish: time_finish})
	}
	return results, rows.Err()
}

func (r *postgresRepository) GetTestByUser(user_id int) ([]models.Test, error) {
	SQLQuery := `
SELECT id,subject,level,topic,test,result
FROM bot.tests
WHERE user_id=$1
`
	queryCtx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	rows, err := r.conn.Query(queryCtx, SQLQuery, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	testList := []models.Test{}
	for rows.Next() {
		var id int
		var subject string
		var level string
		var topic string
		var test string
		var result int
		if err := rows.Scan(&id, &subject, &level, &topic, &test, &result); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		testList = append(testList, models.Test{Id: id, Subject: subject, Level: level, Topic: topic, Test: test, Result: result})
	}
	return testList, rows.Err()
}
func (r *postgresRepository) InitializationRow(table_name string, colum_name string, value any) error {
	SQLQuery := fmt.Sprintf(`
SELECT EXISTS(
SELECT * FROM %s 
WHERE %s = $1);
`, "bot."+table_name, colum_name)
	queryCtx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	var exists bool
	err := r.conn.QueryRow(queryCtx, SQLQuery, value).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check row existence: %w", err)
	}

	if !exists {
		return sql.ErrNoRows // Возвращаем стандартную ошибку "no rows in result set"
	}

	return nil
}
func (r *postgresRepository) GetTeacherLists() ([]models.Teacher, error) {
	SQLQuery := `
SELECT telegram_id,teacher_name
FROM bot.teacher
`
	queryCtx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	rows, err := r.conn.Query(queryCtx, SQLQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	teacherLists := []models.Teacher{}
	for rows.Next() {
		var telegram_id int
		var teacher_name string
		if err := rows.Scan(&telegram_id, &teacher_name); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		teacherLists = append(teacherLists, models.Teacher{Telegram_id: telegram_id, Teacher_name: teacher_name})
	}
	return teacherLists, nil
}

func (r *postgresRepository) Close() error {
	if r.conn != nil {
		return r.conn.Close(context.Background())
	}
	return nil
}

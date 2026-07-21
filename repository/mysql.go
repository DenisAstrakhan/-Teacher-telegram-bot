package repository

import (
	"TeacherBot/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type mysqlRepository struct {
	conn *sql.DB
	ctx  context.Context
}

func (r *mysqlRepository) InsertTecher(telegram_id int, telegram_name *string, teacher_name string) error {
	SQLQuery := `
INSERT INTO teacher (telegram_id, telegram_name, teacher_name)
VALUES (?,?,?);
`
	_, err := r.conn.ExecContext(r.ctx, SQLQuery, telegram_id, telegram_name, teacher_name)
	return err
}

func (r *mysqlRepository) InsertUser(
	telegram_id int,
	telegram_name *string,
	full_name string,
	teacher_ID int) error {
	SQLQuery := `
INSERT INTO users (telegram_id, telegram_name, full_name, teacher_ID)
VALUES (?,?,?,?);
`
	_, err := r.conn.ExecContext(r.ctx, SQLQuery, telegram_id, telegram_name, full_name, teacher_ID)
	return err
}

func (r *mysqlRepository) InsertTests(
	subject string,
	level string,
	topic string,
	test string,
	result int,
	time_finish time.Time,
	user_id int) error {
	SQLQuery := `
INSERT INTO tests (subject,level, topic, test,result,time_finish,user_id)
VALUES (?,?,?,?,?,?,?);
`
	_, err := r.conn.ExecContext(r.ctx, SQLQuery, subject, level, topic, test, result, time_finish, user_id)
	return err
}

func (r *mysqlRepository) UpdateRow(
	table string,
	column string,
	value any,
	indexColumn string,
	index int) error {
	SQLQuery := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s=?", table, column, indexColumn)
	_, err := r.conn.ExecContext(r.ctx, SQLQuery, value, index)
	return err
}

func (r *mysqlRepository) DeleteRow(table string, column string, index int) error {
	SQLQuery := fmt.Sprintf("DELETE FROM %s WHERE %s=?", table, column)
	_, err := r.conn.ExecContext(r.ctx, SQLQuery, index)
	return err
}

func (r *mysqlRepository) GetStudentsByTeacher(teacher_ID int, limit int) ([]models.User, error) {
	SQLQuery := `
SELECT telegram_id,telegram_name ,full_name
FROM users
WHERE teacher_ID = ?
LIMIT ?
`
	rows, err := r.conn.QueryContext(r.ctx, SQLQuery, teacher_ID, limit)
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

func (r *mysqlRepository) GetResultByUser(user_id int, limit int) ([]models.UserResult, error) {
	SQLQuery := `
SELECT id,result,time_finish
FROM tests
WHERE user_id = ?
ORDER BY id ASC LIMIT ?
`
	rows, err := r.conn.QueryContext(r.ctx, SQLQuery, user_id, limit)
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

func (r *mysqlRepository) GetTestByUser(user_id int) ([]models.Test, error) {
	SQLQuery := `
SELECT id,subject,level,topic,test,result
FROM tests
WHERE user_id=?
`

	rows, err := r.conn.QueryContext(r.ctx, SQLQuery, user_id)
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
func (r *mysqlRepository) InitializationRow(table_name string, colum_name string, value any) error {
	SQLQuery := fmt.Sprintf(`
SELECT EXISTS(
SELECT * FROM %s 
WHERE %s = ?);
`, table_name, colum_name)
	var exists bool
	err := r.conn.QueryRowContext(r.ctx, SQLQuery, value).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check row existence: %w", err)
	}

	if !exists {
		return sql.ErrNoRows // Возвращаем стандартную ошибку "no rows in result set"
	}

	return nil
}
func (r *mysqlRepository) GetTeacherLists() ([]models.Teacher, error) {
	SQLQuery := `
SELECT telegram_id,teacher_name
FROM teacher
`
	rows, err := r.conn.QueryContext(r.ctx, SQLQuery)
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

func (r *mysqlRepository) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

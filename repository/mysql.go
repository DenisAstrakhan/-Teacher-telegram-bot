package repository

import (
	"TeacherBot/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type mysqlRepository struct {
	conn    *sql.DB
	timeout time.Duration
}

func (r *mysqlRepository) InsertTecher(telegram_id int, telegram_name *string, teacher_name string, ctx context.Context) error {
	SQLQuery := `
INSERT INTO teacher (telegram_id, telegram_name, teacher_name)
VALUES (?,?,?);
`
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	_, err := r.conn.ExecContext(queryCtx, SQLQuery, telegram_id, telegram_name, teacher_name)
	return err
}

func (r *mysqlRepository) InsertUser(
	telegram_id int,
	telegram_name *string,
	full_name string,
	teacher_ID int,
	ctx context.Context) error {
	SQLQuery := `
INSERT INTO users (telegram_id, telegram_name, full_name, teacher_ID)
VALUES (?,?,?,?);
`
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	_, err := r.conn.ExecContext(queryCtx, SQLQuery, telegram_id, telegram_name, full_name, teacher_ID)
	return err
}

func (r *mysqlRepository) InsertTests(
	subject string,
	level string,
	topic string,
	test string,
	result int,
	time_finish time.Time,
	user_id int,
	ctx context.Context) error {
	SQLQuery := `
INSERT INTO tests (subject,level, topic, test,result,time_finish,user_id)
VALUES (?,?,?,?,?,?,?);
`
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	_, err := r.conn.ExecContext(queryCtx, SQLQuery, subject, level, topic, test, result, time_finish, user_id)
	return err
}

func (r *mysqlRepository) UpdateRow(
	table string,
	column string,
	value any,
	indexColumn string,
	index int,
	ctx context.Context) error {
	SQLQuery := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s=?", table, column, indexColumn)
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	_, err := r.conn.ExecContext(queryCtx, SQLQuery, value, index)
	return err
}

func (r *mysqlRepository) DeleteRow(table string, column string, index int, ctx context.Context) error {
	SQLQuery := fmt.Sprintf("DELETE FROM %s WHERE %s=?", table, column)
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	_, err := r.conn.ExecContext(queryCtx, SQLQuery, index)
	return err
}

func (r *mysqlRepository) GetStudentsByTeacher(teacher_ID int, limit int, ctx context.Context) ([]models.User, error) {
	SQLQuery := `
SELECT telegram_id,telegram_name ,full_name
FROM users
WHERE teacher_ID = ?
LIMIT ?
`
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	rows, err := r.conn.QueryContext(queryCtx, SQLQuery, teacher_ID, limit)
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

func (r *mysqlRepository) GetResultByUser(user_id int, limit int, ctx context.Context) ([]models.UserResult, error) {
	SQLQuery := `
SELECT id,result,time_finish
FROM tests
WHERE user_id = ?
ORDER BY id ASC LIMIT ?
`
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	rows, err := r.conn.QueryContext(queryCtx, SQLQuery, user_id, limit)
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
		results = append(results, models.UserResult{Id: id, Result: result, Time_finish: &time_finish})
	}
	return results, rows.Err()
}

func (r *mysqlRepository) GetTestByUser(user_id int, ctx context.Context) ([]models.Test, error) {
	SQLQuery := `
SELECT id,subject,level,topic,test,result
FROM tests
WHERE user_id=?
`
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	rows, err := r.conn.QueryContext(queryCtx, SQLQuery, user_id)
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
func (r *mysqlRepository) InitializationRow(table_name string, colum_name string, value any, ctx context.Context) error {
	SQLQuery := fmt.Sprintf(`
SELECT EXISTS(
SELECT * FROM %s 
WHERE %s = ?);
`, table_name, colum_name)
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	var exists bool
	err := r.conn.QueryRowContext(queryCtx, SQLQuery, value).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check row existence: %w", err)
	}

	if !exists {
		return sql.ErrNoRows // Возвращаем стандартную ошибку "no rows in result set"
	}

	return nil
}
func (r *mysqlRepository) GetTeacherLists(ctx context.Context) ([]models.Teacher, error) {
	SQLQuery := `
SELECT telegram_id,teacher_name
FROM teacher
`
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	rows, err := r.conn.QueryContext(queryCtx, SQLQuery)
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
	return r.conn.Close()
}

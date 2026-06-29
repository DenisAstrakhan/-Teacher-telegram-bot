package postgres

import (
	"TeacherBot/domain"
	"TeacherBot/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	conn *pgx.Conn
	ctx  context.Context
}

func NewUserRepository(ctx context.Context) (domain.UserRepository, error) {
	pgConnString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	//pgConnString := fmt.Sprintf("postgres://%s:%s@bot-postgres:5432/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	conn, err := pgx.Connect(ctx, pgConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &userRepository{
		conn: conn,
		ctx:  ctx,
	}, nil
}

func (r *userRepository) InsertTecher(telegram_id int, telegram_name *string, teacher_name string) error {
	SQLQuery := `
INSERT INTO bot.teacher (telegram_id, telegram_name, teacher_name)
VALUES ($1,$2,$3);
`
	_, err := r.conn.Exec(r.ctx, SQLQuery, telegram_id, telegram_name, teacher_name)
	return err
}

func (r *userRepository) InsertUser(
	telegram_id int,
	telegram_name *string,
	full_name string,
	teacher_ID int) error {
	SQLQuery := `
INSERT INTO bot.users (telegram_id, telegram_name, full_name, teacher_ID)
VALUES ($1,$2,$3,$4);
`
	_, err := r.conn.Exec(r.ctx, SQLQuery, telegram_id, telegram_name, full_name, teacher_ID)
	return err
}

func (r *userRepository) InsertTests(
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
	_, err := r.conn.Exec(r.ctx, SQLQuery, subject, level, topic, test, result, time_finish, user_id)
	return err
}

func (r *userRepository) UpdateRow(
	table string,
	column string,
	value any,
	indexColumn string,
	index int) error {
	SQLQuery := fmt.Sprintf("UPDATE %s SET %s = $1 WHERE %s=$2", table, column, indexColumn)
	_, err := r.conn.Exec(r.ctx, SQLQuery, value, index)
	return err
}

func (r *userRepository) DeleteRow(table string, column string, index int) error {
	SQLQuery := fmt.Sprintf("DELETE FROM %s WHERE %s=$1", table, column)
	_, err := r.conn.Exec(r.ctx, SQLQuery, index)
	return err
}

func (r *userRepository) GetStudentsByTeacher(teacher_ID int, limit int) ([]models.User, error) {
	SQLQuery := `
SELECT telegram_id,full_name
FROM bot.users
WHERE teacher_ID = $1
LIMIT $2
`
	rows, err := r.conn.Query(r.ctx, SQLQuery, teacher_ID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	students := []models.User{}
	for rows.Next() {
		var telegram_id int
		var full_name string
		if err := rows.Scan(&telegram_id, &full_name); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		students = append(students, models.User{Telegram_id: telegram_id, Full_name: full_name})
	}
	return students, rows.Err()
}

func (r *userRepository) GetResultByUser(user_id int, limit int) ([]models.UserResult, error) {
	SQLQuery := `
SELECT id,result,time_finish
FROM bot.tests
WHERE user_id = $1
ORDER BY id ASC LIMIT $2
`
	rows, err := r.conn.Query(r.ctx, SQLQuery, user_id, limit)
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

func (r *userRepository) GetTestByUser(id int, user_id int) (models.Test, error) {
	SQLQuery := `
SELECT subject,level,topic,test,result
FROM bot.tests
WHERE id = $1 AND user_id=$2
`
	var subject string
	var level string
	var topic string
	var test string
	var result int
	row := r.conn.QueryRow(r.ctx, SQLQuery, id, user_id)
	if err := row.Scan(&subject, &level, &topic, &test, &result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Test{}, fmt.Errorf("test not found: id=%d, user_id=%d", id, user_id)
		}
		return models.Test{}, fmt.Errorf("scan error: %w", err)
	}
	return models.Test{
		Subject: subject,
		Level:   level,
		Topic:   topic,
		Test:    test,
		Result:  result,
	}, nil
}
func (r *userRepository) InitializationRow(table_name string, colum_name string, value any) error {
	SQLQuery := fmt.Sprintf(`
SELECT EXISTS(
SELECT * FROM %s 
WHERE %s = $1);
`, table_name, colum_name)
	var exists bool
	err := r.conn.QueryRow(r.ctx, SQLQuery, value).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check row existence: %w", err)
	}

	if !exists {
		return sql.ErrNoRows // Возвращаем стандартную ошибку "no rows in result set"
	}

	return nil
}
func (r *userRepository) GetTeacherLists() ([]models.Teacher, error) {
	SQLQuery := `
SELECT telegram_id,teacher_name
FROM bot.teacher
`
	rows, err := r.conn.Query(r.ctx, SQLQuery)
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

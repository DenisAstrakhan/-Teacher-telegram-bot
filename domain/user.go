package domain

import (
	"TeacherBot/models"
	"context"
	"time"
)

type UserRepository struct {
	DataRepository  UserDataRepository
	CacheRepository UserCacheRepository
}

func (r UserRepository) Close() {
	r.CacheRepository.Close()
	r.DataRepository.Close()
}

type UserDataRepository interface {
	InsertTecher(
		telegram_id int,
		telegram_name *string,
		teacher_name string,
		ctx context.Context) error
	InsertUser(
		telegram_id int,
		telegram_name *string,
		full_name string,
		teacher_ID int,
		ctx context.Context) error
	InsertTests(
		subject string,
		level string,
		topic string,
		test string,
		result int,
		time_finish time.Time,
		user_id int,
		ctx context.Context) error
	UpdateRow(
		table string,
		column string,
		value any,
		indexColumn string,
		index int,
		ctx context.Context) error
	DeleteRow(table string, column string, index int, ctx context.Context) error
	GetStudentsByTeacher(teacher_ID int, limit int, ctx context.Context) ([]models.User, error)
	GetResultByUser(user_id int, limit int, ctx context.Context) ([]models.UserResult, error)
	GetTestByUser(user_id int, ctx context.Context) ([]models.Test, error)
	InitializationRow(
		table_name string,
		colum_name string,
		value any, ctx context.Context) error
	GetTeacherLists(ctx context.Context) ([]models.Teacher, error)
	Close() error
}
type UserCacheRepository interface {
	SetWithTTL(key int, value models.UserState, ttl time.Duration, ctx context.Context) error
	Get(key int, ctx context.Context) (models.UserState, error)
	Exists(key int, ctx context.Context) (bool, error)
	Close() error
}

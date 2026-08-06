package domain

import (
	"TeacherBot/models"
	"time"
)

/*
	type User struct {
		TelegramID   int
		TelegramName string
		FullName     string
		TeacherName  string
	}
*/
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
		teacher_name string) error
	InsertUser(
		telegram_id int,
		telegram_name *string,
		full_name string,
		teacher_ID int) error
	InsertTests(
		subject string,
		level string,
		topic string,
		test string,
		result int,
		time_finish time.Time,
		user_id int) error
	UpdateRow(
		table string,
		column string,
		value any,
		indexColumn string,
		index int) error
	DeleteRow(table string, column string, index int) error
	GetStudentsByTeacher(teacher_ID int, limit int) ([]models.User, error)
	GetResultByUser(user_id int, limit int) ([]models.UserResult, error)
	GetTestByUser(user_id int) ([]models.Test, error)
	InitializationRow(
		table_name string,
		colum_name string,
		value any) error
	GetTeacherLists() ([]models.Teacher, error)
	Close() error
}
type UserCacheRepository interface {
	Ping() error
	SetWithTTL(key int, value models.UserState, ttl time.Duration) error
	Get(key int) (models.UserState, error)
	GetAllKeys() ([]string, error)
	Delete(key int) error
	Exists(key int) (bool, error)
	Expire(key int, ttl time.Duration) error
	FlushDB() error
	Close() error
}

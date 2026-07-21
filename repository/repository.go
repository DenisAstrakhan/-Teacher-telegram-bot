package repository

import (
	"TeacherBot/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func NewUserRepository(ctx context.Context) (domain.UserRepository, error) {
	switch os.Getenv("REPOSITORY_TYPE") {
	case "postgres":
		return newPostgresRepository(ctx)
	case "mysql":
		return newMysqlRepository(ctx)
	default:
		return nil, errors.New("Could not determine database type. REPOSITORY_TYPE environment variable is not set or has invalid value")
	}

}

func newPostgresRepository(ctx context.Context) (domain.UserRepository, error) {
	pgConnString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	//pgConnString := fmt.Sprintf("postgres://%s:%s@bot-postgres:5432/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	conn, err := pgx.Connect(ctx, pgConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &postgresRepository{
		conn: conn,
		ctx:  ctx,
	}, nil
}
func newMysqlRepository(ctx context.Context) (domain.UserRepository, error) {
	ConnString := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?parseTime=true&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DATABASE"))
	//ConnString := fmt.Sprintf("%s:%s@tcp(bot-mysql:3306)/%s?parseTime=true&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DATABASE"))
	conn, err := sql.Open("mysql", ConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &mysqlRepository{
		conn: conn,
		ctx:  ctx,
	}, nil
}

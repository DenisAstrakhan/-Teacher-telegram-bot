package repository

import (
	"TeacherBot/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

func NewUserRepository(ctx context.Context) (domain.UserRepository, error) {
	var userRepository domain.UserRepository
	var err error
	switch os.Getenv("REPOSITORY_TYPE") {
	case "postgres":
		userRepository.DataRepository, err = newPostgresRepository(ctx)
	case "mysql":
		userRepository.DataRepository, err = newMysqlRepository(ctx)
	case "sqlite":
		userRepository.DataRepository, err = newSQLiteRepository(ctx)
	default:
		return domain.UserRepository{}, errors.New("Could not determine database type. REPOSITORY_TYPE environment variable is not set or has invalid value")
	}
	if err != nil {
		return domain.UserRepository{}, fmt.Errorf("Failed to create repository for data storage: %w", err)
	}
	switch os.Getenv("CACHE_REPOSITORY_TYPE") {
	case "redis":
		userRepository.CacheRepository, err = newRedisRepository(ctx)
	default:
		return domain.UserRepository{}, errors.New("Could not determine database type. CACHE_REPOSITORY_TYPE environment variable is not set or has invalid value")
	}
	if err != nil {
		return domain.UserRepository{}, fmt.Errorf("Failed to create cache repository: %w", err)
	}
	return userRepository, nil
}

func newPostgresRepository(ctx context.Context) (domain.UserDataRepository, error) {
	pgConnString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	//pgConnString := fmt.Sprintf("postgres://%s:%s@bot-postgres:5432/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	conn, err := pgx.Connect(queryCtx, pgConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := conn.Ping(queryCtx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &postgresRepository{
		conn: conn,
	}, nil
}
func newMysqlRepository(ctx context.Context) (domain.UserDataRepository, error) {
	ConnString := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?parseTime=true&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DATABASE"))
	//ConnString := fmt.Sprintf("%s:%s@tcp(bot-mysql:3306)/%s?parseTime=true&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DATABASE"))
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	conn, err := sql.Open("mysql", ConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := conn.PingContext(queryCtx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &mysqlRepository{
		conn: conn,
	}, nil
}

func newSQLiteRepository(ctx context.Context) (domain.UserDataRepository, error) {
	conn, err := sql.Open("sqlite3", "./out/sqlitedata/database.db?_foreign_keys=on&cache=shared")
	//conn, err := sql.Open("sqlite3","/data/database.db?_foreign_keys=on&cache=shared")
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite: %w", err)
	}
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// Проверяем подключение
	if err := conn.PingContext(queryCtx); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite: %w", err)
	}

	// Включаем внешние ключи (дополнительная гарантия)
	if _, err := conn.ExecContext(ctx, "PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return &SQLiteRepository{
		conn: conn,
	}, nil
}

func newRedisRepository(ctx context.Context) (redisRepository, error) {
	addr := "localhost:6379"
	opt := &redis.Options{
		Addr:        addr,
		Password:    "",
		DB:          0,
		PoolSize:    10,
		PoolTimeout: 5 * time.Second,
	}

	client := redis.NewClient(opt)

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(queryCtx).Err(); err != nil {
		return redisRepository{}, fmt.Errorf("Failed to connect to Redis: %w", err)
	}
	return redisRepository{
		client: client,
	}, nil
}

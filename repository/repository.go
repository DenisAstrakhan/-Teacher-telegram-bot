package repository

import (
	"TeacherBot/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
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
	switch os.Getenv("LOGGER_REPOSITORY_TYPE") {
	case "clickhouse":
		userRepository.LoggerRepository, err = newClickhouseRepository(ctx)
	default:
		return domain.UserRepository{}, errors.New("Could not determine database type. LOGGER_REPOSITORY_TYPE environment variable is not set or has invalid value")
	}
	if err != nil {
		return domain.UserRepository{}, fmt.Errorf("Failed to create logger repository: %w", err)
	}
	return userRepository, nil
}

func newPostgresRepository(ctx context.Context) (domain.UserDataRepository, error) {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost:5432"
	}
	pgConnString := os.Getenv("POSTGRES_URL")
	if pgConnString == "" {
		pgConnString = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), host, os.Getenv("POSTGRES_DB"))
	}
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	//Создаём конфиг для пула
	config, err := pgxpool.ParseConfig(pgConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to pars config: %w", err)
	}
	//Создаём пул подключений
	pool, err := pgxpool.NewWithConfig(queryCtx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}
	if err := pool.Ping(queryCtx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &postgresRepository{
		conn:    pool,
		timeout: getEnvDuration("POSTGRES_TIMEOUT"),
	}, nil
}
func newMysqlRepository(ctx context.Context) (domain.UserDataRepository, error) {
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "localhost:3306"
	}
	ConnString := os.Getenv("MYSQL_URL")
	if ConnString == "" {
		ConnString = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), host, os.Getenv("MYSQL_DATABASE"))
	}
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
		conn:    conn,
		timeout: getEnvDuration("MYSQL_TIMEOUT"),
	}, nil
}

func newSQLiteRepository(ctx context.Context) (domain.UserDataRepository, error) {
	datadir := os.Getenv("SQLITE_DATA_DIR")
	if datadir == "" {
		datadir = "/data"
	}
	conn, err := sql.Open("sqlite3", fmt.Sprintf("%s/database.db?_foreign_keys=on&cache=shared", datadir))
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
		conn:    conn,
		timeout: getEnvDuration("SQLITE_TIMEOUT"),
	}, nil
}

func newRedisRepository(ctx context.Context) (domain.UserCacheRepository, error) {
	addr := os.Getenv("REDIS_HOST")
	if addr == "" {
		addr = "localhost:6379"
	}
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
		client:  client,
		timeout: getEnvDuration("REDIS_TIMEOUT"),
	}, nil
}

func newClickhouseRepository(ctx context.Context) (domain.UserLoggerRepository, error) {
	host := os.Getenv("CLICKHOUSE_HOST")
	if host == "" {
		host = "localhost:9000"
	}
	db := os.Getenv("CLICKHOUSE_DB")
	if db == "" {
		db = "logs"
	}
	user := os.Getenv("CLICKHOUSE_USER")
	if user == "" {
		user = "default"
	}
	password := os.Getenv("CLICKHOUSE_PASSWORD")

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{host},
		Auth: clickhouse.Auth{
			Database: db,
			Username: user,
			Password: password,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to open clickhouse: %w", err)
	}
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := conn.Ping(queryCtx); err != nil {
		return nil, fmt.Errorf("Failed to connect to Clickhouse: %w", err)
	}
	return clickhouseRepository{
		conn:    conn,
		timeout: getEnvDuration("CLICKHOUSE_TIMEOUT"),
	}, nil
}

// получаем таймаунт из переменной окружения
func getEnvDuration(key string) time.Duration {
	timeoutstring := os.Getenv(key)
	var timeout time.Duration
	var err error
	if timeoutstring == "" {
		timeout = 5 * time.Second
		fmt.Printf("%s not set, using default: %v\n", key, timeout)

	} else {
		timeout, err = time.ParseDuration(timeoutstring)
		if err != nil {
			timeout = 5 * time.Second
			fmt.Printf("Invalid %s '%s', using default 5s: %v\n", key, timeoutstring, err)
		}
	}
	return timeout
}

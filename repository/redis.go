package repository

import (
	"TeacherBot/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisRepository struct {
	client  *redis.Client
	timeout time.Duration
}

func (r redisRepository) SetWithTTL(key int, value models.UserState, ttl time.Duration, ctx context.Context) error {
	// Проверяем что TTL не отрицательный
	if ttl < 0 {
		return errors.New("TTL must be non-negative")
	}
	// Создаем контекст с ограничением времени на операцию
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	//Создаём строку ключа
	keyStr := fmt.Sprintf("user:state:%d", key)
	// Парсим структуру в json
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("Failed to marshal value: %w", err)
	}
	// Выполняем установку значения с TTL
	err = r.client.Set(queryCtx, keyStr, data, ttl).Err()
	if err != nil {
		// Ошибку оборачиваем с указанием ключа
		return fmt.Errorf("Failed to set key %s: %w", keyStr, err)
	}
	return nil
}

func (r redisRepository) Get(key int, ctx context.Context) (models.UserState, error) {
	// Создаем контекст с ограничением времени на операцию
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	//Создаём строку ключа
	keyStr := fmt.Sprintf("user:state:%d", key)
	// Получить значение по ключу
	data, err := r.client.Get(queryCtx, keyStr).Bytes()
	if err == redis.Nil {
		// Ключ не найден – возвращаем кастомную ошибку
		return models.UserState{}, fmt.Errorf("Key not found: %d", key)
	}
	if err != nil {
		// Ошибка другой nature – оборачиваем её
		return models.UserState{}, fmt.Errorf("Failed to get key %d: %w", key, err)
	}
	// Возврат найденного значения
	var state models.UserState
	if err := json.Unmarshal(data, &state); err != nil {
		return models.UserState{}, fmt.Errorf("Failed to unmarshal retrieved data into JSON: %w", err)
	}
	return state, nil
}

func (r redisRepository) Exists(key int, ctx context.Context) (bool, error) {
	queryCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	//Создаём строку ключа
	keyStr := fmt.Sprintf("user:state:%d", key)
	// Выполняем запрос на существование ключа
	exists, err := r.client.Exists(queryCtx, keyStr).Result()
	if err != nil {
		// Ошибку оборачиваем с указанием ключа
		return false, fmt.Errorf("Failed to check key existence %s: %w", keyStr, err)
	}
	return exists > 0, nil
}

func (r redisRepository) Close() error {
	return r.client.Close()
}

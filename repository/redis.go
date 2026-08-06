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
	ctx    context.Context
	client *redis.Client
}

func (r redisRepository) Set(key int, value models.UserState) error {
	// Создаем контекст с ограничением времени на операцию
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	//Создаём строку ключа
	keyStr := fmt.Sprintf("user:state:%d", key)
	// Парсим структуру в json
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %w", err)
	}
	// Выполняем установку значения
	err = r.client.Set(ctx, keyStr, data, 0).Err()
	if err != nil {
		// Ошибку оборачиваем с указанием ключа
		return fmt.Errorf("ошибка при установке ключа %s: %w", keyStr, err)
	}
	return nil
}

func (r redisRepository) SetWithTTL(key int, value models.UserState, ttl time.Duration) error {
	// Проверяем что TTL не отрицательный
	if ttl < 0 {
		return errors.New("Ошибка ttl отрицательное")
	}
	// Создаем контекст с ограничением времени на операцию
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	//Создаём строку ключа
	keyStr := fmt.Sprintf("user:state:%d", key)
	// Парсим структуру в json
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %w", err)
	}
	// Выполняем установку значения с TTL
	err = r.client.Set(ctx, keyStr, data, ttl).Err()
	if err != nil {
		// Ошибку оборачиваем с указанием ключа
		return fmt.Errorf("ошибка при установке ключа %s: %w", keyStr, err)
	}
	return nil
}

func (r redisRepository) Get(key int) (models.UserState, error) {
	// Создаем контекст с ограничением времени на операцию
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	//Создаём строку ключа
	keyStr := fmt.Sprintf("user:state:%d", key)
	// Получить значение по ключу
	data, err := r.client.Get(ctx, keyStr).Bytes()
	if err == redis.Nil {
		// Ключ не найден – возвращаем кастомную ошибку
		return models.UserState{}, fmt.Errorf("ключ не найден: %d", key)
	}
	if err != nil {
		// Ошибка другой nature – оборачиваем её
		return models.UserState{}, fmt.Errorf("ошибка при получении ключа %d: %w", key, err)
	}
	// Возврат найденного значения
	var state models.UserState
	if err := json.Unmarshal(data, &state); err != nil {
		return models.UserState{}, fmt.Errorf("не удалось перевести полученные данные в json. Ошибка: %w", err)
	}
	return state, nil
}

func (r redisRepository) GetAllKeys() ([]string, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	keys, err := r.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении всех ключей: %w", err)
	}
	return keys, nil
}

func (r redisRepository) Delete(key int) error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	//Создаём строку ключа
	keyStr := fmt.Sprintf("user:state:%d", key)
	// Удаление ключа через клиент Redis
	err := r.client.Del(ctx, keyStr).Err()
	if err != nil {
		// Оборачиваем ошибку с указанием ключа
		return fmt.Errorf("ошибка при удалении ключа %s: %w", keyStr, err)
	}
	return nil
}

func (r redisRepository) Exists(key int) (bool, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	//Создаём строку ключа
	keyStr := fmt.Sprintf("user:state:%d", key)
	// Выполняем запрос на существование ключа
	exists, err := r.client.Exists(ctx, keyStr).Result()
	if err != nil {
		// Ошибку оборачиваем с указанием ключа
		return false, fmt.Errorf("ошибка при проверке существования ключа %s: %w", keyStr, err)
	}
	return exists > 0, nil
}

func (r redisRepository) Expire(key int, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	//Создаём строку ключа
	keyStr := fmt.Sprintf("user:state:%d", key)
	err := r.client.Expire(ctx, keyStr, ttl).Err()
	if err != nil {
		return fmt.Errorf("ошибка при установке TTL для ключа %s: %w", keyStr, err)
	}
	return nil
}

func (r redisRepository) FlushDB() error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	err := r.client.FlushDB(ctx).Err()
	if err != nil {
		return fmt.Errorf("ошибка при очистке базы данных: %w", err)
	}
	return nil
}

func (r redisRepository) Ping() error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	err := r.client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("ошибка при проверке соединения с Redis: %w", err)
	}
	return nil
}

func (r redisRepository) Close() error {
	if err := r.client.Close(); err != nil {
		return err
	}
	return nil
}

package domain

import (
	"TeacherBot/models"
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	sensitive "github.com/LuYongwang/go-sensitive-word"
	"github.com/tigusigalpa/gigachat-go"
	"go.uber.org/zap"
)

type BotContext struct {
	UserRepository UserRepository             //Подключение к базе данных
	GigaChat       *gigachat.Client           // клиент подключения к Giga Chat
	UserStates     map[int64]models.UserState //хранилище состояний пользователей
	Subjects       map[string]struct{}        //хранилеще предметов для формирования теста
	Filter         *sensitive.Manager         //Фильтер для фильтрации мата
	Mtx            sync.RWMutex               // для потокобезопасного доступа к UserStates и Giga Chat
}

func NewBotContext(userRepository UserRepository, client *gigachat.Client, filter *sensitive.Manager, logger *zap.Logger) *BotContext {
	subjects, err := newSubjectList(logger)
	if err != nil {
		subjects = make(map[string]struct{})
	}

	return &BotContext{
		UserRepository: userRepository,
		GigaChat:       client,
		UserStates:     make(map[int64]models.UserState),
		Subjects:       subjects,
		Filter:         filter,
		Mtx:            sync.RWMutex{},
	}
}
func (bc *BotContext) SetUserState(userID int64, state models.UserState) {
	bc.Mtx.Lock()
	defer bc.Mtx.Unlock()

	bc.UserStates[userID] = state
}
func (bc *BotContext) GetUserStattes() map[int64]models.UserState {
	bc.Mtx.RLock()
	defer bc.Mtx.RUnlock()
	return bc.UserStates
}
func newSubjectList(logger *zap.Logger) (map[string]struct{}, error) {
	file, err := os.Open("dictionaries/SubjectList.txt")
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open SubjectList.txt file: %v", err))
		return nil, err
	}
	defer file.Close()
	subjectSet := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Удаляем лишние пробелы и приводим к нижнему регистру
		subject := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if subject != "" {
			subjectSet[subject] = struct{}{}
		}
	}
	//Проверяем небыло ли ошибки при сканировании
	if err := scanner.Err(); err != nil {
		logger.Error(fmt.Sprintf("Failed to read SubjectList.txt file: %v", err))
		return nil, err
	}

	return subjectSet, nil
}

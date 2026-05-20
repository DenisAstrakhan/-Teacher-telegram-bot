package models

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/tigusigalpa/gigachat-go"
	"go.uber.org/zap"
)

// Состояние пользователя
type UserState struct {
	CurrentMenu    string              // текущее меню
	Data           map[string]string   // дополнительные данные
	MessageID      int                 //ID сообщения для изменения
	Conversation   []gigachat.Message  // переписка с чатом
	UserAnswers    []string            //ответы пользователя
	CorrectAnswers []string            //правельные ответы
	AllQuestions   []string            //тест 10 вопросов
	UserLastPress  map[int64]time.Time //Хранилище времени последнего нажатия
}

func NewUserState() UserState {
	return UserState{
		CurrentMenu:   "main",
		Data:          make(map[string]string),
		UserLastPress: make(map[int64]time.Time),
		MessageID:     0,
	}
}

type BotContext struct {
	GigaChat   *gigachat.Client    // клиент подключения к Giga Chat
	UserStates map[int64]UserState //хранилище состояний пользователей
	Subjects   map[string]struct{} //хранилеще предметов для формирования теста
	Mtx        sync.RWMutex        // для потокобезопасного доступа к UserStates и Giga Chat
}

func NewBotContext(client *gigachat.Client, logger *zap.Logger) *BotContext {
	subjects, err := newSubjectList(logger)
	if err != nil {
		subjects = make(map[string]struct{})
	}

	return &BotContext{
		GigaChat:   client,
		UserStates: make(map[int64]UserState),
		Subjects:   subjects,
		Mtx:        sync.RWMutex{},
	}
}
func (bc *BotContext) SetUserState(userID int64, state UserState) {
	bc.Mtx.Lock()
	defer bc.Mtx.Unlock()

	bc.UserStates[userID] = state
}
func (bc *BotContext) GetUserStattes() map[int64]UserState {
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

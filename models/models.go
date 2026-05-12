package models

import (
	"sync"
	"time"

	"github.com/tigusigalpa/gigachat-go"
)

// Состояние пользователя
type UserState struct {
	CurrentMenu    string              // текущее меню
	Data           map[string]string   // дополнительные данные
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
	}
}

type BotContext struct {
	GigaChat   *gigachat.Client    // клиент подключения к Giga Chat
	UserStates map[int64]UserState //хранилище состояний пользователей
	Mtx        sync.RWMutex        // для потокобезопасного доступа к UserStates и Giga Chat
}

func NewBotContext(client *gigachat.Client) *BotContext {
	return &BotContext{
		GigaChat:   client,
		UserStates: make(map[int64]UserState),
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

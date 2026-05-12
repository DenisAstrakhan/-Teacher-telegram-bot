package models

import (
	"sync"

	"github.com/tigusigalpa/gigachat-go"
)

// Состояние пользователя
type UserState struct {
	CurrentMenu string            // текущее меню
	Data        map[string]string // дополнительные данные
}

func NewUserState() UserState {
	return UserState{
		CurrentMenu: "main",
		Data:        make(map[string]string),
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

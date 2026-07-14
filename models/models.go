package models

import (
	"time"

	"github.com/tigusigalpa/gigachat-go"
)

// Состояние пользователя
type UserState struct {
	Teacher        *bool              //является ли пользователь учителем
	CurrentMenu    string             // текущее меню
	Data           map[string]string  // дополнительные данные
	MessageID      int                //ID сообщения для изменения
	Conversation   []gigachat.Message // переписка с чатом
	UserAnswers    []string           //ответы пользователя
	CorrectAnswers []string           //правельные ответы
	AllQuestions   []string           //тест 10 вопросов
	UserLastPress  time.Time          //Хранилище времени последнего нажатия
	TeacherLists   []Teacher          //хранилеще списка учителей
}

func NewUserState(teacher *bool) UserState {
	state := UserState{
		Teacher:     teacher,
		CurrentMenu: "main",
		Data:        make(map[string]string),
		MessageID:   0,
	}
	state.Data["subject"] = ""
	state.Data["Topic"] = ""
	state.Data["level"] = ""
	return state
}

type User struct {
	Telegram_id   int
	Telegram_name string
	Full_name     string
}

type UserResult struct {
	Id          int
	Result      int
	Time_finish time.Time
}

type Test struct {
	Subject string
	Level   string
	Topic   string
	Test    string
	Result  int
}

type Teacher struct {
	Telegram_id  int
	Teacher_name string
}

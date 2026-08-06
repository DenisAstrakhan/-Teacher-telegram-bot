package models

import (
	"time"

	"github.com/tigusigalpa/gigachat-go"
)

// Состояние пользователя
type UserState struct {
	Teacher        *bool              `json:"teacher,omitempty"`         //является ли пользователь учителем
	CurrentMenu    string             `json:"current_menu,omitempty"`    //текущее меню
	Data           map[string]string  `json:"data,omitempty"`            //дополнительные данные
	MessageID      int                `json:"message_id,omitempty"`      //ID сообщения для изменения
	Conversation   []gigachat.Message `json:"conversation,omitempty"`    // переписка с чатом
	UserAnswers    []string           `json:"user_answers,omitempty"`    //ответы пользователя
	CorrectAnswers []string           `json:"correct_answers,omitempty"` //правельные ответы
	AllQuestions   []string           `json:"all_questions,omitempty"`   //тест 10 вопросов
	UserLastPress  time.Time          `json:"user_last_press,omitempty"` //Времени последнего нажатия    ??????????????????????????????????
	TeacherLists   []Teacher          `json:"teacher_lists,omitempty"`   //хранилище списка учителей
	StudentList    []User             `json:"student_list,omitempty"`    //Хранилище списка учеников
	TestList       []Test             `json:"test_list,omitempty"`       //Хранилище списка тестов
	Student        User               `json:"student,omitempty"`         //Выбранный студент            ?????????????????????????????????
	TestID         int                `json:"test_id,omitempty"`         //Id выбранного теста
}

type Message struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type User struct {
	Telegram_id   int    `json:"telegram_id,omitempty"` //Изменить на TelegramID
	Telegram_name string `json:"telegram_name,omitempty"`
	Full_name     string `json:"full_name,omitempty"` //Изменить на FullName
}

type UserResult struct {
	Id          int        `json:"id,omitempty"`
	Result      int        `json:"result,omitempty"`
	Time_finish *time.Time `json:"time_finish,omitempty"` // изменить на TimeFinish
}

type Test struct {
	Id      int    `json:"id,omitempty"`
	Subject string `json:"subject,omitempty"`
	Level   string `json:"level,omitempty"`
	Topic   string `json:"topic,omitempty"`
	Test    string `json:"test,omitempty"`
	Result  int    `json:"result,omitempty"`
}

type Teacher struct {
	Telegram_id  int    `json:"telegram_id,omitempty"`  //Изменить на TelegramID
	Teacher_name string `json:"teacher_name,omitempty"` //Изменить на TeacherName
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

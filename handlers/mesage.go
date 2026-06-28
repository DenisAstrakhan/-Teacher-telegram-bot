package handlers

import (
	"TeacherBot/domain"
	gchat "TeacherBot/gigachat"
	"TeacherBot/menu"
	"TeacherBot/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	sensitive "github.com/LuYongwang/go-sensitive-word"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func HandleMessage(logger *zap.Logger, bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext) {
	userID := update.Message.Chat.ID
	text := update.Message.Text
	logger.Info(fmt.Sprintf("User ID - %v: message \"%s\" ", userID, text))
	// Инициализируем состояние пользователя
	userNew, state, err := initializationUserStates(logger, userID, BotContext, bot, update)
	if err != nil {
		logger.Error(fmt.Sprintf("error looking up user in database: %s", err))
		return
	}
	// Обработка команд
	switch text {
	case "/start":
		if !userNew {
			menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!")
			return
		}
		logger.Debug(fmt.Sprintf("Message ID to delete: %v", state.MessageID))
		msgToDelete := tgbotapi.NewDeleteMessage(userID, state.MessageID)
		if _, err := bot.Send(msgToDelete); err != nil {
			logger.Error(fmt.Sprintf("Error sending message: %v", err))
		}
		state.MessageID = 0
		BotContext.SetUserState(userID, state)
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	default:
		//Проверка на пустой ввод
		if len(strings.Fields(text)) == 0 {
			logger.Info(fmt.Sprintf("User %v entered nothing", userID))
			return
		}
		if !validationMessage(text, userID, logger) {
			logger.Debug("Uncorrect input")
			logger.Debug(fmt.Sprintf("Message ID to delete: %v", state.MessageID))
			msgToDelete := tgbotapi.NewDeleteMessage(userID, state.MessageID)
			if _, err := bot.Send(msgToDelete); err != nil {
				logger.Error(fmt.Sprintf("Error sending message: %v", err))
			}
			state.MessageID = 0
			BotContext.SetUserState(userID, state)
			menu.ShowWarningMenu(bot, update, logger, BotContext)
			return
		}
		logger.Debug("Correct input")
		if _, exists := state.Data["score"]; exists {
			// Пользователь проходит интерактивный тест
			logger.Info(fmt.Sprintf("User ID - %v interactive test message: %s ", userID, text))
			gchat.InteractiveTest(bot, update, BotContext, logger)
			return
		}
		if state.CurrentMenu == "setting" && state.Data["subject"] == "" {
			//Пользователь выбирает предмет теста
			logger.Info(fmt.Sprintf("User ID - %v selected subject test: %s ", userID, text))
			if !validationSubject(text, userID, BotContext, logger) {
				msg := tgbotapi.NewMessage(userID, "Попробуйте ещё раз! Указанного предмета нет в согласованном списке")
				if _, err := bot.Send(msg); err != nil {
					logger.Error(fmt.Sprintf("Error sending message: %v", err))
				}
				return
			}
			state.Data["subject"] = text
			BotContext.SetUserState(userID, state)
			msg := tgbotapi.NewMessage(userID, "Напишите тему теста")
			if _, err := bot.Send(msg); err != nil {
				logger.Error(fmt.Sprintf("Error sending message: %v", err))
			}
			return
		}
		if state.CurrentMenu == "setting" && state.Data["Topic"] == "" && state.Data["subject"] != "" {
			//Пользователь выбирает тему теста
			logger.Info(fmt.Sprintf("User ID - %v selected topic test: %s ", userID, text))
			state.Data["Topic"] = text
			state.Data["level"] = "Базовый"
			state.MessageID = 0
			BotContext.SetUserState(userID, state)
			menu.ShowSetingMenu(bot, update, logger, BotContext)
			return
		}
		if _, exists := state.Data["teacher"]; exists {
			//Пользователь вносит учителя в базу данных
			BotContext.UserRepository.InsertTecher(int(userID), &update.Message.From.UserName, text)
		}
		if _, exists := state.Data["student"]; exists {
			//Пользователь вносит ученика в базу данных
			BotContext.UserRepository.InsertUser(int(userID), &update.Message.From.UserName, state.Data["user name"], 12345)
		}
		if _, exists := state.Data["user name"]; exists {
			//Пользователь вводит своё имя
			state.Data["student"] = ""
			state.Data["user name"] = text
			BotContext.SetUserState(userID, state)
		}
		logger.Debug("Simple input")
	}
}
func validationMessage(text string, userID int64, logger *zap.Logger) bool {
	// Инициализация фильтра
	filter, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err != nil {
		logger.Error(fmt.Sprintf("Error creating filter: %v\n", err))
		return false
	}
	// Загрузка словаря Русских ругательств из файла
	err = filter.LoadDictPath("dictionaries/russian-bad-words.txt")
	if err != nil {
		logger.Error(fmt.Sprintf("Error loading dictionary: %v\n", err))
		return false
	}
	//Проверка запрещённых слов
	if filter.IsSensitive(text) {
		logger.Info(fmt.Sprintf("User %v entered forbidden words", userID))
		return false
	}
	return true
}
func validationSubject(text string, userID int64, BotContext *domain.BotContext, logger *zap.Logger) bool {
	BotContext.Mtx.RLock()
	Subjects := BotContext.Subjects
	BotContext.Mtx.RUnlock()
	if _, exist := Subjects[strings.ToLower(text)]; exist {
		logger.Info(fmt.Sprintf("User %v entered a subject from the list", userID))
		return true
	}
	logger.Info(fmt.Sprintf("User %v entered a subject not in the list", userID))
	return false
}
func initializationUserStates(logger *zap.Logger, userID int64, BotContext *domain.BotContext, bot *tgbotapi.BotAPI, update tgbotapi.Update) (bool, models.UserState, error) {
	userStates := BotContext.GetUserStattes()
	state, exists := userStates[userID]
	if !exists {
		if err := BotContext.UserRepository.InitializationRow("bot.teacher", "telegram_id", userID); !errors.Is(err, sql.ErrNoRows) {
			if err == nil {
				teacher := true
				state = createUserSate(&teacher, userID, logger, BotContext)
				return false, state, nil
			}
			return false, models.UserState{}, err
		}
		if err := BotContext.UserRepository.InitializationRow("bot.users", "telegram_id", userID); !errors.Is(err, sql.ErrNoRows) {
			if err == nil {
				teacher := false
				state = createUserSate(&teacher, userID, logger, BotContext)
				return false, state, nil
			}
			return false, models.UserState{}, err
		}
		createUserSate(nil, userID, logger, BotContext)
		menu.ShowWhoAreYouMenu(bot, update, logger, BotContext)
	}
	return true, state, nil
}
func createUserSate(teacher *bool, userID int64, logger *zap.Logger, BotContext *domain.BotContext) models.UserState {
	logger.Info(fmt.Sprintf("Add New User: ID - %v", userID))
	state := models.NewUserState(teacher)
	state.Data["subject"] = ""
	state.Data["Topic"] = ""
	state.Data["level"] = ""
	BotContext.SetUserState(userID, state)
	return state
}

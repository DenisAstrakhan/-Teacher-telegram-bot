package handlers

import (
	"TeacherBot/domain"
	gchat "TeacherBot/gigachat"
	"TeacherBot/menu"
	"TeacherBot/models"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func HandleMessage(logger *zap.Logger, bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext) {
	userID := update.Message.Chat.ID
	text := update.Message.Text
	logger.Info(fmt.Sprintf("User ID - %v: message \"%s\" ", userID, text))
	// Инициализируем состояние пользователя
	userNew, state, err := initializationUserStates(logger, userID, BotContext, bot, update)
	logger.Debug(fmt.Sprintf("userNew - %t", userNew))
	if userNew && err == nil {
		//Есл пользователь не найден спрашиваем "Кто он?"
		menu.ShowWhoAreYouMenu(bot, update, logger, BotContext)
	}
	if err != nil {
		logger.Error(fmt.Sprintf("error looking up user in database: %s", err))
		return
	}
	// Обработка команд
	switch text {
	case "/start":
		if userNew {
			return
		}
		if state.Teacher != nil && *state.Teacher {
			menu.ShowTeacherMenu(bot, update, logger, BotContext)
			return
		}
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
	default:
		//Проверка на пустой ввод
		if len(strings.Fields(text)) == 0 {
			logger.Info(fmt.Sprintf("User %v entered nothing", userID))
			return
		}
		//Проверка на коректность ввода
		if !validationMessage(text, userID, BotContext, logger) {
			if state.Teacher != nil && *state.Teacher {
				//Учитель написал непреличное слова пока закрое на это глаза
				logger.Warn("Учитель выражается не корректно")
				return
			}
			logger.Debug("Uncorrect input")
			logger.Debug(fmt.Sprintf("Message ID to delete: %v", state.MessageID))
			msgToDelete := tgbotapi.NewDeleteMessage(userID, state.MessageID)
			if _, err := bot.Request(msgToDelete); err != nil {
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
		if state.CurrentMenu == "teacher" {
			//Пользователь вносит учителя в базу данных
			err := BotContext.UserRepository.InsertTecher(int(userID), &update.Message.From.UserName, text)
			if err != nil {
				logger.Error(fmt.Sprintf("Ошибка при добавлении учителя в базу данных: %v", err))
			}
			logger.Info(fmt.Sprintf("Пользователь ID-%d добавлен в базу данных.", userID))
			menu.ShowTeacherMenu(bot, update, logger, BotContext)
			return
		}

		if _, exists := state.Data["user name"]; exists {
			//Пользователь вводит своё имя
			state := menu.ShowTecherList(bot, update, logger, BotContext)
			state.Data["user name"] = text
			state.Data["student"] = ""
			BotContext.SetUserState(userID, state)
			return
		}
		//Проверяем введено ли число
		choice, err := strconv.Atoi(text)
		if err != nil {
			//Введено не число
			return
		}
		if _, exists := state.Data["resoult edit"]; !exists {
			//Введено просто число
			logger.Debug("resoult edit не сущестует")
			return
		}
		//Учитель исправляет результат теста
		if choice <= 100 && choice >= 0 {
			if err := BotContext.UserRepository.UpdateRow("tests", "result", choice, "id", state.TestID); err != nil {
				logger.Error(fmt.Sprintf("Ошибка при попытке редактировать результат теста в БД: %v", err))
				msg := tgbotapi.NewMessage(userID, "Произошла ошибка при внесении изменений в БД. Попробуйте ещё раз")
				if _, err := bot.Send(msg); err != nil {
					logger.Error(fmt.Sprintf("Error sending mesage: %v", err))
					return
				}
				return
			}
			delete(state.Data, "resoult edit")
			logger.Info("Отценка исправлена в базе данных")
			msg := tgbotapi.NewMessage(userID, "Отценка исправлена")
			if _, err := bot.Send(msg); err != nil {
				logger.Error(fmt.Sprintf("Error sending mesage: %v", err))
				return
			}
			return
		}
		msg := tgbotapi.NewMessage(userID, "Введённое число не укладывается в диапазон от 0 до 100. Попробуйте ещё раз.")
		if _, err := bot.Send(msg); err != nil {
			logger.Error(fmt.Sprintf("Error sending mesage: %v", err))
			return
		}
		logger.Debug("Simple input")
	}
}
func validationMessage(text string, userID int64, BotContext *domain.BotContext, logger *zap.Logger) bool {
	//Проверка запрещённых слов
	if BotContext.Filter.IsSensitive(text) {
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
		//Пользователя нет в программе
		if err := BotContext.UserRepository.InitializationRow("teacher", "telegram_id", userID); !errors.Is(err, sql.ErrNoRows) {
			if err == nil {
				//Пользователь есть в базе как учитель
				teacher := true
				state = models.NewUserState(&teacher)
				BotContext.SetUserState(userID, state)
				logger.Info(fmt.Sprintf("Teacher fetched from database, ID: - %d", userID))
				return false, state, nil
			}
			return true, models.UserState{}, err
		}
		if err := BotContext.UserRepository.InitializationRow("users", "telegram_id", userID); !errors.Is(err, sql.ErrNoRows) {
			if err == nil {
				//Пользователь есть в базе как ученик
				teacher := false
				state = models.NewUserState(&teacher)
				BotContext.SetUserState(userID, state)
				logger.Info(fmt.Sprintf("User fetched from database, ID: -  %d", userID))
				return false, state, nil
			}
			return true, models.UserState{}, err
		}
		return true, models.UserState{}, nil
	}
	return false, state, nil
}

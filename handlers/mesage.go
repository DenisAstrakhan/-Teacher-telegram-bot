package handlers

import (
	"TeacherBot/domain"
	gchat "TeacherBot/gigachat"
	"TeacherBot/menu"
	"TeacherBot/models"
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func HandleMessage(logger *zap.Logger, bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, ctx context.Context) {
	userID := update.Message.Chat.ID
	text := update.Message.Text
	logger.Sugar().Infof("User %d: message \"%s\" ", userID, text)
	// Инициализируем состояние пользователя
	userNew, state, err := initializationUserStates(logger, userID, BotContext, bot, update, ctx)
	logger.Sugar().Debugf("userNew - %t", userNew)
	if userNew && err == nil {
		//Есл пользователь не найден спрашиваем "Кто он?"
		menu.ShowWhoAreYouMenu(bot, update, logger, BotContext)
	}
	if err != nil {
		logger.Error("Error looking up user in database: %w", zap.Error(err))
		return
	}
	// Обработка команд
	switch text {
	case "/start":
		if userNew {
			return
		}
		if state.Teacher != nil && *state.Teacher {
			menu.ShowTeacherMenu(bot, update, logger, BotContext, ctx)
			return
		}
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
	default:
		//Проверка на пустой ввод
		if len(strings.Fields(text)) == 0 {
			logger.Sugar().Infof("User %d entered nothing", userID)
			return
		}
		//Проверка на коректность ввода
		if !validationMessage(text, userID, BotContext, logger) {
			if state.Teacher != nil && *state.Teacher {
				//Учитель написал непреличное слова пока закрое на это глаза
				logger.Warn("Teacher used inappropriate language")
				return
			}
			logger.Debug("Uncorrect input")
			logger.Sugar().Debugf("Message ID to delete: %v", state.MessageID)
			msgToDelete := tgbotapi.NewDeleteMessage(userID, state.MessageID)
			if _, err := bot.Request(msgToDelete); err != nil {
				logger.Error("Error sending message: %w", zap.Error(err))
			}
			state.MessageID = 0
			menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
			menu.ShowWarningMenu(bot, update, logger, BotContext, ctx)
			return
		}
		logger.Debug("Correct input")
		if _, exists := state.Data["score"]; exists {
			// Пользователь проходит интерактивный тест
			logger.Sugar().Infof("User %d interactive test message: %s ", userID, text)
			gchat.InteractiveTest(bot, update, BotContext, logger, ctx)
			return
		}
		if state.CurrentMenu == "setting" && state.Data["subject"] == "" {
			//Пользователь выбирает предмет теста
			logger.Sugar().Infof("User %d selected subject test: %s ", userID, text)
			if !validationSubject(text, userID, BotContext, logger) {
				msg := tgbotapi.NewMessage(userID, "Попробуйте ещё раз! Указанного предмета нет в согласованном списке")
				if _, err := bot.Send(msg); err != nil {
					logger.Error("Error sending message: %w", zap.Error(err))
				}
				return
			}
			state.Data["subject"] = text
			menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
			msg := tgbotapi.NewMessage(userID, "Напишите тему теста")
			if _, err := bot.Send(msg); err != nil {
				logger.Error("Error sending message: %w", zap.Error(err))
			}
			return
		}
		if state.CurrentMenu == "setting" && state.Data["Topic"] == "" && state.Data["subject"] != "" {
			//Пользователь выбирает тему теста
			logger.Sugar().Infof("User %d selected topic test: %s ", userID, text)
			state.Data["Topic"] = text
			state.Data["level"] = "Базовый"
			state.MessageID = 0
			menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
			menu.ShowSetingMenu(bot, update, logger, BotContext, ctx)
			return
		}
		if state.CurrentMenu == "teacher" {
			//Пользователь вносит учителя в базу данных
			err := BotContext.UserRepository.DataRepository.InsertTecher(int(userID), &update.Message.From.UserName, text, ctx)
			if err != nil {
				logger.Error("Failed to add teacher to the database: %w", zap.Error(err))
				return
			}
			logger.Sugar().Infof("User %d added to the database", userID)
			menu.ShowTeacherMenu(bot, update, logger, BotContext, ctx)
			return
		}

		if _, exists := state.Data["user name"]; exists {
			//Пользователь вводит своё имя
			state := menu.ShowTecherList(bot, update, logger, BotContext, ctx)
			state.Data["user name"] = text
			state.Data["student"] = ""
			menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
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
			logger.Debug("Map key 'resoult edit' does not exist")
			return
		}
		//Учитель исправляет результат теста
		if choice <= 100 && choice >= 0 {
			if err := BotContext.UserRepository.DataRepository.UpdateRow("tests", "result", choice, "id", state.TestID, ctx); err != nil {
				logger.Error("Failed to update test result in the database: %w", zap.Error(err))
				msg := tgbotapi.NewMessage(userID, "Произошла ошибка при внесении изменений в БД. Попробуйте ещё раз")
				if _, err := bot.Send(msg); err != nil {
					logger.Error("Error sending mesage: %w", zap.Error(err))
					return
				}
				return
			}
			delete(state.Data, "resoult edit")
			logger.Info("Test result updated in the database")
			msg := tgbotapi.NewMessage(userID, "Отценка исправлена")
			if _, err := bot.Send(msg); err != nil {
				logger.Error("Error sending mesage: %w", zap.Error(err))
				return
			}
			return
		}
		msg := tgbotapi.NewMessage(userID, "Введённое число не укладывается в диапазон от 0 до 100. Попробуйте ещё раз.")
		if _, err := bot.Send(msg); err != nil {
			logger.Error("Error sending mesage: %w", zap.Error(err))
			return
		}
		logger.Debug("Simple input")
	}
}
func validationMessage(text string, userID int64, BotContext *domain.BotContext, logger *zap.Logger) bool {
	//Проверка запрещённых слов
	if BotContext.Filter.IsSensitive(text) {
		logger.Sugar().Infof("User %d entered forbidden words", userID)
		return false
	}
	return true
}
func validationSubject(text string, userID int64, BotContext *domain.BotContext, logger *zap.Logger) bool {
	BotContext.Mtx.RLock()
	Subjects := BotContext.Subjects
	BotContext.Mtx.RUnlock()
	if _, exist := Subjects[strings.ToLower(text)]; exist {
		logger.Sugar().Infof("User %d entered a subject from the list", userID)
		return true
	}
	logger.Sugar().Infof("User %d entered a subject not in the list", userID)
	return false
}
func initializationUserStates(logger *zap.Logger, userID int64, BotContext *domain.BotContext, bot *tgbotapi.BotAPI, update tgbotapi.Update, ctx context.Context) (bool, models.UserState, error) {
	exists, err := BotContext.UserRepository.CacheRepository.Exists(int(userID), ctx)
	if err != nil {
		logger.Error("Failed to retrieve data from cache: %w ", zap.Error(err))
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		return true, models.UserState{}, err
	}
	if !exists {
		//Пользователя нет в программе
		if err := BotContext.UserRepository.DataRepository.InitializationRow("teacher", "telegram_id", userID, ctx); !errors.Is(err, sql.ErrNoRows) {
			if err == nil {
				//Пользователь есть в базе как учитель
				teacher := true
				state := models.NewUserState(&teacher)
				menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
				logger.Sugar().Infof("Teacher fetched from database, ID: - %d", userID)
				return false, state, nil
			}
			return true, models.UserState{}, err
		}
		if err := BotContext.UserRepository.DataRepository.InitializationRow("users", "telegram_id", userID, ctx); !errors.Is(err, sql.ErrNoRows) {
			if err == nil {
				//Пользователь есть в базе как ученик
				teacher := false
				state := models.NewUserState(&teacher)
				menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
				logger.Sugar().Infof("User fetched from database, ID: -  %d", userID)
				return false, state, nil
			}
			return true, models.UserState{}, err
		}
		return true, models.UserState{}, nil
	}
	state, err := BotContext.UserRepository.CacheRepository.Get(int(userID), ctx)
	if err != nil {
		logger.Error("Failed to retrieve data from cache: %w", zap.Error(err))
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		return true, models.UserState{}, err
	}
	return false, state, nil
}

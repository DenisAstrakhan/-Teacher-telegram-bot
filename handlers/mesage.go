package handlers

import (
	gchat "TeacherBot/gigachat"
	"TeacherBot/menu"
	"TeacherBot/models"
	"fmt"
	"strings"

	sensitive "github.com/LuYongwang/go-sensitive-word"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func HandleMessage(logger *zap.Logger, bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext) {
	userID := update.Message.Chat.ID
	text := update.Message.Text
	logger.Info(fmt.Sprintf("User ID - %v: message \"%s\" ", userID, text))
	// Инициализируем состояние пользователя
	userStates := BotContext.GetUserStattes()
	state, exists := userStates[userID]
	if !exists {
		logger.Info(fmt.Sprintf("Add New User: ID - %v", userID))
		state = models.NewUserState()
		state.Data["subject"] = ""
		state.Data["Topic"] = ""
		state.Data["level"] = ""
		BotContext.SetUserState(userID, state)
	}
	// Обработка команд
	switch text {
	case "/start":
		menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!")
		return
	default:
		if !validationMessage(text, userID, logger) {
			logger.Debug("Некоректный ввод")
			return
		}
		logger.Debug("Коректный ввод")
		if _, exists := state.Data["score"]; exists {
			// Пользователь проходит интерактивный тест
			logger.Info(fmt.Sprintf("User ID - %v interactive test message: %s ", userID, text))
			gchat.InteractiveTest(bot, update, BotContext, logger)
			return
		}
		if state.Data["subject"] == "" {
			//Пользователь выбирает предмет теста
			logger.Info(fmt.Sprintf("User ID - %v selected subject test: %s ", userID, text))
			state.Data["subject"] = text
			BotContext.SetUserState(userID, state)
			msg := tgbotapi.NewMessage(userID, "Напишите тему теста")
			bot.Send(msg)
			return
		}
		if state.Data["Topic"] == "" && state.Data["subject"] != "" {
			//Пользователь выбирает тему теста
			logger.Info(fmt.Sprintf("User ID - %v selected topic test: %s ", userID, text))
			state.Data["Topic"] = text
			state.Data["level"] = "Базовый"
			state.MessageID = 0
			BotContext.SetUserState(userID, state)
			menu.ShowSetingMenu(bot, update, logger, BotContext)
		}
	}
}
func validationMessage(text string, userID int64, logger *zap.Logger) bool {
	// Инициализация фильтра
	filter, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка создания фильтра: %v\n", err))
		return false
	}
	//Проверка на пустой ввод
	if len(strings.Fields(text)) == 0 {
		logger.Info(fmt.Sprintf("Пользователь %v ничего не ввёл", userID))
		return false
	}
	// Загрузка словаря Русских ругательств из файла
	err = filter.LoadDictPath("handlers/russian-bad-words.txt")
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка загрузки словаря: %v\n", err))
		return false
	}
	//Проверка запрещённых слов
	if filter.IsSensitive(text) {
		logger.Info(fmt.Sprintf("Пользователь %v ввёл запрещённые слова", userID))
		return false
	}
	return true
}

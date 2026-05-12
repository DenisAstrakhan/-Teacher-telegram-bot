package handlers

import (
	"TeacherBot/menu"
	"TeacherBot/models"
	"fmt"
	"strings"

	sensitive "github.com/LuYongwang/go-sensitive-word"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func HandleMesage(logger *zap.Logger, bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext) {
	userID := update.Message.Chat.ID
	text := update.Message.Text
	logger.Info(fmt.Sprintf("User ID - %v: message \"%s\" ", userID, text))
	// Инициализируем состояние пользователя
	userStates := BotContext.GetUserStattes()
	state, exists := userStates[userID]
	if !exists {
		logger.Info(fmt.Sprintf("Add New User: ID - %v", userID))
		state = models.UserState{
			CurrentMenu: "main",
			Data:        make(map[string]string),
		}
		state.Data["subject"] = ""
		state.Data["Topic"] = ""
		state.Data["level"] = ""
	}
	// Обработка команд
	switch text {
	case "/start":
		menu.ShowStartMenu(bot, update, logger)
		return
	default:
		if !validationMessage(text, userID, logger) {
			logger.Debug("Некоректный ввод")
			return
		}
		logger.Debug("Коректный ввод")
		return
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

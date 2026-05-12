package handlers

import (
	"TeacherBot/menu"
	"TeacherBot/models"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func HandleCallback(logger *zap.Logger, bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext) {
	// Отвечаем на callback (убираем "часики")
	bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
	userStates := BotContext.UserStates
	userID := update.CallbackQuery.From.ID
	data := update.CallbackQuery.Data
	logger.Info(fmt.Sprintf("User ID - %v: press \"%s\" ", userID, data))
	// Инициализируем состояние пользователя
	if _, exists := userStates[userID]; !exists {
		userStates[userID] = models.NewUserState()
		logger.Info(fmt.Sprintf("New User: ID - %v, Name: %s", userID, update.CallbackQuery.From.FirstName))
	}
	state := userStates[userID]
	// Обработка callback данных
	switch data {
	case "simple":
		state.CurrentMenu = "simple"
		state.Data["test"] = "simple"
		menu.ShowLevelMenu(bot, update, logger)
	case "interactive":
		state.CurrentMenu = "interactive"
		state.Data["test"] = "interactive"
		menu.ShowLevelMenu(bot, update, logger)
	case "settings":

	case "back":
		goBack(bot, update, BotContext, logger)
		return
	case "Beginner":
		state.CurrentMenu = "Beginner"
		state.Data["level"] = "Beginner"
	case "Intermediate":
		state.CurrentMenu = "Intermediate"
		state.Data["level"] = "Intermediate"
	case "Advanc":
		state.CurrentMenu = "Advanc"
		state.Data["level"] = "Advanc"
	case "Topic1":
		switch state.Data["level"] {
		case "Beginner":
			state.Data["Topic"] = "Present Simple & Present Continuous (базовое сравнение)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		case "Intermediate":
			state.Data["Topic"] = "Present Perfect vs. Past Simple"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		case "Advanc":
			state.Data["Topic"] = "Инверсия (Never have I seen... / Not only did he...)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		default:
			logger.Warn(fmt.Sprintf("User ID - %v: Failed to distribute Topic1 across levels ", userID))
			return
		}
	case "Topic2":
		switch state.Data["level"] {
		case "Beginner":
			state.Data["Topic"] = "There is / There are + предлоги места"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		case "Intermediate":
			state.Data["Topic"] = "Условные предложения (Conditionals: 0, 1, 2 типы)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		case "Advanc":
			state.Data["Topic"] = "Смешанные условные предложения (Mixed Conditionals)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		default:
			logger.Warn(fmt.Sprintf("User ID - %v: Failed to distribute Topic2 across levels ", userID))
			return
		}
	case "Topic3":
		switch state.Data["level"] {
		case "Beginner":
			state.Data["Topic"] = "Модальные глаголы (can / must)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		case "Intermediate":
			state.Data["Topic"] = "Пассивный залог (Present / Past Simple Passive)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		case "Advanc":
			state.Data["Topic"] = "Сослагательное наклонение (I suggest that he go / It’s crucial that she be)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		default:
			logger.Warn(fmt.Sprintf("User ID - %v: Failed to distribute Topic3 across levels ", userID))
			return
		}
	case "Topic4":
		switch state.Data["level"] {
		case "Beginner":
			state.Data["Topic"] = "Простое прошедшее время (правильные и топ-неправильных глаголов)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		case "Intermediate":
			state.Data["Topic"] = "Косвенная речь (Reported Speech)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		case "Advanc":
			state.Data["Topic"] = "Эллипсис и замена (So do I / Neither can she / He does)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		default:
			logger.Warn(fmt.Sprintf("User ID - %v: Failed to distribute Topic4 across levels ", userID))
			return
		}
	case "Topic5":
		switch state.Data["level"] {
		case "Beginner":
			state.Data["Topic"] = "Базовая лексика: семья, еда, дом, время на часах"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		case "Intermediate":
			state.Data["Topic"] = "Фразовые глаголы (look after, give up, run out of и др.)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		case "Advanc":
			state.Data["Topic"] = "Расширенная идиоматика и стилистическая синонимия (разница между ask / inquire / demand)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			return
		default:
			logger.Warn(fmt.Sprintf("User ID - %v: Failed to distribute Topic5 across levels ", userID))
			return
		}
	default:
		logger.Info(fmt.Sprintf("User ID - %v: Failed to process callback", userID))
		menu.ShowStartMenu(bot, update, logger)
	}
	BotContext.SetUserState(userID, state)
}
func goBack(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext, logger *zap.Logger) {
	userID := update.CallbackQuery.From.ID
	userStates := BotContext.UserStates
	state := userStates[userID]
	switch state.CurrentMenu {
	case "simple", "interactive":
		state.CurrentMenu = "main"
		menu.ShowStartMenu(bot, update, logger)
	case "Beginner", "Intermediate", "Advanc":
		state.CurrentMenu = state.Data["test"]
		menu.ShowLevelMenu(bot, update, logger)
	default:
		menu.ShowStartMenu(bot, update, logger)
		state.CurrentMenu = "main"
	}

	BotContext.SetUserState(userID, state)
}

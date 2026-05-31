package handlers

import (
	gchat "TeacherBot/gigachat"
	"TeacherBot/menu"
	"TeacherBot/models"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func HandleCallback(logger *zap.Logger, bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext) {
	// Отвечаем на callback (убираем "часики")
	bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
	userID := update.CallbackQuery.From.ID
	data := update.CallbackQuery.Data
	logger.Info(fmt.Sprintf("User ID - %v: press \"%s\" ", userID, data))
	// Инициализируем состояние пользователя
	exists, state := initializationUserStates(logger, userID, BotContext)
	if !exists {
		logger.Info(fmt.Sprintf("User %v, Name: %s, not found in system", userID, update.CallbackQuery.From.FirstName))
		msg := tgbotapi.NewMessage(userID, "После перезапуска ваши данные в системе были утеряны. Начните с команды \"/start\"")
		if _, err := bot.Send(msg); err != nil {
			logger.Error(fmt.Sprintf("Error sending message: %v", err))
		}
		return
	}
	//Защита от повторного нажатий
	BotContext.Mtx.Lock()
	if time.Since(state.UserLastPress[userID]) < 1000*time.Millisecond {
		logger.Warn(fmt.Sprintf("User ID - %v: press again", userID))
		BotContext.Mtx.Unlock()
		return
	}
	state.UserLastPress[userID] = time.Now()
	BotContext.Mtx.Unlock()
	// Обработка callback данных
	switch data {
	case "simple":
		state.CurrentMenu = "simple"
		state.Data["test"] = "simple"
		if state.Data["subject"] != "" && state.Data["Topic"] != "" && state.Data["level"] != "" {
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
			return
		}
		menu.ShowLevelMenu(bot, update, logger, BotContext)
	case "interactive":
		state.CurrentMenu = "interactive"
		state.Data["test"] = "interactive"
		if state.Data["subject"] != "" && state.Data["Topic"] != "" && state.Data["level"] != "" {
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
			return
		}
		menu.ShowLevelMenu(bot, update, logger, BotContext)
	case "settings":
		gchat.SelectSubject(bot, update, BotContext, logger)
		state.CurrentMenu = "setting"
	case "back":
		goBack(bot, update, BotContext, logger)
		return
	case "sorry":
		menu.ReturnStartMenu(bot, update, BotContext, logger, "Попробуй всё заново")
		return
	case "Beginner":
		state.CurrentMenu = "Beginner"
		state.Data["level"] = "Beginner"
		menu.ShowBeginnerMenu(bot, update, logger, BotContext)
	case "Intermediate":
		state.CurrentMenu = "Intermediate"
		state.Data["level"] = "Intermediate"
		menu.ShowIntermediateMenu(bot, update, logger, BotContext)
	case "Advanc":
		state.CurrentMenu = "Advanc"
		state.Data["level"] = "Advanc"
		menu.ShowAdvancMenu(bot, update, logger, BotContext)
	case "Topic1":
		switch state.Data["level"] {
		case "Beginner":
			state.Data["Topic"] = "Present Simple & Present Continuous (базовое сравнение)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
			return
		case "Intermediate":
			state.Data["Topic"] = "Present Perfect vs. Past Simple"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
			return
		case "Advanc":
			state.Data["Topic"] = "Инверсия (Never have I seen... / Not only did he...)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
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
			gchat.StartTest(bot, update, BotContext, logger)
			return
		case "Intermediate":
			state.Data["Topic"] = "Условные предложения (Conditionals: 0, 1, 2 типы)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
			return
		case "Advanc":
			state.Data["Topic"] = "Смешанные условные предложения (Mixed Conditionals)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
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
			gchat.StartTest(bot, update, BotContext, logger)
			return
		case "Intermediate":
			state.Data["Topic"] = "Пассивный залог (Present / Past Simple Passive)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
			return
		case "Advanc":
			state.Data["Topic"] = "Сослагательное наклонение (I suggest that he go / It’s crucial that she be)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
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
			gchat.StartTest(bot, update, BotContext, logger)
			return
		case "Intermediate":
			state.Data["Topic"] = "Косвенная речь (Reported Speech)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
			return
		case "Advanc":
			state.Data["Topic"] = "Эллипсис и замена (So do I / Neither can she / He does)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
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
			gchat.StartTest(bot, update, BotContext, logger)
			return
		case "Intermediate":
			state.Data["Topic"] = "Фразовые глаголы (look after, give up, run out of и др.)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
			return
		case "Advanc":
			state.Data["Topic"] = "Расширенная идиоматика и стилистическая синонимия (разница между ask / inquire / demand)"
			state.Data["subject"] = "Английский язык"
			BotContext.SetUserState(userID, state)
			gchat.StartTest(bot, update, BotContext, logger)
			return
		default:
			logger.Warn(fmt.Sprintf("User ID - %v: Failed to distribute Topic5 across levels ", userID))
			return
		}
	case "A":
		state.UserAnswers = append(state.UserAnswers, "A")
		BotContext.SetUserState(userID, state)
		gchat.SimpleTest(bot, update, BotContext, logger)
		return
	case "B":
		state.UserAnswers = append(state.UserAnswers, "B")
		BotContext.SetUserState(userID, state)
		gchat.SimpleTest(bot, update, BotContext, logger)
		return
	case "C":
		state.UserAnswers = append(state.UserAnswers, "C")
		BotContext.SetUserState(userID, state)
		gchat.SimpleTest(bot, update, BotContext, logger)
		return
	case "D":
		state.UserAnswers = append(state.UserAnswers, "D")
		BotContext.SetUserState(userID, state)
		gchat.SimpleTest(bot, update, BotContext, logger)
		return
	default:
		logger.Info(fmt.Sprintf("User ID - %v: Failed to process callback", userID))
		menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!")
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
		menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!")
	case "Beginner", "Intermediate", "Advanc":
		state.CurrentMenu = state.Data["test"]
		menu.ShowLevelMenu(bot, update, logger, BotContext)
	case "setting":
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	default:
		menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!")
		state.CurrentMenu = "main"
	}

	BotContext.SetUserState(userID, state)
}

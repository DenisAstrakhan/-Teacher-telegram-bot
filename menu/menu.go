package menu

import (
	"TeacherBot/domain"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func ShowStartMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext, Caption string) {
	// Создаем инлайн клавиатуру
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Простой тест", "simple"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Интерактивный тест", "interactive"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "settings"), //tgbotapi.NewInlineKeyboardButtonURL
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🌐 Repositorie", "https://github.com/DenisAstrakhan/-Teacher-telegram-bot"),
		),
	)
	var chatID int64
	if update.Message == nil {
		chatID = update.CallbackQuery.From.ID
	} else {
		chatID = update.Message.Chat.ID
	}
	userStates := BotContext.GetUserStattes()
	state := userStates[chatID]
	if state.MessageID == 0 {
		sendMenu(bot, update, Caption, keyboard, logger, BotContext, "Image/start.jpg")
	} else {
		editMenu(bot, update, Caption, keyboard, logger, BotContext)
	}
}

func ShowLevelMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext) {
	//level: Beginner (A1-A2) - Новичок, Intermediate (B1-B2) - Средний, Advanced (C1-C2) - Продвинутый,
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Новичёк", "Beginner"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Средний", "Intermediate"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Продвинутый", "Advanc"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back"),
		),
	)
	editMenu(bot, update, "Выберите сложность", keyboard, logger, BotContext)
}
func ShowBeginnerMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext) {
	/*
	   Topic1: Present Simple & Present Continuous (базовое сравнение)
	   Topic2: There is / There are + предлоги места
	   Topic3: Модальные глаголы (can / must)
	   Topic4: Простое прошедшее время (правильные и топ-неправильных глаголов)
	   Topic5: Базовая лексика: семья, еда, дом, время на часах
	*/

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Present Simple & Present Continuous", "Topic1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 СThere is / There are", "Topic2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Модальные глаголы (can / must)", "Topic3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Простое прошедшее время", "Topic4"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝  Базовая лексика", "Topic5"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back"),
		),
	)
	editMenu(bot, update, "Выберите тему", keyboard, logger, BotContext)
}
func ShowIntermediateMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext) {
	/*
		Topic1: Present Perfect vs. Past Simple
		Topic2: Условные предложения (Conditionals: 0, 1, 2 типы)
		Topic3: Пассивный залог (Present / Past Simple Passive)
		Topic4: Косвенная речь (Reported Speech)
		Topic5: Фразовые глаголы (look after, give up, run out of и др.)
	*/
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Present Perfect vs. Past Simple", "Topic1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Условные предложения", "Topic2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Пассивный залог (can / must)", "Topic3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Косвенная речь", "Topic4"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Фразовые глаголы", "Topic5"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back"),
		),
	)
	editMenu(bot, update, "Выберите тему", keyboard, logger, BotContext)
}

func ShowAdvancMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext) {
	/*
		Topic1: Инверсия (Never have I seen... / Not only did he...)
		Topic2: Смешанные условные предложения (Mixed Conditionals)
		Topic3: Сослагательное наклонение (I suggest that he go / It’s crucial that she be)
		Topic4: Эллипсис и замена (So do I / Neither can she / He does)
		Topic5: Расширенная идиоматика и стилистическая синонимия (разница между ask / inquire / demand)

	*/
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Инверсия", "Topic1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Смешанные условные предложения", "Topic2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Сослагательное наклонение", "Topic3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Эллипсис и замена", "Topic4"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Расширенная идиоматика и стилистическая синонимия", "Topic5"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back"),
		),
	)
	editMenu(bot, update, "Выберите тему", keyboard, logger, BotContext)
}
func ShowTestMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, Caption string, logger *zap.Logger, BotContext *domain.BotContext) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("A", "A"),
			tgbotapi.NewInlineKeyboardButtonData("B", "B"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("C", "C"),
			tgbotapi.NewInlineKeyboardButtonData("D", "D"),
		),
	)
	editMenu(bot, update, Caption, keyboard, logger, BotContext)
}
func ShowSetingMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Простой тест", "simple"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Интерактивный тест", "interactive"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back"),
		),
	)
	sendMenu(bot, update, "Выберите тип теста", keyboard, logger, BotContext, "Image/start.jpg")
}
func ShowWarningMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Простите пожалуйста я так больше не буду", "sorry"),
		),
	)
	sendMenu(bot, update, "Ненормативная лексика! За тобой уже выехали.", keyboard, logger, BotContext, "Image/warning.jpg")
}
func ShowWhoAreYouMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext) {
	logger.Debug("start ShowWhoAreYouMenu")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Учитель", "teacher"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ученик", "student"),
		),
	)
	var chatID int64
	if update.Message == nil {
		chatID = update.CallbackQuery.From.ID
	} else {
		chatID = update.Message.Chat.ID
	}
	photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("Image/WhoAreYou.jpeg"))
	photoMsg.Caption = "Кто ты, воин?"
	photoMsg.ReplyMarkup = keyboard
	_, err := bot.Send(photoMsg)
	if err != nil {
		//Не удалось отправить сообщение
		logger.Error(fmt.Sprintf("Error send photo: %v", err))
		return
	}
}
func sendMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, Caption string, keyboard tgbotapi.InlineKeyboardMarkup, logger *zap.Logger, BotContext *domain.BotContext, imageName string) {

	var chatID int64
	if update.Message == nil {
		chatID = update.CallbackQuery.From.ID
	} else {
		chatID = update.Message.Chat.ID
	}
	userStates := BotContext.GetUserStattes()
	state := userStates[chatID]
	photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(imageName))
	photoMsg.Caption = Caption
	photoMsg.ReplyMarkup = keyboard
	sendMessage, err := bot.Send(photoMsg)
	if err != nil {
		//Не удалось отправить сообщение
		logger.Error(fmt.Sprintf("Error send photo: %v", err))
		return
	}
	state.MessageID = sendMessage.MessageID
	BotContext.SetUserState(chatID, state)

}

func editMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, Caption string, keyboard tgbotapi.InlineKeyboardMarkup, logger *zap.Logger, BotContext *domain.BotContext) {
	var chatID int64
	if update.Message == nil {
		chatID = update.CallbackQuery.From.ID
	} else {
		chatID = update.Message.Chat.ID
	}
	userStates := BotContext.GetUserStattes()
	state := userStates[chatID]
	editMessage := tgbotapi.NewEditMessageCaption(chatID, state.MessageID, Caption)
	editMessage.ReplyMarkup = &keyboard
	_, err := bot.Send(editMessage)
	if err != nil {
		logger.Debug(fmt.Sprintf("Error edit photo message: %v", err))
	}
}

func ReturnStartMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, Caption string) {
	var userID int64
	if update.Message == nil {
		userID = update.CallbackQuery.From.ID
	} else {
		userID = update.Message.Chat.ID
	}
	userStates := BotContext.GetUserStattes()
	state := userStates[userID]
	state.AllQuestions = nil
	state.UserAnswers = nil
	state.CorrectAnswers = nil
	state.Conversation = nil
	state.CurrentMenu = "main"
	state.Data["subject"] = ""
	state.Data["Topic"] = ""
	state.Data["level"] = ""
	delete(state.Data, "score")
	BotContext.SetUserState(userID, state)
	ShowStartMenu(bot, update, logger, BotContext, Caption)
}

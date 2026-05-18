package menu

import (
	"TeacherBot/models"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func ShowStartMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *models.BotContext, Caption string) {
	// Создаем инлайн клавиатуру
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Простой тест", "simple"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Интерактивный тест", "interactive"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "settings"),
		),
	)
	sendMenu(bot, update, Caption, keyboard, logger, BotContext)
}

func ShowLevelMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *models.BotContext) {
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

	sendMenu(bot, update, "Выберите сложность", keyboard, logger, BotContext)
}
func ShowBeginnerMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *models.BotContext) {
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

	sendMenu(bot, update, "Выберите тему", keyboard, logger, BotContext)
}
func ShowIntermediateMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *models.BotContext) {
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

	sendMenu(bot, update, "Выберите тему", keyboard, logger, BotContext)
}

func ShowAdvancMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *models.BotContext) {
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

	sendMenu(bot, update, "Выберите тему", keyboard, logger, BotContext)
}
func ShowTestMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, Caption string, logger *zap.Logger, BotContext *models.BotContext) {
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
	sendMenu(bot, update, Caption, keyboard, logger, BotContext)
}
func ShowSetingMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *models.BotContext) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Простой тест", "simple"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Интерактивный тест", "interactive"),
		),
	)
	sendMenu(bot, update, "Выберите тип теста", keyboard, logger, BotContext)
}
func sendMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, Caption string, keyboard tgbotapi.InlineKeyboardMarkup, logger *zap.Logger, BotContext *models.BotContext) {
	var chatID int64
	if update.Message == nil {
		chatID = update.CallbackQuery.From.ID
		userStates := BotContext.GetUserStattes()
		state := userStates[chatID]
		logger.Debug(fmt.Sprintf("MessageID: %v", state.MessageID))
		if _, exist := state.Data["nopoto"]; exist {
			editMessage := tgbotapi.NewEditMessageText(chatID, state.MessageID, Caption)
			editMessage.ReplyMarkup = &keyboard
			_, err := bot.Send(editMessage)
			if err != nil {
				logger.Debug(fmt.Sprintf("Error edit photo message: %v", err))
			}
			return
		}
		editMessage := tgbotapi.NewEditMessageCaption(chatID, state.MessageID, Caption)
		editMessage.ReplyMarkup = &keyboard
		_, err := bot.Send(editMessage)
		if err != nil {
			logger.Debug(fmt.Sprintf("Error edit photo message: %v", err))
		}
		return
	}
	chatID = update.Message.Chat.ID
	userStates := BotContext.GetUserStattes()
	state := userStates[chatID]
	if state.MessageID == 0 {
		//Первое сообщение пользователю
		photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("Image/start.jpg"))
		photoMsg.Caption = Caption
		photoMsg.ReplyMarkup = keyboard
		sendMessage, err := bot.Send(photoMsg)
		if err != nil {
			// Если фото не отправилось (файл не найден), отправляем только текст
			logger.Error(fmt.Sprintf("Error send photo: %v", err))
			textMsg := tgbotapi.NewMessage(chatID, Caption)
			textMsg.ReplyMarkup = keyboard
			sendMessage, _ = bot.Send(textMsg)
			state.MessageID = sendMessage.MessageID
			state.Data["nopoto"] = ""
			BotContext.SetUserState(chatID, state)
			return
		}

		state.MessageID = sendMessage.MessageID
		BotContext.SetUserState(chatID, state)
		logger.Info(fmt.Sprintf("Message ID: %v", sendMessage.MessageID))
		return
	}
	//Повторное сообщение пользователю
	if update.Message.Photo != nil {
		//Пользователь отправил сообщение с фото
		editMessage := tgbotapi.NewEditMessageCaption(chatID, state.MessageID, Caption)
		editMessage.ReplyMarkup = &keyboard
		_, err := bot.Send(editMessage)
		if err != nil {
			logger.Debug(fmt.Sprintf("Error edit photo message: %v", err))
		}
		return
	}
	//Пользователь отправил сообщения без фото
	editMessage := tgbotapi.NewEditMessageText(chatID, state.MessageID, Caption)
	editMessage.ReplyMarkup = &keyboard
	_, err := bot.Send(editMessage)
	if err != nil {
		logger.Debug(fmt.Sprintf("Error edit text message: %v", err))
	}
}

package menu

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func getChatID(update tgbotapi.Update) int64 {
	var chatID int64
	if update.Message == nil {
		chatID = update.CallbackQuery.From.ID
		return chatID
	}
	chatID = update.Message.Chat.ID
	return chatID
}

func ShowStartMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger) {
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
	sendMenu(bot, update, "👋 Добро пожаловать в бот!", keyboard, logger)
}

func ShowLevelMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger) {
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

	sendMenu(bot, update, "Выберите сложность", keyboard, logger)
}
func ShowBeginnerMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger) {
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

	sendMenu(bot, update, "Выберите тему", keyboard, logger)
}
func ShowIntermediateMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger) {
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

	sendMenu(bot, update, "Выберите тему", keyboard, logger)
}

func ShowAdvancMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger) {
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

	sendMenu(bot, update, "Выберите тему", keyboard, logger)
}
func sendMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, Caption string, keyboard tgbotapi.InlineKeyboardMarkup, logger *zap.Logger) {
	chatID := getChatID(update)
	photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("Image/start.jpg"))
	photoMsg.Caption = Caption
	photoMsg.ReplyMarkup = keyboard
	if _, err := bot.Send(photoMsg); err != nil {
		// Если фото не отправилось (файл не найден), отправляем только текст
		logger.Error(fmt.Sprintf("Ошибка отправки фото: %v", err))
		textMsg := tgbotapi.NewMessage(chatID, Caption)
		textMsg.ReplyMarkup = keyboard
		bot.Send(textMsg)
	}
}

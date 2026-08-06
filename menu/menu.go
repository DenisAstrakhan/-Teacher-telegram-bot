package menu

import (
	"TeacherBot/domain"
	"TeacherBot/models"
	"fmt"
	"strconv"
	"strings"

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
	chatID := getUserID(update)
	state, err := BotContext.UserRepository.CacheRepository.Get(int(chatID))
	if err != nil {
		logger.Error("Ошибка при получении данных из cache: %w", zap.Error(err))
		ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	}
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
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Учитель", "teacher"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ученик", "student"),
		),
	)
	chatID := getUserID(update)
	photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("Image/WhoAreYou.jpeg"))
	photoMsg.Caption = "Кто ты, воин?"
	photoMsg.ReplyMarkup = keyboard
	_, err := bot.Send(photoMsg)
	if err != nil {
		//Не удалось отправить сообщение
		logger.Error("Error send photo: %w", zap.Error(err))
		return
	}
}

func ShowTeacherMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext) {
	teacherID := getUserID(update)
	//Получаем список учеников
	usersList, err := BotContext.UserRepository.DataRepository.GetStudentsByTeacher(int(teacherID), 100)
	if err != nil {
		logger.Error("Неудолось получить список учеников: %w", zap.Error(err))
		return
	}
	//проверяем не пустой л список
	if len(usersList) == 0 {
		logger.Info(fmt.Sprintf("У учителя ID - %d нет учеников", teacherID))
		return
	}
	state, err := BotContext.UserRepository.CacheRepository.Get(int(teacherID))
	if err != nil {
		logger.Error("Ошибка при получении данных из cache: %w", zap.Error(err))
		ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	}
	state.StudentList = usersList
	state.CurrentMenu = "teacher"
	SetState(bot, update, BotContext, logger, teacherID, state)
	// Формируем текст с нумерованным списком
	var sb strings.Builder
	sb.WriteString("📋 *Выберите ученика:*\n\n")
	// Создаём клавиатуру с номерами
	var rows [][]tgbotapi.InlineKeyboardButton
	for i, user := range usersList {
		var url string
		if user.Telegram_name == "" {
			url = "Нет ссылки!"
		} else {
			url = fmt.Sprintf("https://t.me/%s", user.Telegram_name)
		}
		fmt.Fprintf(&sb, "%d. %s\n", i+1, url)
		// Добавляем кнопку с номером
		button := tgbotapi.NewInlineKeyboardButtonData(user.Full_name, strconv.Itoa(i+1))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	lists := sb.String()
	if state.MessageID == 0 {
		sendMenu(bot, update, lists, tgbotapi.NewInlineKeyboardMarkup(rows...), logger, BotContext, "Image/Techer.jpg")
	} else {
		editMenu(bot, update, lists, tgbotapi.NewInlineKeyboardMarkup(rows...), logger, BotContext)
	}
}
func ShowTecherList(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext) models.UserState {
	userID := getUserID(update)
	state, err := BotContext.UserRepository.CacheRepository.Get(int(userID))
	if err != nil {
		logger.Error("Ошибка при получении данных из cache: %w", zap.Error(err))
		ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return state
	}
	teacherLists, err := BotContext.UserRepository.DataRepository.GetTeacherLists()
	if err != nil {
		logger.Error("Не удалось получит полный список учетелей. Ошибка: %w", zap.Error(err))
		return state
	}
	state.TeacherLists = teacherLists
	if len(teacherLists) == 0 {
		//Список учителей пуст
		logger.Debug("Список учителей пуст!!")
		return state
	}

	// Формируем текст с нумерованным списком
	lists := "📋 *Выберите своего учителя:*\n\n"

	// Создаём клавиатуру с номерами
	var rows [][]tgbotapi.InlineKeyboardButton

	for i, teacher := range teacherLists {
		number := i + 1
		// Добавляем кнопку с номером
		button := tgbotapi.NewInlineKeyboardButtonData(teacher.Teacher_name, strconv.Itoa(number))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	msg := tgbotapi.NewMessage(userID, lists)
	if _, err := bot.Send(msg); err != nil {
		logger.Error("Error sending message: %w", zap.Error(err))
	}
	// Отправляем клавиатуру с номерами
	keyboardMsg := tgbotapi.NewMessage(userID, "👇 *Нажмите номер учителя:*")
	keyboardMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	if _, err := bot.Send(keyboardMsg); err != nil {
		logger.Error("Error sending keyboard: %w", zap.Error(err))
		return state
	}
	return state
}
func ShowTestList(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext, Caption string) {
	userID := getUserID(update)
	state, err := BotContext.UserRepository.CacheRepository.Get(int(userID))
	if err != nil {
		logger.Error("Ошибка при получении данных из cache: %w", zap.Error(err))
		ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	}
	testList, err := BotContext.UserRepository.DataRepository.GetTestByUser(state.Student.Telegram_id)
	if err != nil {
		logger.Error(fmt.Sprintf("Не удалось получить список тестов. Ошибка: %v", err))
	}
	state.TestList = testList
	SetState(bot, update, BotContext, logger, userID, state)
	if len(testList) == 0 {
		logger.Info("У ученика нет пройденных тестов")
		return
	}
	//Создаём клавиатуру
	var rows [][]tgbotapi.InlineKeyboardButton
	button := tgbotapi.NewInlineKeyboardButtonData("Показать отценки", "result")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	for _, test := range testList {

		button := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Предмет: %s. Тема: %s. Отценка: %d", test.Subject, test.Topic, test.Result), strconv.Itoa(test.Id))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	button = tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))

	// Отправляем клавиатуру с номерами
	editMenu(bot, update, Caption, tgbotapi.NewInlineKeyboardMarkup(rows...), logger, BotContext)

}

func ShowTestListMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext, Caption string) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удолить тест", "delete test"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Исправить отценку", "resoult edit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back"),
		),
	)
	//отправляем меню
	editMenu(bot, update, "👇 Выберите действие:", keyboard, logger, BotContext)

}

func ShowResultMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger, BotContext *domain.BotContext) {
	chatID := getUserID(update)
	state, err := BotContext.UserRepository.CacheRepository.Get(int(chatID))
	if err != nil {
		logger.Error("Ошибка при получении данных из cache: %w", zap.Error(err))
		ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	}
	result, err := BotContext.UserRepository.DataRepository.GetResultByUser(state.Student.Telegram_id, 100)
	if err != nil {
		logger.Warn(fmt.Sprintf("Не удолось получить результаты тестов. Ошибка: %v", err))
		return
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Отценки ученика %s \n\n", state.Student.Full_name)
	for i, text := range result {
		fmt.Fprintf(&sb, "%d. Результат: %d. Время окончания: %v \n", i+1, text.Result, text.Time_finish)
	}
	resultList := sb.String()
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back"),
		),
	)
	editMenu(bot, update, resultList, keyboard, logger, BotContext)
}

func sendMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, Caption string, keyboard tgbotapi.InlineKeyboardMarkup, logger *zap.Logger, BotContext *domain.BotContext, imageName string) {
	chatID := getUserID(update)
	state, err := BotContext.UserRepository.CacheRepository.Get(int(chatID))
	if err != nil {
		logger.Error("Ошибка при получении данных из cache: %w", zap.Error(err))
		ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	}
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
	logger.Debug(fmt.Sprintf("Saving MessageID: %d for chat: %d", state.MessageID, chatID))
	SetState(bot, update, BotContext, logger, chatID, state)
}

func editMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, Caption string, keyboard tgbotapi.InlineKeyboardMarkup, logger *zap.Logger, BotContext *domain.BotContext) {
	chatID := getUserID(update)
	state, errget := BotContext.UserRepository.CacheRepository.Get(int(chatID))
	if errget != nil {
		logger.Error("Ошибка при получении данных из cache: %w", zap.Error(errget))
		ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	}
	logger.Debug(fmt.Sprintf("Trying to edit MessageID: %d for chat: %d", state.MessageID, chatID))
	editMessage := tgbotapi.NewEditMessageCaption(chatID, state.MessageID, Caption)
	editMessage.ReplyMarkup = &keyboard
	_, err := bot.Send(editMessage)
	if err != nil {
		logger.Sugar().Debugf("Error edit photo message: %v", err)
	}
}

func ReturnStartMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, Caption string) {
	var userID int64
	if update.Message == nil {
		userID = update.CallbackQuery.From.ID
	} else {
		userID = update.Message.Chat.ID
	}
	state, err := BotContext.UserRepository.CacheRepository.Get(int(userID))
	if err != nil {
		logger.Error("Ошибка при получении данных из cache: %w", zap.Error(err))
		ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	}
	state.AllQuestions = nil
	state.UserAnswers = nil
	state.CorrectAnswers = nil
	state.Conversation = nil
	state.CurrentMenu = "main"
	state.Data["subject"] = ""
	state.Data["Topic"] = ""
	state.Data["level"] = ""
	delete(state.Data, "score")
	SetState(bot, update, BotContext, logger, userID, state)
	ShowStartMenu(bot, update, logger, BotContext, Caption)
}
func getUserID(update tgbotapi.Update) int64 {
	var chatID int64
	if update.Message == nil {
		chatID = update.CallbackQuery.From.ID
		return chatID
	}
	chatID = update.Message.Chat.ID
	return chatID
}

func SetState(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, userID int64, state models.UserState) {
	if err := BotContext.UserRepository.CacheRepository.SetWithTTL(int(userID), state, 0); err != nil {
		logger.Error("Ошибка при получении данных из cache: %w", zap.Error(err))
		ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
	}
}

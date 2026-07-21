package handlers

import (
	"TeacherBot/domain"
	gchat "TeacherBot/gigachat"
	"TeacherBot/menu"
	"TeacherBot/models"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func HandleCallback(logger *zap.Logger, bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext) {
	// Отвечаем на callback (убираем "часики")
	if _, err := bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "")); err != nil {
		logger.Error(fmt.Sprintf("Error answering callback: %v", err))
		return
	}
	userID := update.CallbackQuery.From.ID
	data := update.CallbackQuery.Data
	logger.Info(fmt.Sprintf("User ID - %v: press \"%s\" ", userID, data))
	// Инициализируем состояние пользователя
	userNew, state, err := initializationUserStates(logger, userID, BotContext, bot, update)
	logger.Debug(fmt.Sprintf("userNew: %t", userNew))
	if err != nil {
		logger.Error(fmt.Sprintf("error looking up user in database: %s", err))
		return
	}
	if !userNew {
		//Защита от повторного нажатий
		BotContext.Mtx.Lock()
		if time.Since(state.UserLastPress) < 1000*time.Millisecond {
			logger.Warn(fmt.Sprintf("User ID - %v: press again", userID))
			BotContext.Mtx.Unlock()
			return
		}
		state.UserLastPress = time.Now()
		BotContext.Mtx.Unlock()
	}
	// Обработка callback данных
	switch data {
	case "result":
		state.CurrentMenu = "test list"
		menu.ShowResultMenu(bot, update, logger, BotContext)
	case "delete test":
		if err := BotContext.UserRepository.DeleteRow("tests", "id", state.TestID); err != nil {
			menu.ShowTestListMenu(bot, update, logger, BotContext, "Не удалось удолить тест.")
			logger.Warn(fmt.Sprintf("Не удолось удолить строку. Ошибка: %v", err))
			return
		}
		state.CurrentMenu = "test"
		menu.ShowTestList(bot, update, logger, BotContext, "Тест успешно удалён")
		logger.Info(fmt.Sprintf("Тест ID - %d удолён.", state.TestID))
	case "resoult edit":
		state.Data["resoult edit"] = ""
		msg := tgbotapi.NewMessage(userID, "Введите новую отценку, от 0 до 100.")
		if _, err := bot.Send(msg); err != nil {
			logger.Error(fmt.Sprintf("Error sending mesage: %v", err))
		}
	case "teacher":
		if userNew {
			teacher := true
			state = models.NewUserState(&teacher)
			state.CurrentMenu = "teacher"
			BotContext.SetUserState(userID, state)
			logger.Info(fmt.Sprintf("New teacher added, ID: - %d", userID))
			msg := tgbotapi.NewMessage(userID, "Введите своё имя")
			if _, err := bot.Send(msg); err != nil {
				logger.Error(fmt.Sprintf("Error sending message: %v", err))
				return
			}
			menu.ShowTeacherMenu(bot, update, logger, BotContext)
			return
		}

	case "student":
		if userNew {
			teacher := false
			state = models.NewUserState(&teacher)
			state.Data["user name"] = ""
			BotContext.SetUserState(userID, state)
			logger.Info(fmt.Sprintf("New user added, ID: - %d", userID))
			msg := tgbotapi.NewMessage(userID, "Введите своё имя")
			if _, err := bot.Send(msg); err != nil {
				logger.Error(fmt.Sprintf("Error sending message: %v", err))
				return
			}
			return
		}

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

		//проверяем является ли введённый текст числом
		choice, err := strconv.Atoi(data)
		if err != nil {
			logger.Info(fmt.Sprintf("User ID - %v: Failed to process callback", userID))
			menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!")
			return
		}
		//Проверем кто ввёл число учитель или ученик
		if state.Teacher != nil && *state.Teacher {
			if state.CurrentMenu == "teacher" {
				//Учитель выбирает ученика
				if choice > 0 && len(state.StudentList) >= choice {
					//Ученик есть в списке
					state.Student = state.StudentList[choice-1]
					state.CurrentMenu = "test"
					BotContext.SetUserState(userID, state)
					menu.ShowTestList(bot, update, logger, BotContext, "👇 *Выберите тест:*")
					return
				}
				logger.Warn("Ученика не оказалось в списке")
				return
			}
			if state.CurrentMenu == "test" {
				//Учитель выбирает тест
				state.TestID = choice
				state.CurrentMenu = "test list"
				BotContext.SetUserState(userID, state)
				menu.ShowTestListMenu(bot, update, logger, BotContext, "👇 Выберите действие:")
				return
			}

		}
		if choice > 0 && len(state.TeacherLists) >= choice {
			//Учитель есть в списке учителей
			err := BotContext.UserRepository.InsertUser(int(userID), &update.CallbackQuery.From.UserName, state.Data["user name"], state.TeacherLists[choice-1].Telegram_id)
			if err != nil {
				logger.Error(fmt.Sprintf("Ошибка при добавлении пользователя в базу данных: %v", err))
				return
			}
			logger.Info(fmt.Sprintf("Пользователь ID-%d добавлен в базу данных.", userID))
			BotContext.SetUserState(userID, state)
			menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!")
			return
		}

	}
	BotContext.SetUserState(userID, state)
}
func goBack(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger) {
	userID := update.CallbackQuery.From.ID
	userStates := BotContext.UserStates
	state := userStates[userID]
	switch state.CurrentMenu {
	case "test list":
		state.CurrentMenu = "test"
		menu.ShowTestList(bot, update, logger, BotContext, "👇 Выберите действие:")
	case "test":
		state.CurrentMenu = "teacher"
		menu.ShowTeacherMenu(bot, update, logger, BotContext)
	case "simple", "interactive":
		state.CurrentMenu = "main"
		menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!")
	case "Beginner", "Intermediate", "Advanc":
		state.CurrentMenu = state.Data["test"]
		menu.ShowLevelMenu(bot, update, logger, BotContext)
	case "setting":
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	/*case "teacher":
	menu.ShowTeacherMenu(bot, update, logger, BotContext)
	return*/
	default:
		menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!")
		state.CurrentMenu = "main"
	}

	BotContext.SetUserState(userID, state)
}

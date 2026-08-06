package handlers

import (
	"TeacherBot/domain"
	gchat "TeacherBot/gigachat"
	"TeacherBot/menu"
	"TeacherBot/models"
	"context"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func HandleCallback(logger *zap.Logger, bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, ctx context.Context) {
	// Отвечаем на callback (убираем "часики")
	if _, err := bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "")); err != nil {
		logger.Error("Error answering callback: %w", zap.Error(err))
		return
	}
	userID := update.CallbackQuery.From.ID
	data := update.CallbackQuery.Data
	logger.Sugar().Infof("User %d: press \"%s\" ", userID, data)
	// Инициализируем состояние пользователя
	userNew, state, err := initializationUserStates(logger, userID, BotContext, bot, update, ctx)
	logger.Sugar().Debugf("userNew: %t", userNew)
	if err != nil {
		logger.Error("Error looking up user in database: %w", zap.Error(err))
		return
	}
	if !userNew {
		//Защита от повторного нажатий
		BotContext.Mtx.Lock()
		if time.Since(state.UserLastPress) < 1000*time.Millisecond {
			logger.Sugar().Warnf("User %d: press again", userID)
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
		menu.ShowResultMenu(bot, update, logger, BotContext, ctx)
	case "delete test":
		if err := BotContext.UserRepository.DataRepository.DeleteRow("tests", "id", state.TestID, ctx); err != nil {
			menu.ShowTestListMenu(bot, update, logger, BotContext, "Не удалось удолить тест.", ctx)
			logger.Sugar().Warnf("Failed to delete record: %w", err)
			return
		}
		state.CurrentMenu = "test"
		menu.ShowTestList(bot, update, logger, BotContext, "Тест успешно удалён", ctx)
		logger.Sugar().Infof("Test %d deleted", state.TestID)
	case "resoult edit":
		state.Data["resoult edit"] = ""
		msg := tgbotapi.NewMessage(userID, "Введите новую отценку, от 0 до 100.")
		if _, err := bot.Send(msg); err != nil {
			logger.Error("Error sending mesage: %w", zap.Error(err))
		}
	case "teacher":
		if userNew {
			teacher := true
			state = models.NewUserState(&teacher)
			state.CurrentMenu = "teacher"
			menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
			logger.Sugar().Infof("New teacher added, ID: - %d", userID)
			msg := tgbotapi.NewMessage(userID, "Введите своё имя")
			if _, err := bot.Send(msg); err != nil {
				logger.Error("Error sending message: %w", zap.Error(err))
				return
			}
			menu.ShowTeacherMenu(bot, update, logger, BotContext, ctx)
			return
		}

	case "student":
		if userNew {
			teacher := false
			state = models.NewUserState(&teacher)
			state.Data["user name"] = ""
			menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
			logger.Sugar().Infof("New user added, ID: - %d", userID)
			msg := tgbotapi.NewMessage(userID, "Введите своё имя")
			if _, err := bot.Send(msg); err != nil {
				logger.Error("Error sending message: %w", zap.Error(err))
				return
			}
			return
		}

	case "simple":
		state.CurrentMenu = "simple"
		state.Data["test"] = "simple"
		if state.Data["subject"] != "" && state.Data["Topic"] != "" && state.Data["level"] != "" {
			menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
			gchat.StartTest(bot, update, BotContext, logger, ctx)
			return
		}
		menu.ShowLevelMenu(bot, update, logger, BotContext, ctx)
	case "interactive":
		state.CurrentMenu = "interactive"
		state.Data["test"] = "interactive"
		if state.Data["subject"] != "" && state.Data["Topic"] != "" && state.Data["level"] != "" {
			menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
			gchat.StartTest(bot, update, BotContext, logger, ctx)
			return
		}
		menu.ShowLevelMenu(bot, update, logger, BotContext, ctx)
	case "settings":
		gchat.SelectSubject(bot, update, BotContext, logger, ctx)
		state.CurrentMenu = "setting"
	case "back":
		goBack(bot, update, BotContext, logger, ctx)
		return
	case "sorry":
		menu.ReturnStartMenu(bot, update, BotContext, logger, "Попробуй всё заново", ctx)
		return
	case "Beginner":
		state.CurrentMenu = "Beginner"
		state.Data["level"] = "Beginner"
		menu.ShowBeginnerMenu(bot, update, logger, BotContext, ctx)
	case "Intermediate":
		state.CurrentMenu = "Intermediate"
		state.Data["level"] = "Intermediate"
		menu.ShowIntermediateMenu(bot, update, logger, BotContext, ctx)
	case "Advanc":
		state.CurrentMenu = "Advanc"
		state.Data["level"] = "Advanc"
		menu.ShowAdvancMenu(bot, update, logger, BotContext, ctx)
	case "Topic1":
		topic1(bot, update, BotContext, logger, userID, state, ctx)
		return
	case "Topic2":
		topic2(bot, update, BotContext, logger, userID, state, ctx)
		return
	case "Topic3":
		topic3(bot, update, BotContext, logger, userID, state, ctx)
		return
	case "Topic4":
		topic4(bot, update, BotContext, logger, userID, state, ctx)
		return
	case "Topic5":
		topic5(bot, update, BotContext, logger, userID, state, ctx)
		return
	case "A":
		state.UserAnswers = append(state.UserAnswers, "A")
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.SimpleTest(bot, update, BotContext, logger, ctx)
		return
	case "B":
		state.UserAnswers = append(state.UserAnswers, "B")
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.SimpleTest(bot, update, BotContext, logger, ctx)
		return
	case "C":
		state.UserAnswers = append(state.UserAnswers, "C")
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.SimpleTest(bot, update, BotContext, logger, ctx)
		return
	case "D":
		state.UserAnswers = append(state.UserAnswers, "D")
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.SimpleTest(bot, update, BotContext, logger, ctx)
		return
	default:

		//проверяем является ли введённый текст числом
		choice, err := strconv.Atoi(data)
		if err != nil {
			logger.Sugar().Infof("User %d: Failed to process callback", userID)
			menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!", ctx)
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
					menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
					menu.ShowTestList(bot, update, logger, BotContext, "👇 *Выберите тест:*", ctx)
					return
				}
				logger.Warn("Student not found in the list")
				return
			}
			if state.CurrentMenu == "test" {
				//Учитель выбирает тест
				state.TestID = choice
				state.CurrentMenu = "test list"
				menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
				menu.ShowTestListMenu(bot, update, logger, BotContext, "👇 Выберите действие:", ctx)
				return
			}

		}
		if choice > 0 && len(state.TeacherLists) >= choice {
			//Учитель есть в списке учителей
			err := BotContext.UserRepository.DataRepository.InsertUser(int(userID), &update.CallbackQuery.From.UserName, state.Data["user name"], state.TeacherLists[choice-1].Telegram_id, ctx)
			if err != nil {
				logger.Error("Failed to add user to the database: %w", zap.Error(err))
				return
			}
			logger.Sugar().Infof("User %d added to the database", userID)
			menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
			menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!", ctx)
			return
		}

	}
	menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
}
func goBack(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, ctx context.Context) {
	userID := update.CallbackQuery.From.ID
	state, err := BotContext.UserRepository.CacheRepository.Get(int(userID), ctx)
	if err != nil {
		logger.Error("Failed to retrieve data from cache: %w", zap.Error(err))
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		return
	}
	switch state.CurrentMenu {
	case "test list":
		state.CurrentMenu = "test"
		menu.ShowTestList(bot, update, logger, BotContext, "👇 Выберите действие:", ctx)
	case "test":
		state.CurrentMenu = "teacher"
		menu.ShowTeacherMenu(bot, update, logger, BotContext, ctx)
	case "simple", "interactive":
		state.CurrentMenu = "main"
		menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!", ctx)
	case "Beginner", "Intermediate", "Advanc":
		state.CurrentMenu = state.Data["test"]
		menu.ShowLevelMenu(bot, update, logger, BotContext, ctx)
	case "setting":
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		return
	default:
		menu.ShowStartMenu(bot, update, logger, BotContext, "👋 Добро пожаловать в бот!", ctx)
		state.CurrentMenu = "main"
	}
	menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
}
func topic1(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, userID int64, state models.UserState, ctx context.Context) {
	switch state.Data["level"] {
	case "Beginner":
		state.Data["Topic"] = "Present Simple & Present Continuous (базовое сравнение)"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	case "Intermediate":
		state.Data["Topic"] = "Present Perfect vs. Past Simple"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	case "Advanc":
		state.Data["Topic"] = "Инверсия (Never have I seen... / Not only did he...)"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	default:
		logger.Sugar().Warnf("User %d: Failed to distribute Topic1 across levels ", userID)
	}
}
func topic2(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, userID int64, state models.UserState, ctx context.Context) {
	switch state.Data["level"] {
	case "Beginner":
		state.Data["Topic"] = "There is / There are + предлоги места"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	case "Intermediate":
		state.Data["Topic"] = "Условные предложения (Conditionals: 0, 1, 2 типы)"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	case "Advanc":
		state.Data["Topic"] = "Смешанные условные предложения (Mixed Conditionals)"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	default:
		logger.Sugar().Warnf("User %d: Failed to distribute Topic2 across levels ", userID)
	}
}
func topic3(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, userID int64, state models.UserState, ctx context.Context) {
	switch state.Data["level"] {
	case "Beginner":
		state.Data["Topic"] = "Модальные глаголы (can / must)"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	case "Intermediate":
		state.Data["Topic"] = "Пассивный залог (Present / Past Simple Passive)"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	case "Advanc":
		state.Data["Topic"] = "Сослагательное наклонение (I suggest that he go / It’s crucial that she be)"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	default:
		logger.Sugar().Warnf("User %d: Failed to distribute Topic3 across levels ", userID)
	}
}
func topic4(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, userID int64, state models.UserState, ctx context.Context) {
	switch state.Data["level"] {
	case "Beginner":
		state.Data["Topic"] = "Простое прошедшее время (правильные и топ-неправильных глаголов)"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	case "Intermediate":
		state.Data["Topic"] = "Косвенная речь (Reported Speech)"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	case "Advanc":
		state.Data["Topic"] = "Эллипсис и замена (So do I / Neither can she / He does)"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	default:
		logger.Sugar().Warnf("User %d: Failed to distribute Topic4 across levels ", userID)
	}
}
func topic5(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, userID int64, state models.UserState, ctx context.Context) {
	switch state.Data["level"] {
	case "Beginner":
		state.Data["Topic"] = "Базовая лексика: семья, еда, дом, время на часах"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	case "Intermediate":
		state.Data["Topic"] = "Фразовые глаголы (look after, give up, run out of и др.)"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	case "Advanc":
		state.Data["Topic"] = "Расширенная идиоматика и стилистическая синонимия (разница между ask / inquire / demand)"
		state.Data["subject"] = "Английский язык"
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		gchat.StartTest(bot, update, BotContext, logger, ctx)
	default:
		logger.Sugar().Warnf("User %d: Failed to distribute Topic5 across levels ", userID)
	}
}

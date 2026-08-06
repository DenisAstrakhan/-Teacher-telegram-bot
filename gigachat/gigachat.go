package gchat

import (
	"TeacherBot/domain"
	"TeacherBot/menu"
	"TeacherBot/models"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tigusigalpa/gigachat-go"
	"go.uber.org/zap"
)

func StartBot() *gigachat.Client {
	// Создание ключа авторизации из учетных данных
	authKey := base64.StdEncoding.EncodeToString(
		[]byte(os.Getenv("GIGACHAT_CLIENT_ID") + ":" + os.Getenv("GIGACHAT_CLIENT_SECRET")),
	)
	// Создание менеджера токенов
	tokenManager := gigachat.NewTokenManager(authKey,
		gigachat.WithScope("GIGACHAT_API_PERS"),
		gigachat.WithInsecureSkipVerify(true), // отключает проверку сертификата и делает соединение уязвимым для MITM-атак.только для разработки.
	)

	// Создание клиента
	client := gigachat.NewClient(tokenManager,
		gigachat.WithClientInsecureSkipVerify(true), // отключает проверку сертификата и делает соединение уязвимым для MITM-атак.только для разработки.
	)
	return client
}
func StartTest(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, ctx context.Context) {
	userID := getUserID(update)
	state, err := BotContext.UserRepository.CacheRepository.Get(int(userID), ctx)
	if err != nil {
		logger.Error("Failed to retrieve data from cache: %w ", zap.Error(err))
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		return
	}
	switch state.Data["test"] {
	case "interactive":
		logger.Sugar().Infof("User %d: Start interactive test Subject: \"%s\", Topic:  \"%s\", Level  \"%s\"", userID, state.Data["subject"], state.Data["Topic"], state.Data["level"])
		InteractiveTest(bot, update, BotContext, logger, ctx)
		return
	case "simple":
		logger.Sugar().Infof("User %d: Start simple test Subject: \"%s\", Topic:  \"%s\", Level  \"%s\"", userID, state.Data["subject"], state.Data["Topic"], state.Data["level"])
		SimpleTest(bot, update, BotContext, logger, ctx)
		return
	default:
		logger.Sugar().Warnf("User %d: Failed to start tes", userID)
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
	}

}
func InteractiveTest(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, ctx context.Context) {

	client := BotContext.GigaChat
	userID := getUserID(update)
	state, err := BotContext.UserRepository.CacheRepository.Get(int(userID), ctx)
	if err != nil {
		logger.Error("Failed to retrieve data from cache: %w", zap.Error(err))
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		return
	}

	//Проверяем счёт
	if _, exists := state.Data["score"]; !exists {
		//Пользователь только начал тест
		msgToDelete := tgbotapi.NewDeleteMessage(userID, state.MessageID)
		if _, err := bot.Request(msgToDelete); err != nil {
			logger.Error("Error sending message: %w", zap.Error(err))
		}
		logger.Sugar().Infof("User %d: Is at the beginning of the test", userID)
		state.Data["score"] = "0"
		// Создаём "учителя" с памятью о ходе теста
		promptfile, err := getPrompt("RunInteractiveTes.txt")
		if err != nil {
			logger.Sugar().Errorf("User %d: Failed to read prompt file Error: %w", userID, err)
			menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
			return
		}
		conversation := gigachat.Conversation(
			fmt.Sprintf(promptfile, state.Data["subject"], state.Data["Topic"], state.Data["level"]),
			"Начни тестирование. Задай первый вопрос.",
		)
		// Получаем вопрос от учителя
		response, err := client.Chat(conversation)
		if err != nil {
			logger.Sugar().Errorf("User %d: Failed to get question: Error: %w", userID, err)
			menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		}
		question := gigachat.ExtractContent(response)
		logger.Sugar().Infof("User %d: Teacher: %s", userID, question)
		conversation = append(conversation, gigachat.Message{Role: "assistant", Content: question})
		// Отправляем вопрос пользоватпелю
		msg := tgbotapi.NewMessage(userID, question)
		if _, err := bot.Send(msg); err != nil {
			logger.Error("Error sending message: %w", zap.Error(err))
		}
		state.Conversation = conversation
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		return
	}

	state.Conversation = append(state.Conversation, gigachat.Message{Role: "user", Content: update.Message.Text})
	// Получаем вопрос от учителя
	response, err := client.Chat(state.Conversation)
	if err != nil {
		logger.Sugar().Errorf("User %d: Failed to get question: Error: %w", userID, err)
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
	}
	question := gigachat.ExtractContent(response)
	logger.Sugar().Infof("User %d: Teacher: %s", userID, question)
	scoreII, end := parseScoreDigit(question)
	if end {
		// Тест окончен
		response := getLetterGrade(scoreII)
		logger.Sugar().Infof("Test finish! User %d result: %s", userID, response)
		jsonData, err := json.Marshal(state.Conversation)
		if err != nil {
			logger.Error("Failed to marshal user conversation to JSON: %w", zap.Error(err))
			finishInteractiveTest(bot, update, BotContext, logger, state, userID, ctx)
			menu.ReturnStartMenu(bot, update, BotContext, logger, response, ctx)
			return
		}
		if err := BotContext.UserRepository.DataRepository.InsertTests(state.Data["subject"], state.Data["level"], state.Data["Topic"], string(jsonData), scoreII, time.Now(), int(userID), ctx); err != nil {
			logger.Error("Failed to add test to the database: %w", zap.Error(err))
			finishInteractiveTest(bot, update, BotContext, logger, state, userID, ctx)
			menu.ReturnStartMenu(bot, update, BotContext, logger, response, ctx)
			return
		}
		logger.Sugar().Infof("User %d added a test to the database", userID)
		finishInteractiveTest(bot, update, BotContext, logger, state, userID, ctx)
		menu.ReturnStartMenu(bot, update, BotContext, logger, response, ctx)
		return
	}
	state.Conversation = append(state.Conversation,
		gigachat.Message{Role: "assistant", Content: question})
	msg := tgbotapi.NewMessage(userID, question)
	if _, err := bot.Send(msg); err != nil {
		logger.Error("Error sending message: %w", zap.Error(err))
	}
	menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
}
func SimpleTest(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, ctx context.Context) {
	client := BotContext.GigaChat
	userID := getUserID(update)
	state, err := BotContext.UserRepository.CacheRepository.Get(int(userID), ctx)
	if err != nil {
		logger.Error("Failed to retrieve data from cache: %w", zap.Error(err))
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		return
	}
	//Проверяем наличие теста
	if len(state.AllQuestions) == 10 {
		len := len(state.UserAnswers)
		// Проверяем количество ответов пользователя
		if len == 10 {
			// Тест завершился
			checkscore, err := checkingAnswer(state.UserAnswers, state.CorrectAnswers)
			if err != nil {
				logger.Sugar().Errorf("User ID - %v: Failed test: Error: %w", userID, err)
				menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
				return
			}
			response := getLetterGrade(checkscore)
			jsonData, err := json.Marshal(state.AllQuestions)
			if err != nil {
				logger.Error("Failed to marshal question list to JSON: %w", zap.Error(err))
				return
			}
			if err := BotContext.UserRepository.DataRepository.InsertTests(state.Data["subject"], state.Data["level"], state.Data["Topic"], string(jsonData), checkscore, time.Now(), int(userID), ctx); err != nil {
				logger.Error("Failed to add test to the database: %w", zap.Error(err))
				return
			}
			logger.Sugar().Infof("User %d added a test to the database ", userID)
			logger.Sugar().Infof("Test finish! User %d result: %s", userID, response)
			menu.ReturnStartMenu(bot, update, BotContext, logger, response, ctx)
			return
		}
		question, correctAnswer, err := parseQuestion(state.AllQuestions[len])
		if err != nil {
			logger.Sugar().Errorf("User %d: Failed test: Error: %v", userID, err)
			menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		}
		state.CorrectAnswers = append(state.CorrectAnswers, correctAnswer)
		menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
		menu.ShowTestMenu(bot, update, question, logger, BotContext, ctx)
		return
	}
	//Пользователь только начал тест
	logger.Sugar().Infof("User %d: Is at the beginning of the test", userID)
	//Получаем тест
	promptfile, err := getPrompt("RunOneRequestTest.txt")
	if err != nil {
		logger.Sugar().Errorf("User %d: Failed to read prompt file Error: %w", userID, err)
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		return
	}
	prompt := fmt.Sprintf(promptfile, state.Data["subject"], state.Data["Topic"], state.Data["level"])
	messages := []gigachat.Message{
		{Role: "user", Content: prompt},
	}
	response, err := client.Chat(messages)
	if err != nil {
		logger.Sugar().Errorf("User %d: Failed to get question: Error: %w", userID, err)
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		return
	}
	allquestions := splitByQuestionNumber(response.Choices[0].Message.Content)
	if len(allquestions) != 10 {
		logger.Sugar().Errorf("User %d: Failed to get 10 test questions", userID)
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		return
	}
	state.AllQuestions = allquestions
	question, correctAnswer, err := parseQuestion(state.AllQuestions[0])
	if err != nil {
		logger.Sugar().Errorf("User %d: Failed test: Error: %w", userID, err)
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		return
	}
	state.CorrectAnswers = append(state.CorrectAnswers, correctAnswer)
	state.CurrentMenu = ""
	menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
	menu.ShowTestMenu(bot, update, question, logger, BotContext, ctx)
}
func SelectSubject(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, ctx context.Context) {
	userID := getUserID(update)
	state, err := BotContext.UserRepository.CacheRepository.Get(int(userID), ctx)
	if err != nil {
		logger.Error("Failed to retrieve data from cache: %w ", zap.Error(err))
		menu.ReturnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!", ctx)
		return
	}
	state.Data["subject"] = ""
	if state.MessageID != 0 {
		logger.Sugar().Debugf("Attempting to delete message - UserID: %d, MessageID: %d", userID, state.MessageID)
		msgToDelete := tgbotapi.NewDeleteMessage(userID, state.MessageID)
		if _, err := bot.Request(msgToDelete); err != nil {
			logger.Sugar().Warnf("Error sending message: %v", err)
		}
	}
	menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
	msg := tgbotapi.NewMessage(userID, "Напишите школьный предмет для которого нужно создать тест")
	sentMsg, err := bot.Send(msg)
	if err != nil {
		logger.Error("Failed to send subject request: %w", zap.Error(err))
		return
	}
	state.MessageID = sentMsg.MessageID
	menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
}
func parseScoreDigit(input string) (int, bool) {
	re := regexp.MustCompile(`SCORE\s*(\d+)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return 0, false
	}
	if score, err := strconv.Atoi(matches[1]); err == nil {
		return score, true
	}
	return 0, false
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
func getPrompt(filename string) (string, error) {
	data, err := os.ReadFile("prompts/" + filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
func splitByQuestionNumber(text string) []string {
	// Регулярное выражение: ищем начало строки, затем цифры, затем точку и пробелы
	re := regexp.MustCompile(`(?m)^\d+\.\s+`)

	// Находим все позиции, где начинаются вопросы
	indices := re.FindAllStringIndex(text, -1)

	// Если нет ни одного вопроса, возвращаем весь текст
	if len(indices) == 0 {
		return []string{text}
	}

	// Создаем слайс для вопросов
	questions := make([]string, 0, len(indices))

	// Проходим по каждому найденному индексу
	for i := range len(indices) {
		// Начало текущего вопроса
		start := indices[i][0]

		// Определяем конец текущего вопроса
		var end int
		if i+1 < len(indices) {
			// Если есть следующий вопрос, вопрос заканчивается перед ним
			end = indices[i+1][0]
		} else {
			// Если это последний вопрос, берем до конца текста
			end = len(text)
		}

		// Извлекаем вопрос и удаляем лишние пробелы
		question := text[start:end]
		questions = append(questions, strings.TrimSpace(question))
	}

	return questions
}
func parseQuestion(question string) (string, string, error) {
	re := regexp.MustCompile(`Правильный ответ:\s*([A-D])`)
	matches := re.FindStringSubmatch(question)
	if len(matches) < 1 {
		//В question нет строки "Правильный ответ: "
		return "", "", errors.New("Answer not found")
	}
	return re.ReplaceAllString(question, ""), matches[1], nil
}
func checkingAnswer(user []string, correct []string) (int, error) {
	if len(user) != len(correct) {
		err := errors.New("Failed to compare user's answers with the correct ones")
		return 0, err
	}
	var score int
	for i, v := range user {
		if v == correct[i] {
			score += 10
		}
	}
	return score, nil
}
func getLetterGrade(percentage int) string {
	switch {
	case percentage >= 90:
		return "A (Отлично!)"
	case percentage >= 80:
		return "B (Хорошо)"
	case percentage >= 70:
		return "C (Удовлетворительно)"
	case percentage >= 60:
		return "D (Проходной)"
	default:
		return "F (Нужно повторить материал)"
	}
}
func finishInteractiveTest(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *domain.BotContext, logger *zap.Logger, state models.UserState, userID int64, ctx context.Context) {
	delete(state.Data, "score")
	state.MessageID = 0
	menu.SetState(bot, update, BotContext, logger, userID, state, ctx)
}

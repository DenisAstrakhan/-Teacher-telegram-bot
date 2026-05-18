package gchat

import (
	"TeacherBot/menu"
	"TeacherBot/models"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

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
func StartTest(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext, logger *zap.Logger) {
	userID := getUserID(update)
	userStates := BotContext.GetUserStattes()
	state := userStates[userID]
	switch state.Data["test"] {
	case "interactive":
		logger.Info(fmt.Sprintf("User ID - %v: Start interactive test Subject: \"%s\", Topic:  \"%s\", Level  \"%s\"", userID, state.Data["subject"], state.Data["Topic"], state.Data["level"]))
		InteractiveTest(bot, update, BotContext, logger)
		return
	case "simple":
		logger.Info(fmt.Sprintf("User ID - %v: Start simple test Subject: \"%s\", Topic:  \"%s\", Level  \"%s\"", userID, state.Data["subject"], state.Data["Topic"], state.Data["level"]))
		SimpleTest(bot, update, BotContext, logger)
		return
	default:
		logger.Warn(fmt.Sprintf("User ID - %v: Failed to start tes", userID))
		returnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
	}

}
func InteractiveTest(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext, logger *zap.Logger) {

	client := BotContext.GigaChat
	userID := getUserID(update)
	state := BotContext.UserStates[userID]

	//Проверяем счёт
	if _, exists := state.Data["score"]; !exists {
		//Пользователь только начал тест
		logger.Info(fmt.Sprintf("User ID - %v: Is at the beginning of the test", userID))
		state.Data["score"] = "0"
		// Создаём "учителя" с памятью о ходе теста
		promptfile, err := getPrompt("RunInteractiveTes.txt")
		if err != nil {
			logger.Error(fmt.Sprintf("User ID - %v: Failed to read prompt file Error: %v", userID, err))
			returnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
			return
		}
		conversation := gigachat.Conversation(
			fmt.Sprintf(promptfile, state.Data["subject"], state.Data["Topic"], state.Data["level"]),
			"Начни тестирование. Задай первый вопрос.",
		)
		// Получаем вопрос от учителя
		response, err := client.Chat(conversation)
		if err != nil {
			logger.Error(fmt.Sprintf("User ID - %v: Failed to get question: Error: %v", userID, err))
			returnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		}
		question := gigachat.ExtractContent(response)
		logger.Info(fmt.Sprintf("User ID - %v: Teacher: %s", userID, question))
		conversation = append(conversation, gigachat.Message{Role: "assistant", Content: question})
		// Отправляем вопрос пользоватпелю
		msg := tgbotapi.NewMessage(userID, question)
		bot.Send(msg)
		state.Conversation = conversation
		BotContext.SetUserState(userID, state)
		return
	}

	state.Conversation = append(state.Conversation, gigachat.Message{Role: "user", Content: update.Message.Text})
	// Получаем вопрос от учителя
	response, err := client.Chat(state.Conversation)
	if err != nil {
		logger.Error(fmt.Sprintf("User ID - %v: Failed to get question: Error: %v", userID, err))
		returnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
	}
	question := gigachat.ExtractContent(response)
	logger.Info(fmt.Sprintf("User ID - %v: Teacher: %s", userID, question))
	scoreII, end := parseScoreDigit(question)
	if end {
		// Тест окончен
		response := getLetterGrade(scoreII)
		logger.Info(fmt.Sprintf("Test finish! User ID - %v result: %s", userID, response))
		delete(state.Data, "score")
		state.MessageID = 0
		BotContext.SetUserState(userID, state)
		returnStartMenu(bot, update, BotContext, logger, response)
		return
	}
	state.Conversation = append(state.Conversation,
		gigachat.Message{Role: "assistant", Content: question})
	msg := tgbotapi.NewMessage(userID, question)
	bot.Send(msg)
	BotContext.SetUserState(userID, state)
}
func SimpleTest(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext, logger *zap.Logger) {
	client := BotContext.GigaChat
	userID := getUserID(update)
	userStates := BotContext.GetUserStattes()
	state := userStates[userID]
	//Проверяем наличие теста
	if len(state.AllQuestions) == 10 {
		len := len(state.UserAnswers)
		// Проверяем количество ответов пользователя
		if len == 10 {
			// Тест завершился
			checkscore, err := checkingAnswer(state.UserAnswers, state.CorrectAnswers)
			if err != nil {
				logger.Error(fmt.Sprintf("User ID - %v: Failed test: Error: %v", userID, err))
				returnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
				return
			}
			response := getLetterGrade(checkscore)
			logger.Info(fmt.Sprintf("Test finish! User ID - %v result: %s", userID, response))
			returnStartMenu(bot, update, BotContext, logger, response)
			return
		}
		question, correctAnswer, err := parseQuestion(state.AllQuestions[len])
		if err != nil {
			logger.Error(fmt.Sprintf("User ID - %v: Failed test: Error: %v", userID, err))
			returnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		}
		state.CorrectAnswers = append(state.CorrectAnswers, correctAnswer)
		BotContext.SetUserState(userID, state)
		menu.ShowTestMenu(bot, update, question, logger, BotContext)
		return
	}
	//Пользователь только начал тест
	logger.Info(fmt.Sprintf("User ID - %v: Is at the beginning of the test", userID))
	//Получаем тест
	promptfile, err := getPrompt("RunOneRequestTest.txt")
	if err != nil {
		logger.Error(fmt.Sprintf("User ID - %v: Failed to read prompt file Error: %v", userID, err))
		returnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	}
	prompt := fmt.Sprintf(promptfile, state.Data["subject"], state.Data["Topic"], state.Data["level"])
	messages := []gigachat.Message{
		{Role: "user", Content: prompt},
	}
	response, err := client.Chat(messages)
	if err != nil {
		logger.Error(fmt.Sprintf("User ID - %v: Failed to get question: Error: %v", userID, err))
		returnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	}
	allquestions := splitByQuestionNumber(response.Choices[0].Message.Content)
	if len(allquestions) != 10 {
		logger.Error(fmt.Sprintf("User ID - %v: Failed to get 10 test questions", userID))
		returnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	}
	state.AllQuestions = allquestions
	question, correctAnswer, err := parseQuestion(state.AllQuestions[0])
	if err != nil {
		logger.Error(fmt.Sprintf("User ID - %v: Failed test: Error: %v", userID, err))
		returnStartMenu(bot, update, BotContext, logger, "👋 Добро пожаловать в бот!")
		return
	}
	state.CorrectAnswers = append(state.CorrectAnswers, correctAnswer)
	state.CurrentMenu = ""
	BotContext.SetUserState(userID, state)
	menu.ShowTestMenu(bot, update, question, logger, BotContext)
}
func SelectSubject(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext) {
	userID := getUserID(update)
	state := BotContext.UserStates[userID]
	state.Data["subject"] = ""
	msgToDelete := tgbotapi.NewDeleteMessage(userID, state.MessageID)
	bot.Send(msgToDelete)
	BotContext.SetUserState(userID, state)
	msg := tgbotapi.NewMessage(userID, "Напишите школьный предмет для которого нужно создать тест")
	bot.Send(msg)
}
func parseScoreDigit(input string) (int, bool) {
	re := regexp.MustCompile(`--SCORE--\s*(\d+)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return 0, false
	}
	if score, err := strconv.Atoi(matches[1]); err == nil {
		return score, true
	}
	return 0, false
}

func returnStartMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext, logger *zap.Logger, Caption string) {
	userID := getUserID(update)
	userStates := BotContext.GetUserStattes()
	state := userStates[userID]
	state.AllQuestions = nil
	state.UserAnswers = nil
	state.CorrectAnswers = nil
	state.CurrentMenu = "main"
	state.Data["subject"] = ""
	state.Data["Topic"] = ""
	state.Data["level"] = ""
	BotContext.SetUserState(userID, state)
	menu.ShowStartMenu(bot, update, logger, BotContext, Caption)
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
	for i := 0; i < len(indices); i++ {
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

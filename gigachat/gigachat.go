package gchat

import (
	"TeacherBot/models"
	"encoding/base64"
	"fmt"
	"os"

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
	state := BotContext.UserStates[userID]
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
		returnStartMenu(bot, update, BotContext)
	}

}
func InteractiveTest(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext, logger *zap.Logger) {

}
func SimpleTest(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext, logger *zap.Logger) {

}
func returnStartMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, BotContext *models.BotContext) {

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

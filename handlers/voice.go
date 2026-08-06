package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func SaveVoice(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger) {
	logger.Sugar().Infof("User %d: sent the voice ", update.Message.From.ID)
	voice := update.Message.Voice
	// Получаем информацию о файле
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: voice.FileID})
	if err != nil {
		logger.Sugar().Warnf("User %d: Error getting file: %w", update.Message.From.ID, err)
		return
	}
	// Создаем папку "out", если её не существует
	if err := os.MkdirAll("out", 0755); err != nil {
		logger.Sugar().Warnf("User %d: Error creating directory: %w", update.Message.From.ID, err)
		return
	}
	// Формируем путь для сохранения
	now := time.Now()
	fileName := fmt.Sprintf("out/voice_ID_%v_%s.ogg", update.Message.From.ID, now.Format("2006-01-02_15-04-05"))

	//Получаем URL для скачивания
	fileURL := file.Link(bot.Token)

	//Скачиваем файл
	httpClient := &http.Client{
		Timeout: 30 * time.Second, // Устанавливаем таймаут
	}
	resp, err := httpClient.Get(fileURL)
	if err != nil {
		logger.Sugar().Warnf("User %d: Error downloading file: %vw", update.Message.From.ID, err)
		return
	}
	defer resp.Body.Close()

	// Создаем файл на диске
	outFile, err := os.Create(fileName)
	if err != nil {
		logger.Sugar().Warnf("User %d: Error creating file: %w", update.Message.From.ID, err)
		return
	}
	defer outFile.Close()
	// Копируем содержимое
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		logger.Sugar().Warnf("User %d: Error saving file: %w", update.Message.From.ID, err)
		return
	}

	// Подтверждаем пользователю
	logger.Sugar().Infof("User %d: Save voice: %s", update.Message.From.ID, fileName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Звуковое сообщение сохранено!")
	if _, err := bot.Send(msg); err != nil {
		logger.Error("Error sending message: %w", zap.Error(err))
	}

}

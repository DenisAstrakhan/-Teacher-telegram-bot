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

func SavePhoto(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *zap.Logger) {
	// Получаем фото максимального размера
	logger.Info(fmt.Sprintf("User ID - %v: sent the photo ", update.Message.From.ID))
	photo := update.Message.Photo[len(update.Message.Photo)-1]

	// Получаем информацию о файле
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
	if err != nil {
		logger.Warn(fmt.Sprintf("User ID - %v: Error getting file: %v", update.Message.From.ID, err))
		return
	}

	// Создаем папку "out", если её не существует
	if err := os.MkdirAll("out", 0755); err != nil {
		logger.Warn(fmt.Sprintf("User ID - %v: Error creating directory: %v", update.Message.From.ID, err))
		return
	}

	// Формируем путь для сохранения
	now := time.Now()
	fileName := fmt.Sprintf("out/photo_ID_%v_%s.jpg", update.Message.From.ID, now.Format("2006-01-02_15-04-05"))

	// Формируем URL для скачивания файла
	// Используем токен бота и file.FilePath из ответа GetFile
	downloadURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s",
		os.Getenv("BOT_TOKEN"),
		file.FilePath)

	// Скачиваем файл через HTTP запрос
	resp, err := http.Get(downloadURL)
	if err != nil {
		logger.Warn(fmt.Sprintf("User ID - %v: Error downloading file: %v", update.Message.From.ID, err))
		return
	}
	defer resp.Body.Close()

	// Создаем файл на диске
	outFile, err := os.Create(fileName)
	if err != nil {
		logger.Warn(fmt.Sprintf("User ID - %v: Error creating file: %v", update.Message.From.ID, err))
		return
	}
	defer outFile.Close()

	// Копируем содержимое
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		logger.Warn(fmt.Sprintf("User ID - %v: Error saving file: %v", update.Message.From.ID, err))
		return
	}

	// Подтверждаем пользователю
	logger.Info(fmt.Sprintf("User ID - %v: Save photo: %s", update.Message.From.ID, fileName))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Изображение сохранено!")
	bot.Send(msg)

}

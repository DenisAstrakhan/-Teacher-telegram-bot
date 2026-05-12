package main

import (
	"TeacherBot/handlers"
	"TeacherBot/logger"
	"TeacherBot/models"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	//Создаём логер
	logger, logFileClose, err := logger.NewLogger("DEBUG")
	if err != nil {
		panic(err)
	}
	defer logFileClose()
	//Подтягиваем переменные окружения
	err = godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
		panic(err)
	}
	BotContext := models.NewBotContext()
	// Инициализируем бот
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		logger.Error("Failed to initialize bot")
		log.Panic(err)
	}

	//Логируем соощеие об успешной инициализации бота
	logger.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))
	//Создаём обдейт конфиг
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	//Создаём канал куда будут приходить обновления пользователей
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		// Обрабатываем callback от инлайн кнопок
		if update.CallbackQuery != nil {
			handlers.HandleCallback(logger, bot, update, BotContext)
			continue
		}
		// Проверяем обычные сообщения
		if update.Message != nil {
			handlers.HandleMesage(logger, bot, update, BotContext)
			continue
		}
		// Проверяем получения изображения
		if len(update.Message.Photo) > 0 {
			handlers.SavePhoto(logger)
			continue
		}
	}
}

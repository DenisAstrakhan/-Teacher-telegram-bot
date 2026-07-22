package main

import (
	filter "TeacherBot/dictionaries"
	"TeacherBot/domain"
	gchat "TeacherBot/gigachat"
	"TeacherBot/handlers"
	"TeacherBot/logger"
	"TeacherBot/repository"
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	//Подтягиваем переменные окружения
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	//Создаём логер
	logger, logFileClose, err := logger.NewLogger(os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(err)
	}
	defer logFileClose()
	//Создаём контекст для работы с репозиторием
	RepositoryContext, RepositoryCancel := context.WithCancel(context.Background())
	defer RepositoryCancel()
	//Создаём подключение к базе данных
	UserRepository, err := repository.NewUserRepository(RepositoryContext)
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных: %w", err)
	}
	defer UserRepository.Close()

	// Создаём Giga chat клиента
	GigaChat := gchat.StartBot()
	//Инициализируем фильтер матерных слов
	filter, err := filter.InitFilter(logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка при создании фильтра нецензурных слов: %v", err))
	}
	//Создаём контекст бота
	BotContext := domain.NewBotContext(UserRepository, GigaChat, filter, logger)
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
		// Проверяем получения изображения
		if len(update.Message.Photo) > 0 {
			handlers.SavePhoto(bot, update, logger)
			continue
		}
		//Проверяем звуковые сообщения
		if update.Message.Voice != nil {
			handlers.SaveVoice(bot, update, logger)
			continue
		}
		// Проверяем обычные сообщения
		if update.Message != nil {
			handlers.HandleMessage(logger, bot, update, BotContext)
			continue
		}

	}
}

package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	//Подтягиваем переменные окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Инициализируем бот
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	//Выстовляем уровень логирования бота Debug
	bot.Debug = true
	//Логируем соощеие об успешной инициализации бота
	log.Printf("Authorized on account %s", bot.Self.UserName)
	//Создаём обдейт конфиг
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	//Создаём канал куда будут приходить обновления пользователей
	updates, _ := bot.GetUpdatesChan(u)
	//Проходим циклом по каналу и проверяем нет ли новых сообщений от пользователя
	for update := range updates {
		if update.Message != nil { // Пришло новое сообщение
			//Логируем имя автора и текст сообщения
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			//Создаём новое сообщение
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//означает что наше сообщение будет направлено ответом на сообщеия польователя
			msg.ReplyToMessageID = update.Message.MessageID
			//Отправляем сообщение
			bot.Send(msg)
		}
	}
}

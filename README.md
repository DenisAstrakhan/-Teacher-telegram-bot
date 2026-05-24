# TeacherBot — Telegram bot for English language testing with GigaChat

[![Go Version](https://img.shields.io/badge/Go-1.25.4-00ADD8?style=flat&logo=go)](https://go.dev)
[![Telegram Bot API](https://img.shields.io/badge/telegram--bot--api-v5-blue?logo=telegram)](https://github.com/go-telegram-bot-api/telegram-bot-api)
[![GigaChat Go](https://img.shields.io/badge/gigachat--go-v1.0.2-purple)](https://github.com/tigusigalpa/gigachat-go)

A bot for conducting English language tests using **GigaChat** from Sber. The user selects their knowledge level and topic, after which GigaChat dynamically generates a test of 10 questions. Two testing formats are supported: classic (multiple choice with 4 options) and interactive (free conversation with AI). If GigaChat returns an invalid test, the bot automatically returns to the main menu.

---

## 🚀 Key Features

- **Topic selection** — topics are predefined in the bot's menu.
- **Two testing modes**:
  - *Classic* — question + 4 answer options.
  - *Interactive* — dialogue with GigaChat as a tutor.
- **Manual subject/topic setup** — users can specify custom values (validated against the subject list).
- **Media saving** — all user images and voice messages are saved to the `out/` folder.
- **Configuration without recompilation**:
  - Allowed subjects list (`SubjectList`).
  - Profanity dictionary (`russian-bad-words`).
  - GigaChat prompts (`RunStepByStepTest`, `RunInteractiveTest`).
  - Logging level (via environment variable).
- **Profanity filtering** — automatic checking of user messages.
- **Logging** — structured logging with rotation (supports multiple levels).

---

## 🛠 Technologies and Libraries

| Library | Version | Purpose |
|---------|---------|---------|
| [go-telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api) | v5.5.1 | Telegram Bot API interaction |
| [gigachat-go](https://github.com/tigusigalpa/gigachat-go) | v1.0.2 | GigaChat API client (Sber) |
| [go-sensitive-word](https://github.com/LuYongwang/go-sensitive-word) | v1.1.0 | Profanity filtering |
| [zap](https://go.uber.org/zap) | v1.28.0 | High-performance structured logging |
| [godotenv](https://github.com/joho/godotenv) | v1.5.1 | Load environment variables from `.env` file |

---

## 📁 Project Structure
TeacherBot/
├── dictionaries/
│ ├── SubjectList # Allowed subjects (one per line)
│ └── russian-bad-words # Profanity dictionary
├── gigachat/ # GigaChat client and test logic
├── handlers/ # User action handlers
├── Image/ # Menu images
├── logger/ # Configurable logger (zap)
├── logs/ # Log files directory
├── menu/ # Bot menu with inline keyboard
├── models/ # Models and constructors
├── out/ # Saved images and voice messages
├── prompts/
│ ├── RunStepByStepTest # Prompt for classic test
│ └── RunInteractiveTest # Prompt for interactive mode
├── .env # Environment variables (not committed)
├── .env.example # Example environment variables
├── go.mod
├── go.sum
└── main.go # Entry point

text

---

## ⚙️ Setup and Launch

### 1. Obtaining Credentials

| Service | Action |
|---------|--------|
| **Telegram** | Register and get `BOT_TOKEN` from [@BotFather](https://t.me/botfather) |
| **Sber AI (GigaChat)** | Register in [Sber AI personal account](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project). Create a project and get `Client ID` and `Client Secret` |

### 2. Environment Variables

Create a `.env` file in the project root:

```env
BOT_TOKEN=your_telegram_bot_token
GIGACHAT_CLIENT_ID=your_client_id
GIGACHAT_CLIENT_SECRET=your_client_secret
LOG_LEVEL=debug   # debug, info, warn, error
3. Configuration Without Recompilation
File	Purpose
dictionaries/SubjectList	Allowed subjects (one per line)
dictionaries/russian-bad-words	Profanity dictionary (used by go-sensitive-word)
prompts/RunStepByStepTest	Prompt for classic test (10 questions with options)
prompts/RunInteractiveTest	Prompt for interactive mode (free dialogue)
LOG_LEVEL variable	Logging level
4. Launch
bash
go mod tidy
go run main.go
Or build a binary:

bash
go build -o teacher-bot.exe main.go
./teacher-bot.exe
🧠 How It Works
Main Flow
User sends /start command → bot shows main menu (graphics from Image/ folder).

User selects knowledge level and topic.

Bot sends a request to GigaChat (via gigachat-go library) with the corresponding prompt.

If GigaChat returns an invalid test (not 10 questions, incorrect format, etc.):

Bot sends an error message

Returns to the start page

Event is logged via zap

If test is valid:

Classic mode: sequential questions with 4 options

Interactive mode: user freely communicates with AI tutor

Media Saving
All user-sent images and voice messages are saved to the out/ folder with the following naming format:

Type	Name Format
Images	photo_<userID>_<timestamp>.jpg
Voice	voice_<userID>_<timestamp>.ogg
Profanity Filtering
Uses the go-sensitive-word library

Dictionary loaded from dictionaries/russian-bad-words

Every message sent to GigaChat is checked by the filter

If profanity is detected → bot sends a warning, request to GigaChat is blocked

Logging (zap)
Log file: logs/2006-01-02T15.04.05.000000.log

Structured format — convenient for parsing and analysis

Levels via LOG_LEVEL:

Level	Description
debug	Full debugging info (including raw GigaChat responses)
info	Main events (startup, topic selection, test completion)
warn	Non-critical errors (retries, timeouts)
error	Critical errors (GigaChat unavailable, file issues)
🔧 Planned Improvements (TODO)
User authentication

Save users and test results to a database

User statistics (number of tests completed, progress)

Export results to PDF

Automatic cleanup of old files in out/ folder

Webhook support instead of Long Polling (if the bot gains a million users)

❓ FAQ
Q: How do I add a new subject?
A: Add a line to dictionaries/SubjectList. No bot restart needed — changes are picked up automatically.

Q: How do I add a new profane word?
A: Add the word to dictionaries/russian-bad-words (one per line). The bot uses go-sensitive-word for filtering.

Q: The bot isn't responding or GigaChat is unavailable?
A: Set LOG_LEVEL=debug in .env and check the logs in logs/. If the issue persists, the bot will return the user to the main menu.

Q: Where are user-uploaded files stored?
A: In the out/ folder. It is recommended to periodically clean it or set up automatic deletion of old files.

Q: Can I change the number of questions?
A: Yes — edit the prompt in prompts/RunStepByStepTest, replacing "10 questions" with the desired number. The bot will automatically adapt.

📄 License
MIT

🤝 Feedback
For questions about improvements and bugs — create an Issue in the project repository.
When editing dictionaries/ or prompts/, no bot restart is required.

📖 Русская версия / Russian Version
TeacherBot — Telegram бот для тестирования по английскому языку с GigaChat
https://img.shields.io/badge/Go-1.25.4-00ADD8?style=flat&logo=go
https://img.shields.io/badge/telegram--bot--api-v5-blue?logo=telegram
https://img.shields.io/badge/gigachat--go-v1.0.2-purple

Бот для проведения тестов по английскому языку с использованием GigaChat от Сбера. Пользователь выбирает уровень знаний и тему, после чего GigaChat динамически генерирует тест из 10 вопросов. Поддерживаются два формата тестирования: классический (с выбором варианта) и интерактивный (свободная беседа с ИИ). При некорректном ответе от GigaChat бот автоматически возвращается в главное меню.

🚀 Основные возможности
Выбор темы — темы предопределены в меню бота.

Два режима тестирования:

Классический — вопрос + 4 варианта ответа.

Интерактивный — диалог с GigaChat в роли репетитора.

Настройка предмета и темы вручную — можно задать свой вариант (с валидацией по списку предметов).

Сохранение медиа — все отправленные пользователем изображения и голосовые сообщения сохраняются в папку out/.

Гибкая конфигурация без перекомпиляции:

Список допустимых предметов (SubjectList).

Словарь ненормативной лексики (russian-bad-words).

Промпты для GigaChat (RunStepByStepTest, RunInteractiveTest).

Уровень логирования (через переменную окружения).

Фильтрация ненормативной лексики — автоматическая проверка сообщений пользователя.

Логирование — структурированное логирование с ротацией (поддержка разных уровней).

🛠 Технологии и библиотеки
Библиотека	Версия	Назначение
go-telegram-bot-api/v5	v5.5.1	Взаимодействие с Telegram Bot API
gigachat-go	v1.0.2	Клиент для GigaChat API (Сбер)
go-sensitive-word	v1.1.0	Фильтрация ненормативной лексики
zap	v1.28.0	Высокопроизводительное структурированное логирование
godotenv	v1.5.1	Загрузка переменных окружения из .env файла

📁 Структура проекта

TeacherBot/
├── dictionaries/
│   ├── SubjectList              # Список допустимых предметов (построчно)
│   └── russian-bad-words        # Словарь ненормативной лексики
├── gigachat/                    # Клиент для GigaChat и основная логика тестов
├── handlers/                    # Обработчики действий пользователя
├── Image/                       # Изображения для меню бота
├── logger/                      # Настраиваемый логер (zap)
├── logs/                        # Директория с лог-файлами
├── menu/                        # Меню бота с инлайн-клавиатурой
├── models/                      # Модели и их конструкторы
├── out/                         # Сохранённые изображения и голосовые сообщения
├── prompts/
│   ├── RunStepByStepTest        # Промпт для классического теста
│   └── RunInteractiveTest       # Промпт для интерактивного режима
├── .env                         # Переменные окружения (не коммитится)
├── .env.example                 # Пример переменных окружения
├── go.mod
├── go.sum
└── main.go                      # Точка входа
⚙️ Настройка и запуск
1. Получение авторизационных данных
Сервис	Действие
Telegram	Зарегистрируйтесь и получите BOT_TOKEN у @BotFather
Sber AI (GigaChat)	Зарегистрируйтесь в личном кабинете Sber AI. Создайте проект и получите Client ID и Client Secret
2. Переменные окружения
Создайте файл .env в корне проекта:

env
BOT_TOKEN=your_telegram_bot_token
GIGACHAT_CLIENT_ID=your_client_id
GIGACHAT_CLIENT_SECRET=your_client_secret
LOG_LEVEL=debug   # debug, info, warn, error
3. Настройка без перекомпиляции
Файл	Назначение
dictionaries/SubjectList	Список разрешённых предметов (по одному на строку)
dictionaries/russian-bad-words	Словарь ненормативной лексики (используется go-sensitive-word)
prompts/RunStepByStepTest	Промпт для классического теста (10 вопросов с вариантами)
prompts/RunInteractiveTest	Промпт для интерактивного режима (свободный диалог)
Переменная LOG_LEVEL	Уровень логирования
4. Запуск
bash
go mod tidy
go run main.go
Или сборка бинарного файла:

bash
go build -o teacher-bot.exe main.go
./teacher-bot.exe
🧠 Как это работает
Основной поток
Пользователь отправляет команду /start → бот показывает главное меню (графика из папки Image/).

Пользователь выбирает уровень знаний и тему.

Бот отправляет запрос к GigaChat (через библиотеку gigachat-go) с соответствующим промптом.

Если GigaChat вернул некорректный тест (не 10 вопросов, неправильный формат и т.д.):

Бот отправляет сообщение об ошибке

Возвращается на стартовую страницу

Событие логируется через zap

Если тест корректен:

Классический режим: последовательная выдача вопросов с 4 вариантами

Интерактивный режим: пользователь свободно общается с ИИ-учителем

Сохранение медиа
Все отправленные пользователем изображения и голосовые сообщения сохраняются в папку out/ с форматом имени:

Тип	Формат имени
Изображения	photo_<userID>_<timestamp>.jpg
Голосовые	voice_<userID>_<timestamp>.ogg
Фильтрация ненормативной лексики
Используется библиотека go-sensitive-word

Словарь загружается из dictionaries/russian-bad-words

Перед отправкой любого сообщения в GigaChat проверяется фильтром

При обнаружении мата → бот отправляет предупреждение, запрос в GigaChat не уходит

Логирование (zap)
Файл логов: logs/2006-01-02T15.04.05.000000.log

Структурированный формат — удобно для парсинга и анализа

Уровни через LOG_LEVEL:

Уровень	Описание
debug	Полная отладочная информация (включая сырые ответы GigaChat)
info	Основные события (запуск, выбор темы, завершение теста)
warn	Некритичные ошибки (повторные запросы, таймауты)
error	Критические ошибки (недоступность GigaChat, проблемы с файлами)
🔧 Возможные улучшения (TODO)
Авторизация пользователей

Сохранение пользователей и результатов тестирования в базе данных

Статистика пользователя (количество пройденных тестов, прогресс)

Экспорт результатов в PDF

Автоматическая очистка папки out/ от старых файлов

Поддержка Webhook вместо Long Polling (Если бот обзаведётся миллионом пользователей)

❓ Частые вопросы
Q: Как добавить новое нецензурное слово?
A: Добавьте слово в dictionaries/russian-bad-words (по одному на строку). Бот использует go-sensitive-word для фильтрации.

Q: Бот не отвечает или GigaChat недоступен?
A: Проверьте LOG_LEVEL=debug в .env и посмотрите файл logs/bot.log. Если проблема повторяется — бот вернёт пользователя в главное меню.

Q: Где лежат загруженные пользователями файлы?
A: В папке out/. Рекомендуется периодически чистить её или настроить автоматическое удаление старых файлов.

Q: Можно ли изменить количество вопросов?
A: Да — отредактируйте промпт в prompts/RunStepByStepTest, заменив "10 вопросов" на нужное число. Бот автоматически адаптируется.


📄 Лицензия
MIT

🤝 Обратная связь
По вопросам доработки и багам — создавайте Issue в репозитории проекта.
При редактировании dictionaries/ или prompts/ перезапуск бота не требуется.

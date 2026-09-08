# TeacherBot — Telegram bot for English language testing with GigaChat

[![Go Version](https://img.shields.io/badge/Go-1.25.4-00ADD8?style=flat&logo=go)](https://go.dev)
[![Telegram Bot API](https://img.shields.io/badge/telegram--bot--api-v5-blue?logo=telegram)](https://github.com/go-telegram-bot-api/telegram-bot-api)
[![GigaChat Go](https://img.shields.io/badge/gigachat--go-v1.0.2-purple)](https://github.com/tigusigalpa/gigachat-go)

A bot for conducting English language tests using GigaChat from Sber. Users select their knowledge level and topic, after which GigaChat dynamically generates a test of 10 questions. Two testing formats are supported: classic (multiple choice with 4 options) and interactive (free conversation with AI). If GigaChat returns an invalid test, the bot automatically returns to the main menu.

## 🚀 Key Features

- User roles — support for two roles: student and teacher.
- Taking tests — students take tests in one of two formats (classic or interactive).
- Saving results — all test results are stored in the database.
- Teacher result management — teachers can view, edit, and delete their students' test results.
- Topic selection — topics are predefined in the bot's menu.
- Manual subject/topic setup — users can specify custom values (validated against the subject list).
- Media saving — all user images and voice messages are saved to the `out/` folder.
- Configuration without recompilation:
  - Allowed subjects list (`SubjectList`).
  - Profanity dictionary (`russian-bad-words`).
  - GigaChat prompts (`RunStepByStepTest`, `RunInteractiveTest`).
  - Logging level (via environment variable).
- Profanity filtering — automatic checking of user messages.
- **Logging** — structured logging with rotation and ClickHouse database support (console, file, and ClickHouse).

## 🗄️ Data Storage

- **Database** — configurable via environment variables. Three database systems are supported:
  - SQLite
  - MySQL
  - PostgreSQL
- **Caching** — implemented using Redis.
- **Logging Storage** — logs can be stored in:
  - Console output
  - Log files (with rotation)
  - ClickHouse database (for centralized log storage and analytics)

---

## 🛠 Technologies and Libraries

| Library | Version | Purpose |
|---------|---------|---------|
| go-telegram-bot-api/v5 | v5.5.1 | Telegram Bot API interaction |
| gigachat-go | v1.0.2 | GigaChat API client (Sber) |
| go-sensitive-word | v1.1.0 | Profanity filtering |
| zap | v1.28.0 | High-performance structured logging |
| godotenv | v1.5.1 | Load environment variables from .env file |
| go-redis/redis/v8 | v8.11.5 | Redis client for caching |
| go-sql-driver/mysql | v1.10.0 | MySQL driver |
| pgx/v5 | v5.10.0 | PostgreSQL driver |
| go-sqlite3 | v1.14.48 | SQLite driver |
| clickhouse-go/v2 | v2.48.0 | ClickHouse driver for log storage |

---

## 📁 Project Structure

TeacherBot/
├── dictionaries/
│ ├── SubjectList # Allowed subjects (one per line)
│ └── russian-bad-words # Profanity dictionary
├── domain/ # Contains data models, repository interface, and application context
├── gigachat/ # GigaChat client and test logic
├── handlers/ # User action handlers
├── Image/ # Menu images
├── logger/ # Configurable logger (zap) with multiple outputs
├── logs/ # Log files directory
├── menu/ # Bot menu with inline keyboard
├── models/ # Models and constructors
├── out/ # Saved images and voice messages
├── prompts/
│ ├── RunStepByStepTest # Prompt for classic test
│ └── RunInteractiveTest # Prompt for interactive mode
├── repository # Implementation of specific databases and cache
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
# Telegram Bot
BOT_TOKEN=your_telegram_bot_token

# GigaChat
GIGACHAT_CLIENT_ID=your_client_id
GIGACHAT_CLIENT_SECRET=your_client_secret

# Logging
LOG_LEVEL=debug   # debug, info, warn, error

# Repository Type (sqlite, mysql, postgres)
REPOSITORY_TYPE=

# PostgreSQL Configuration
POSTGRES_USER=your_user_name
POSTGRES_PASSWORD=your_password
POSTGRES_DB=your_database
POSTGRES_HOST=localhost:5432
POSTGRES_TIMEOUT=5s

# MySQL Configuration
MYSQL_ROOT_PASSWORD=your_root_password
MYSQL_DATABASE=your_database
MYSQL_USER=your_user_name
MYSQL_PASSWORD=your_password
MYSQL_HOST=localhost:3306
MYSQL_TIMEOUT=5s

# SQLite Configuration
SQLITE_DATA_DIR=./out/sqlitedata
SQLITE_TIMEOUT=5s

# Cache Repository Type (redis)
CACHE_REPOSITORY_TYPE=

# Redis Configuration
REDIS_HOST=localhost:6379
REDIS_TIMEOUT=5s
MAX_MEMORY=redis_memory_limit

# Logger Repository Type (clickhouse)
LOGGER_REPOSITORY_TYPE=clickhouse

# ClickHouse Configuration (for log storage)
CLICKHOUSE_HOST=localhost:9000
CLICKHOUSE_DB=logs
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=
CLICKHOUSE_TIMEOUT=5s

# Service Configuration
SERVICE_NAME=my-app

3. Configuration Without Recompilation
File	Purpose
dictionaries/SubjectList	Allowed subjects (one per line)
dictionaries/russian-bad-words	Profanity dictionary (used by go-sensitive-word)
prompts/RunStepByStepTest	Prompt for classic test (10 questions with options)
prompts/RunInteractiveTest	Prompt for interactive mode (free dialogue)
4. Launch
bash

go mod tidy
go run main.go

Or build a binary:
bash

go build -o teacher-bot.exe main.go
./teacher-bot.exe

🧠 How It Works
1. Authorization and Role Detection

    User sends the /start command.

    The program checks if the user exists in the cache (Redis):

        If found in cache → the program determines their status (teacher or student) and displays the corresponding menu.

    If the user is not in cache → the program checks if the user exists in the database:

        If found in the database → the program determines their status and displays the corresponding menu.

        If not found in cache or database → the bot asks: "Who are you?"

2. New User Registration
User Choice	Action
"Teacher"	The program saves the user to the database as a teacher.
"Student"	The program asks the student to select a teacher from the database. After selection, the student is saved with a reference to the chosen teacher.
3. Student Functionality

    The bot displays the main menu (graphics from the Image/ folder).

    The user selects their knowledge level and topic.

    The bot sends a request to GigaChat (via the gigachat-go library) with the corresponding prompt.

    If GigaChat returns an invalid test (not 10 questions, incorrect format, etc.):

        The bot sends an error message.

        Returns to the start page.

        The event is logged via zap.

    If the test is valid:

        Classic mode: sequential questions with 4 answer options.

        Interactive mode: the user freely communicates with the AI tutor.

4. Teacher Functionality
Step	Action
1	The bot displays a list of students linked to this teacher.
2	When a student is selected → the bot shows a list of completed tests for that student.
3	When a specific test is selected → the bot displays a management menu with options:
	• Delete test
	• Edit test result (score/grade)
Media Saving

All user-sent images and voice messages are saved to the out/ folder with the following naming format:
Type	Name Format
Images	photo_<userID>_<timestamp>.jpg
Voice	voice_<userID>>_<timestamp>.ogg
Profanity Filtering

    Uses the go-sensitive-word library

    Dictionary loaded from dictionaries/russian-bad-words

    Every message sent to GigaChat is checked by the filter

    If profanity is detected → bot sends a warning, request to GigaChat is blocked

📊 Logging (zap)

Logging is implemented using the zap library and supports three simultaneous outputs:
Log Outputs
Output	Description
Console	Real-time log output to stdout/stderr
File	Rotating log files in logs/ directory (format: 2006-01-02T15.04.05.000000.log)
ClickHouse	Structured logs stored in ClickHouse database for centralized analytics
Log Levels

Log level is set via the LOG_LEVEL environment variable:
Level	Description
debug	Full debugging info (including raw GigaChat responses)
info	Main events (startup, topic selection, test completion)
warn	Non-critical errors (retries, timeouts)
error	Critical errors (GigaChat unavailable, file issues)
ClickHouse Storage
Table Structure
sql

CREATE TABLE IF NOT EXISTS app_logs (
    timestamp DateTime DEFAULT now(),
    level String,
    service String,
    user_id String,
    message String
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, level);

Table Fields
Field	Type	Description
timestamp	DateTime	Event timestamp (auto-generated)
level	String	Log level (debug, info, warn, error)
service	String	Service name (from SERVICE_NAME)
user_id	String	Telegram user ID
message	String	Log message text
Storage Features

    Partitioning by year and month (toYYYYMM(timestamp)) — for efficient storage and fast access

    Sorting by timestamp and level — for optimal query performance

    Auto-generation of current timestamp when creating a record

Example Queries
sql

-- Last 100 error logs
SELECT *
FROM app_logs
WHERE level = 'error'
ORDER BY timestamp DESC
LIMIT 100;

-- Log count by day for the last week
SELECT 
    toDate(timestamp) as date,
    level,
    count() as count
FROM app_logs
WHERE timestamp > now() - INTERVAL 7 DAY
GROUP BY date, level
ORDER BY date DESC;

-- Search logs by specific user
SELECT *
FROM app_logs
WHERE user_id = '123456789'
  AND timestamp > now() - INTERVAL 1 DAY
ORDER BY timestamp DESC;

ClickHouse Configuration

To enable ClickHouse logging, add to .env:
env

LOGGER_REPOSITORY_TYPE=clickhouse
CLICKHOUSE_HOST=localhost:9000
CLICKHOUSE_DB=logs
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=
CLICKHOUSE_TIMEOUT=5s
SERVICE_NAME=my-app

🔧 Planned Improvements (TODO)

    User authentication

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

Q: How do I enable ClickHouse logging?
A: Set LOGGER_REPOSITORY_TYPE=clickhouse in .env and configure ClickHouse connection settings. All logs will be automatically written to both ClickHouse and file/console.

Q: Can I use multiple log outputs simultaneously?
A: Yes — logs are written to console, file, and ClickHouse (if configured) simultaneously by default.
📄 License

MIT
🤝 Feedback

For questions about improvements and bugs — create an Issue in the repository.
When editing dictionaries/ or prompts/, no bot restart is required.
📖 Русская версия / Russian Version
TeacherBot — Telegram бот для тестирования по английскому языку с GigaChat

https://img.shields.io/badge/Go-1.25.4-00ADD8?style=flat&logo=go
https://img.shields.io/badge/telegram--bot--api-v5-blue?logo=telegram
https://img.shields.io/badge/gigachat--go-v1.0.2-purple

Бот для проведения тестов по английскому языку с использованием GigaChat от Сбера. Пользователь выбирает уровень знаний и тему, после чего GigaChat динамически генерирует тест из 10 вопросов. Поддерживаются два формата тестирования: классический (множественный выбор с 4 вариантами) и интерактивный (свободная беседа с ИИ). Если GigaChat возвращает некорректный тест, бот автоматически возвращается в главное меню.
🚀 Основные возможности

    Роли пользователей — поддержка двух ролей: ученик и учитель.

    Прохождение тестов — ученики проходят тесты в одном из двух форматов (классический или интерактивный).

    Сохранение результатов — все результаты тестов сохраняются в базу данных.

    Управление результатами для учителя — учитель может просматривать, редактировать и удалять результаты тестов своих учеников.

    Выбор темы — темы предопределены в меню бота.

    Ручная настройка предмета/темы — пользователи могут указывать произвольные значения (проверяются по списку предметов).

    Сохранение медиафайлов — все изображения и голосовые сообщения пользователей сохраняются в папку out/.

    Настройка без перекомпиляции:

        Список разрешённых предметов (SubjectList).

        Словарь нецензурной лексики (russian-bad-words).

        Промпты GigaChat (RunStepByStepTest, RunInteractiveTest).

        Уровень логирования (через переменную окружения).

    Фильтрация нецензурной лексики — автоматическая проверка сообщений пользователей.

    Логирование — структурированное логирование с ротацией и поддержкой ClickHouse (консоль, файл, ClickHouse).

🗄️ Хранение данных

    База данных — настраивается через переменные окружения. Поддерживаются три СУБД:

        SQLite

        MySQL

        PostgreSQL

    Кэширование — осуществляется с использованием Redis.

    Хранение логов — логи могут сохраняться в:

        Консоль

        Файлы (с ротацией)

        ClickHouse (для централизованного хранения и аналитики)

🛠 Технологии и библиотеки
Библиотека	Версия	Назначение
go-telegram-bot-api/v5	v5.5.1	Взаимодействие с Telegram Bot API
gigachat-go	v1.0.2	Клиент для GigaChat API (Сбер)
go-sensitive-word	v1.1.0	Фильтрация нецензурной лексики
zap	v1.28.0	Высокопроизводительное структурированное логирование
godotenv	v1.5.1	Загрузка переменных окружения из файла .env
go-redis/redis/v8	v8.11.5	Клиент Redis для кэширования
go-sql-driver/mysql	v1.10.0	Драйвер для MySQL
pgx/v5	v5.10.0	Драйвер для PostgreSQL
go-sqlite3	v1.14.48	Драйвер для SQLite
clickhouse-go/v2	v2.48.0	Драйвер ClickHouse для хранения логов
📁 Структура проекта
text

TeacherBot/
├── dictionaries/
│   ├── SubjectList          # Список разрешённых предметов (по одному на строку)
│   └── russian-bad-words    # Словарь нецензурной лексики
├── domain/                  # Содержит модели данных, интерфейс репозитория и контекст приложения
├── gigachat/                # Клиент GigaChat и логика тестов
├── handlers/                # Обработчики действий пользователя
├── Image/                   # Изображения для меню
├── logger/                  # Настраиваемый логер (zap) с множественными выводами
├── logs/                    # Директория с лог-файлами
├── menu/                    # Меню бота с инлайн-клавиатурой
├── models/                  # Модели и их конструкторы
├── out/                     # Сохранённые изображения и голосовые сообщения
├── prompts/
│   ├── RunStepByStepTest    # Промпт для классического теста
│   └── RunInteractiveTest   # Промпт для интерактивного режима
├── repository               # Реализация конкретных БД и кэша
├── .env                     # Переменные окружения (не коммитится)
├── .env.example             # Пример переменных окружения
├── go.mod
├── go.sum
└── main.go                  # Точка входа

⚙️ Настройка и запуск
1. Получение авторизационных данных
Сервис	Действие
Telegram	Зарегистрируйтесь и получите BOT_TOKEN у @BotFather
Sber AI (GigaChat)	Зарегистрируйтесь в личном кабинете Sber AI. Создайте проект и получите Client ID и Client Secret
2. Переменные окружения

Создайте файл .env в корне проекта:
env

# Telegram Bot
BOT_TOKEN=your_telegram_bot_token

# GigaChat
GIGACHAT_CLIENT_ID=your_client_id
GIGACHAT_CLIENT_SECRET=your_client_secret

# Логирование
LOG_LEVEL=debug   # debug, info, warn, error

# Тип репозитория (sqlite, mysql, postgres)
REPOSITORY_TYPE=

# PostgreSQL
POSTGRES_USER=your_user_name
POSTGRES_PASSWORD=your_password
POSTGRES_DB=your_database
POSTGRES_HOST=localhost:5432
POSTGRES_TIMEOUT=5s

# MySQL
MYSQL_ROOT_PASSWORD=your_root_password
MYSQL_DATABASE=your_database
MYSQL_USER=your_user_name
MYSQL_PASSWORD=your_password
MYSQL_HOST=localhost:3306
MYSQL_TIMEOUT=5s

# SQLite
SQLITE_DATA_DIR=./out/sqlitedata
SQLITE_TIMEOUT=5s

# Кэш (redis)
CACHE_REPOSITORY_TYPE=

# Redis
REDIS_HOST=localhost:6379
REDIS_TIMEOUT=5s
MAX_MEMORY=redis_memory_limit

# Тип репозитория для логов (clickhouse)
LOGGER_REPOSITORY_TYPE=clickhouse

# ClickHouse (для хранения логов)
CLICKHOUSE_HOST=localhost:9000
CLICKHOUSE_DB=logs
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=
CLICKHOUSE_TIMEOUT=5s

# Конфигурация сервиса
SERVICE_NAME=my-app

3. Настройка без перекомпиляции
Файл	Назначение
dictionaries/SubjectList	Разрешённые предметы (по одному на строку)
dictionaries/russian-bad-words	Словарь нецензурной лексики (используется go-sensitive-word)
prompts/RunStepByStepTest	Промпт для классического теста (10 вопросов с вариантами)
prompts/RunInteractiveTest	Промпт для интерактивного режима (свободный диалог)
4. Запуск
bash

go mod tidy
go run main.go

Или сборка бинарного файла:
bash

go build -o teacher-bot.exe main.go
./teacher-bot.exe

🧠 Как это работает
1. Авторизация и определение роли

    Пользователь отправляет команду /start.

    Программа проверяет наличие пользователя в кэше (Redis):

        Если найден в кэше → определяется статус (учитель или ученик) и выдаётся соответствующее меню.

    Если пользователя нет в кэше → программа проверяет наличие пользователя в базе данных:

        Если найден в БД → определяется статус и выдаётся соответствующее меню.

        Если не найден ни в кэше, ни в БД → бот задаёт вопрос: "Кто ты?"

2. Регистрация нового пользователя
Выбор пользователя	Действие
"Учитель"	Программа вносит пользователя в БД как учителя
"Ученик"	Программа предлагает выбрать учителя из списка (из БД). После выбора — сохраняет ученика с привязкой к учителю
3. Функционал ученика

    Бот показывает главное меню (графика из папки Image/).

    Пользователь выбирает уровень знаний и тему.

    Бот отправляет запрос в GigaChat (через библиотеку gigachat-go) с соответствующим промптом.

    Если GigaChat возвращает некорректный тест (не 10 вопросов, неверный формат и т.д.):

        Бот отправляет сообщение об ошибке

        Возвращается на стартовую страницу

        Событие логируется через zap

    Если тест валидный:

        Классический режим: последовательные вопросы с 4 вариантами ответов

        Интерактивный режим: пользователь свободно общается с ИИ-репетитором

4. Функционал учителя
Шаг	Действие
1	Бот отправляет список учеников, привязанных к этому учителю
2	При выборе ученика → показывает список пройденных тестов этого ученика
3	При выборе конкретного теста → выдаёт меню управления с возможностями:
	• Удалить тест
	• Изменить результат (оценку/баллы)
Сохранение медиафайлов

Все изображения и голосовые сообщения, отправленные пользователем, сохраняются в папку out/ со следующим форматом имени:
Тип	Формат имени
Изображения	photo_<userID>_<timestamp>.jpg
Голосовые	voice_<userID>>_<timestamp>.ogg
Фильтрация нецензурной лексики

    Используется библиотека go-sensitive-word

    Словарь загружается из dictionaries/russian-bad-words

    Каждое сообщение, отправляемое в GigaChat, проверяется фильтром

    Если обнаружен мат → бот отправляет предупреждение, запрос к GigaChat блокируется

📊 Логирование (zap)

Логирование реализовано с использованием библиотеки zap и поддерживает три вывода одновременно:
Выводы логов
Вывод	Описание
Консоль	Вывод логов в реальном времени в stdout/stderr
Файл	Ротируемые файлы логов в директории logs/ (формат: 2006-01-02T15.04.05.000000.log)
ClickHouse	Структурированные логи, сохраняемые в ClickHouse для централизованной аналитики
Уровни логирования

Уровень задается через переменную окружения LOG_LEVEL:
Уровень	Описание
debug	Полная отладочная информация (включая сырые ответы GigaChat)
info	Основные события (запуск, выбор темы, завершение теста)
warn	Некритичные ошибки (повторы, таймауты)
error	Критические ошибки (GigaChat недоступен, проблемы с файлами)
Хранение в ClickHouse
Структура таблицы
sql

CREATE TABLE IF NOT EXISTS app_logs (
    timestamp DateTime DEFAULT now(),
    level String,
    service String,
    user_id String,
    message String
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, level);

Поля таблицы
Поле	Тип	Описание
timestamp	DateTime	Время события (автоматически)
level	String	Уровень логирования (debug, info, warn, error)
service	String	Имя сервиса (из SERVICE_NAME)
user_id	String	ID пользователя Telegram
message	String	Текст сообщения лога
Особенности хранения

    Партиционирование по году и месяцу (toYYYYMM(timestamp)) — для эффективного хранения и быстрого доступа

    Сортировка по времени и уровню — для оптимальных запросов

    Автоматическая вставка текущего времени при создании записи

Примеры запросов
sql

-- Последние 100 логов с уровнем error
SELECT *
FROM app_logs
WHERE level = 'error'
ORDER BY timestamp DESC
LIMIT 100;

-- Количество логов по дням за последнюю неделю
SELECT 
    toDate(timestamp) as date,
    level,
    count() as count
FROM app_logs
WHERE timestamp > now() - INTERVAL 7 DAY
GROUP BY date, level
ORDER BY date DESC;

-- Поиск логов по конкретному пользователю
SELECT *
FROM app_logs
WHERE user_id = '123456789'
  AND timestamp > now() - INTERVAL 1 DAY
ORDER BY timestamp DESC;

Настройка ClickHouse

Для включения логирования в ClickHouse добавьте в .env:
env

LOGGER_REPOSITORY_TYPE=clickhouse
CLICKHOUSE_HOST=localhost:9000
CLICKHOUSE_DB=logs
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=
CLICKHOUSE_TIMEOUT=5s
SERVICE_NAME=my-app

🔧 Планируемые улучшения (TODO)

    Аутентификация пользователей

    Статистика пользователей (количество пройденных тестов, прогресс)

    Экспорт результатов в PDF

    Автоматическая очистка старых файлов в папке out/

    Поддержка Webhook вместо Long Polling (если бот наберёт миллион пользователей)

❓ Частые вопросы

В: Как добавить новый предмет?
О: Добавьте строку в dictionaries/SubjectList. Перезапуск бота не требуется — изменения подхватываются автоматически.

В: Как добавить новое нецензурное слово?
О: Добавьте слово в dictionaries/russian-bad-words (по одному на строку). Бот использует go-sensitive-word для фильтрации.

В: Бот не отвечает или GigaChat недоступен?
О: Установите LOG_LEVEL=debug в .env и проверьте логи в logs/. Если проблема сохраняется, бот вернёт пользователя в главное меню.

В: Где хранятся файлы, загруженные пользователями?
О: В папке out/. Рекомендуется периодически очищать её или настроить автоматическое удаление старых файлов.

В: Можно ли изменить количество вопросов?
О: Да — отредактируйте промпт в prompts/RunStepByStepTest, заменив "10 вопросов" на нужное количество. Бот автоматически адаптируется.

В: Как включить логирование в ClickHouse?
О: Установите LOGGER_REPOSITORY_TYPE=clickhouse в .env и настройте параметры подключения к ClickHouse. Все логи будут автоматически записываться одновременно в ClickHouse, файл и консоль.

В: Можно ли использовать несколько выводов логов одновременно?
О: Да — логи по умолчанию одновременно записываются в консоль, файл и ClickHouse (если настроен).
📄 Лицензия

MIT
🤝 Обратная связь

По вопросам улучшений и ошибок — создавайте Issue в репозитории проекта.
При редактировании словарей и промптов перезапуск бота не требуется.
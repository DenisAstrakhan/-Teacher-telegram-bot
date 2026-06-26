CREATE SCHEMA bot;

CREATE TABLE bot.teacher(
    telegram_id BIGINT PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    telegram_name VARCHAR (250),
    teacher_name VARCHAR (250) NOT NULL
);

CREATE TABLE bot.users(
    telegram_id BIGINT PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    telegram_name VARCHAR (250),
    full_name VARCHAR (250) NOT NULL,
    teacher_ID BIGINT NOT NULL REFERENCES bot.teacher(telegram_id)
);

CREATE TABLE bot.tests(
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    subject VARCHAR (50) NOT NULL,
    level VARCHAR (50) NOT NULL,
    topic VARCHAR (250) NOT NULL,
    test TEXT NOT NULL,
    result INTEGER NOT NULL CHECK (result BETWEEN 0 AND 100),
    time_finish TIMESTAMPTZ NOT NULL,
    user_id BIGINT NOT NULL REFERENCES bot.users(telegram_id)
);
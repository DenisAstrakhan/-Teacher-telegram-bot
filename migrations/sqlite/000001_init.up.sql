PRAGMA foreign_keys = ON;

CREATE TABLE teacher (
    telegram_id INTEGER PRIMARY KEY, 
    version INTEGER NOT NULL DEFAULT 1,
    telegram_name TEXT,            
    teacher_name TEXT NOT NULL
);

CREATE INDEX idx_teacher_teacher_name ON teacher(teacher_name);
CREATE INDEX idx_teacher_version ON teacher(version);

CREATE TABLE users (
    telegram_id INTEGER PRIMARY KEY,
    version INTEGER NOT NULL DEFAULT 1,
    telegram_name TEXT,
    full_name TEXT NOT NULL,
    teacher_id INTEGER NOT NULL,
    FOREIGN KEY (teacher_id) REFERENCES teacher(telegram_id) ON DELETE CASCADE
);

CREATE INDEX idx_users_teacher_id ON users(teacher_id);
CREATE INDEX idx_users_version ON users(version);

CREATE TABLE tests (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    version INTEGER NOT NULL DEFAULT 1,
    subject TEXT NOT NULL,
    level TEXT NOT NULL,
    topic TEXT NOT NULL,
    test TEXT NOT NULL,
    result INTEGER NOT NULL CHECK (result >= 0 AND result <= 100),
    time_finish DATETIME NOT NULL,  
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(telegram_id) ON DELETE CASCADE
);

CREATE INDEX idx_tests_user_id ON tests(user_id);
CREATE INDEX idx_tests_version ON tests(version);


CREATE TABLE IF NOT EXISTS teacher(
    telegram_id BIGINT PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    telegram_name VARCHAR(250),
    teacher_name VARCHAR(250) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_teacher_teacher_name ON teacher(teacher_name);
CREATE INDEX idx_teacher_version ON teacher(version);

CREATE TABLE IF NOT EXISTS users(
    telegram_id BIGINT PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    telegram_name VARCHAR(250),
    full_name VARCHAR(250) NOT NULL,
    teacher_id BIGINT NOT NULL,
    CONSTRAINT fk_users_teacher 
        FOREIGN KEY (teacher_id) 
        REFERENCES teacher(telegram_id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_users_teacher_id ON users(teacher_id);
CREATE INDEX idx_users_version ON users(version);

CREATE TABLE IF NOT EXISTS tests(
    id INT PRIMARY KEY AUTO_INCREMENT,
    version BIGINT NOT NULL DEFAULT 1,
    subject VARCHAR(50) NOT NULL,
    level VARCHAR(50) NOT NULL,
    topic VARCHAR(250) NOT NULL,
    test TEXT NOT NULL,
    result INTEGER NOT NULL,
    CONSTRAINT chk_tests_result CHECK (result BETWEEN 0 AND 100),
    time_finish DATETIME NOT NULL,
    user_id BIGINT NOT NULL,
    CONSTRAINT fk_tests_user 
        FOREIGN KEY (user_id) 
        REFERENCES users(telegram_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_tests_user_id ON tests(user_id);
CREATE INDEX idx_tests_version ON tests(version); 
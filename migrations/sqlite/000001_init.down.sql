DROP INDEX IF EXISTS idx_tests_version;
DROP INDEX IF EXISTS idx_tests_user_id;
DROP INDEX IF EXISTS idx_users_version;
DROP INDEX IF EXISTS idx_users_teacher_id;
DROP INDEX IF EXISTS idx_teacher_version;
DROP INDEX IF EXISTS idx_teacher_teacher_name;

DROP TABLE IF EXISTS tests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teacher;
DROP TABLE IF EXISTS schema_migrations;
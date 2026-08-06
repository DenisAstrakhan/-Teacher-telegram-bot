include .env
export
POSTGRES_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@bot-postgres:5432/${POSTGRES_DB}?sslmode=disable
MYSQL_URL=mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(bot-mysql:3306)/${MYSQL_DATABASE}
SQLITE_EXEC=docker exec -i bot-env-sqlite sqlite3 /data/database.db

sqlite-up:
	@docker compose up -d bot-sqlite
sqlite-down:
	@docker compose down bot-sqlite
mysql-up:
	@docker compose up -d bot-mysql
mysql-down:
	@docker compose down bot-mysql
# Дать права пользователю MySQL на создание схем и таблиц
mysql-grant-privileges:
	@echo "Выдача прав пользователю ${MYSQL_USER}..."
	docker exec -it bot-env-mysql mysql -u root -p${MYSQL_ROOT_PASSWORD} -e "GRANT ALL PRIVILEGES ON *.* TO '${MYSQL_USER}'@'%' WITH GRANT OPTION; FLUSH PRIVILEGES;"
	@echo "Права успешно выданы!"
	@docker exec -it bot-env-mysql mysql -u root -p${MYSQL_ROOT_PASSWORD} -e "SHOW GRANTS FOR '${MYSQL_USER}'@'%';"
env-up:
	@docker compose up -d bot-postgres

env-down:
	@docker compose down bot-postgres

env-cleanup:
	@read -p "Отчистить все volume окружения? Опасность утери данных. [y/n]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down bot-postgres && \
		sudo rm -rf out/pgdata && \
		echo "Файлы окружения отчищены"; \
	else \
		echo "Отчистка окружения отменена."; \
	fi;
env-port-forwarder:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-postgres-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-postgres-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm --user root bot-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations/postgres \
		-seq "$(seq)"; \
	sudo chown -R $(shell id -u):$(shell id -g) ./migrations/postgres

migrate-postgres-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-postgres-action action=up 1"; \
		exit 1; \
	fi; \
	docker compose run --rm --user root bot-postgres-migrate \
	-path /migrations/postgres \
	-database "$(POSTGRES_URL)" \
	"$(action)"

migrate-postgres-up:
	@$(MAKE) migrate-postgres-action action=up 

migrate-postgres-down:
	@$(MAKE) migrate-postgres-action action=down

migrate-mysql-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-mysql-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm --user root bot-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations/mysql \
		-seq "$(seq)"; \
	sudo chown -R $(shell id -u):$(shell id -g) ./migrations/mysql


migrate-mysql-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-mysql-action action=up 1"; \
		exit 1; \
	fi; \
	docker compose run --rm --user root bot-postgres-migrate \
	-path /migrations/mysql \
	-database "$(MYSQL_URL)" \
	"$(action)"

migrate-mysql-up:
	@$(MAKE) migrate-mysql-action action=up 

migrate-mysql-down:
	@$(MAKE) migrate-mysql-action action=down
# Создание файла базы данных
migrate-sqlite-create-db:
	@echo "📋 Создание файла базы данных на хосте..."
	@sudo mkdir -p ./out/sqlitedata
	@sudo chmod 777 ./out/sqlitedata
	@if [ ! -f ./out/sqlitedata/database.db ]; then \
		echo "  Файл не существует, создаю..."; \
		sudo touch ./out/sqlitedata/database.db; \
		sudo chmod 666 ./out/sqlitedata/database.db; \
		echo "  ✅ Файл создан: ./out/sqlitedata/database.db"; \
	else \
		echo "  ✅ Файл уже существует"; \
	fi
	@echo "  📊 Проверка:"
	@sudo ls -la ./out/sqlitedata/database.db

migrate-sqlite-create:
	@if [ -z "$(seq)" ]; then \
		echo "❌ Отсутствует параметр seq. Пример: make migrate-sqlite-create seq=init"; \
		exit 1; \
	fi; \
	@echo "📄 Создание миграции: $(seq)"; \
	mkdir -p ./migrations/sqlite; \
	NEXT_NUMBER=$$(printf "%06d" $$(($$(ls -1 ./migrations/sqlite/*.up.sql 2>/dev/null | wc -l) + 1))); \
	UP_FILE="./migrations/sqlite/$${NEXT_NUMBER}_$(seq).up.sql"; \
	DOWN_FILE="./migrations/sqlite/$${NEXT_NUMBER}_$(seq).down.sql"; \
	echo "-- +migrate Up" > $$UP_FILE; \
	echo "-- +migrate Down" > $$DOWN_FILE; \
	echo "✅ Созданы файлы:"; \
	echo "   $$UP_FILE"; \
	echo "   $$DOWN_FILE"

# Применить все миграции (ИСПРАВЛЕНО)
migrate-sqlite-up:
	@echo "⬆️  Применение миграций..."; \
	for f in ./migrations/sqlite/*.up.sql; do \
		if [ -f "$$f" ]; then \
			echo "  Применение: $$(basename $$f)"; \
			cat "$$f" | $(SQLITE_EXEC); \
			if [ $$? -eq 0 ]; then \
				echo "  ✅ $$(basename $$f)"; \
			else \
				echo "  ❌ Ошибка в $$(basename $$f)"; \
				exit 1; \
			fi; \
		fi; \
	done; \
	echo "✅ Все миграции применены"; \
	echo ""; \
	echo "📊 Таблицы:"; \
	docker exec bot-env-sqlite sqlite3 /data/database.db ".tables"

# Откатить последнюю миграцию
migrate-sqlite-down:
	@echo "⬇️  Откат последней миграции..."; \
	LAST_FILE=$$(ls -1 ./migrations/sqlite/*.down.sql 2>/dev/null | sort -n | tail -1); \
	if [ -n "$$LAST_FILE" ]; then \
		echo "  Откат: $$(basename $$LAST_FILE)"; \
		cat "$$LAST_FILE" | $(SQLITE_EXEC); \
		if [ $$? -eq 0 ]; then \
			echo "  ✅ Откат выполнен"; \
		else \
			echo "  ❌ Ошибка отката"; \
			exit 1; \
		fi; \
	else \
		echo "  ⚠️  Нет миграций для отката"; \
	fi

redis-up:
	docker compose up -d redis

redis-down:
	docker compose down redis
run-service:
	go run main.go

delet-image:
	docker rmi telegram-bot

bot-run:
	@docker compose up -d bot-postgres
	 docker compose up telegram-bot

bot-down:
	docker compose down telegram-bot

vpn-up:
	adguardvpn-cli connect -l DE

vpn-down:
	adguardvpn-cli disconnect

redis-clear:
	sudo rm -f ./out/redisdata/dump.rdb
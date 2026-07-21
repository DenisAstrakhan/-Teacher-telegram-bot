include .env
export
POSTGRES_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@bot-postgres:5432/${POSTGRES_DB}?sslmode=disable
MYSQL_URL=mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(bot-mysql:3306)/${MYSQL_DATABASE}

mysql-up:
	docker compose up -d bot-mysql
mysql-down:
	docker compose down bot-mysql
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

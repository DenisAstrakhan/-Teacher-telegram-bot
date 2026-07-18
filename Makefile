include .env
export

mysql-up:
	docker compose up -d bot-mysql
mysql-down:
	docker compose down bot-mysql
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

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm --user root bot-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"; \
	sudo chown -R $(shell id -u):$(shell id -g) ./migrations

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-action action=up 1"; \
		exit 1; \
	fi; \
	docker compose run --rm --user root bot-postgres-migrate \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@bot-postgres:5432/${POSTGRES_DB}?sslmode=disable \
	"$(action)"

migrate-up:
	@$(MAKE) migrate-action action=up 

migrate-down:
	@$(MAKE) migrate-action action=down

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

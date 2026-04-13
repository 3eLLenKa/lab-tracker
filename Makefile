-include .env
export

DB_DSN = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@db:5432/$(POSTGRES_DB)?sslmode=disable

up:
	docker compose up -d db
	docker compose run --rm --entrypoint goose migrator \
		-dir /migrations postgres "$(DB_DSN)" up
	docker compose up --build -d app

down:
	docker compose down

build:
	docker compose build app

logs:
	docker compose logs -f app

db-shell:
	docker compose exec db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

migrate-create:
	@if [ -z "$(name)" ]; then \
	 echo 'usage: make migrate-create name=add_column'; exit 1; \
	fi
	docker compose run --rm \
		--entrypoint goose migrator \
		-dir /migrations \
		create $(name) sql

db-up:
	docker compose up -d db
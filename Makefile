up:
	docker compose up -d

down:
	docker compose down

down-v:
	docker compose down -v

build:
	docker compose build

seed:
	docker compose run --rm seed

logs:
	docker compose logs -f server

dev:
	go run ./cmd/seed && go run ./cmd/server

help:
	@echo "make up       - start server + db"
	@echo "make down     - stop containers"
	@echo "make down-v   - stop and delete db data"
	@echo "make build    - build images"
	@echo "make seed     - run seed"
	@echo "make logs     - view logs"
	@echo "make dev      - run locally"
# Makefile for managing docker-compose commands

PROJECT_NAME=goldsteps

.PHONY: local

local:
	(cd server && go run main.go) & \
	(cd client && npm run dev) & \
	wait

up:
	docker-compose up -d

down:
	docker-compose down

build:
	docker-compose build --no-cache

restart:
	make down
	make up

logs:
	docker-compose logs -f

ps:
	docker-compose ps

exec-server:
	docker-compose exec server sh

exec-client:
	docker-compose exec client sh

clean:
	docker-compose down -v  # Remove volumes as well

destroy:
	docker-compose down --rmi all --volumes --remove-orphans

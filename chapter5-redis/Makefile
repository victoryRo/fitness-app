# database name
DB_NAME ?= fitness

# database type
DB_TYPE ?= postgres

# database username
DB_USER ?= expert

# database password
DB_PASSWORD ?= asuncionAutosostenida

# host IP
IP=127.0.0.1

PSQL_URL ?= $(DB_TYPE)://$(DB_USER):$(DB_PASSWORD)@$(IP):5432/$(DB_NAME)

# yaml file
SQL_C ?= ./sqlc.yaml

.PHONY : postgres_app create_db check_db generate redis recreate_redis run_project

postgres_app:
	docker run --name fitness_db \
	-e POSTGRES_USER=$(DB_USER) \
	-e POSTGRES_DB=$(DB_NAME) \
	-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
	-v $(PWD):/usr/share/chapter5-redis \
	-p 5432:5432 \
	-d postgres

# task to create database without typing it manually
create_db:
	docker exec -it fitness_db psql $(PSQL_URL) -c "\i /usr/share/chapter5-redis/db/schema.sql"

# check tables
check_db:
	docker exec -it fitness_db psql $(PSQL_URL) -c "\dt app.*"

# generating Go models with
generate:
	sqlc generate -f $(SQL_C)
	
# run redis container
redis:
	docker run --name local-redis \
	-v $(PWD)/redisdata:/data \
	-p 6379:6379 \
	-d redis

recreate_redis:
	docker stop local-redis || true && docker rm local-redis || true

run_project:
	$(MAKE) postgres_app && sleep 3 && $(MAKE) create_db && sleep 3 && $(MAKE) check_db && sleep 2 && $(MAKE) redis && sleep 5 && go run .


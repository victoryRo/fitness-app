# database name
DB_NAME ?= postgres

# database type
DB_TYPE ?= postgres

# database username
DB_USER ?= postgres

# database password
DB_PWD ?= SecretPassword

# psql URL
IP=`podman inspect --format '{{ .NetworkSettings.IPAddress }}' postgres-db`

PSQLURL ?= $(DB_TYPE)://$(DB_USER):$(DB_PWD)@$(IP):5432/$(DB_NAME)

# sqlc yaml file
SQLC_YAML ?= ./sqlc.yaml

.PHONY : db_up db_down teardown_recreate run_db generate build create_db

db_up:
	podman container run --name postgres-db \
	-e POSTGRES_PASSWORD=SecretPassword \
	-v $(PWD):/usr/share/chapter2 \
	-p 54320:5432 \
	-d $(DB_NAME)

run_db:
	podman exec -it postgres-db \
	psql -h localhost \
	-p 5432 \
	-U $(DB_NAME) \
	-d $(DB_NAME)

db_down:
	podman container stop postgres-db
	podman container rm postgres-db


teardown_recreate: postgresdown postgresup
	# sleep 5
	# $(MAKE) createdb

generate:
	@echo "Generating Go models with sqlc "
	sqlc generate -f $(SQLC_YAML)

build:
	@echo "Building database main sample app"
	go build -o sampledb .

# task to create database without typing it manually
create_db:
	podman run -v $(PWD):/usr/share/chapter2 -it --rm jbergknoff/postgresql-client -p 5432:5432 $(PSQLURL) -c "\i /usr/share/chapter2/db/schema.sql"


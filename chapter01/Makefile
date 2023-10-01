# database name
DB_NAME ?= postgres

# database type
DB_TYPE ?= postgres

# database username
DB_USER ?= postgres

# database password
DB_PWD ?= mysecretpassword

# psql URL
IP=127.0.0.1

# database path
PSQLURL ?= $(DB_TYPE)://$(DB_USER):$(DB_PWD)@$(IP):5432/$(DB_NAME)

# sqlc yaml file
SQLC_YAML ?= ./sqlc.yaml

.PHONY : postgresup postup postgresdown psql createdb teardown_recreate generate

postgresup:
	@echo "Run the database container, volume, password ok"
	podman container run -d --name postgres-db \
	-e POSTGRES_PASSWORD=$(DB_PWD) \
	-v $(PWD):/usr/share/chapter01 \
	-p 54320:5432 \
	$(DB_NAME)

postup:
	podman container run --name store \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_DB=postgres \
	-v $(PWD):/usr/share/chapter01 \
    -p 54320:5432 \
    -d postgres:latest

postgresdown:
	podman container stop postgres-db || true && podman container rm postgres-db || true

psql:
	podman exec -it postgres-db psql $(PSQLURL)

# task to create database without typing it manually
createdb:
	podman exec -it postgres-db psql $(PSQLURL) -c "\i /usr/share/chapter01/db/schema.sql"

teardown_recreate: postgresdown postgresup
	sleep 5
	$(MAKE) createdb
	
generate:
	@echo "Generating Go models with sqlc"
	sqlc generate -f $(SQLC_YAML)

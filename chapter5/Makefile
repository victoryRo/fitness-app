# database name
DB_NAME ?= fitness

# database type
DB_TYPE ?= postgres

# database username
DB_USER ?= victory

# database password
DB_PWD ?= yosoyelcamino

# psql URL
IP=127.0.0.1

PSQL_URL ?= $(DB_TYPE)://$(DB_USER):$(DB_PWD)@$(IP):5432/$(DB_NAME)

# sqlc yaml file
SQLC_YAML ?= ./sqlc.yaml

.PHONY : postgres_up postgres_down psql create_db check_db teardown_recreate generate

# we run the database with a volume
postgres_up:
	podman run --name test-postgres \
	-e POSTGRES_USER=victory \
	-e POSTGRES_DB=fitness \
	-e POSTGRES_PASSWORD=yosoyelcamino \
	-v $(PWD):/usr/share/chapter5 \
	-p 5432:5432 \
	-d postgres

# task to create database without typing it manually
create_db:
	echo $(PWD)
	podman exec -it test-postgres psql $(PSQL_URL) -c "\i /usr/share/chapter5/db/schema.sql"

# check database table
check_db:
	podman exec -it test-postgres psql $(PSQL_URL) -c "\dt app.*"

# we use psql to interact with the database
psql:
	podman exec -it test-postgres psql -h localhost -p 5432 -U victory -d fitness

# remove the container from the database
postgres_down:
	podman stop test-postgres  || true && 	podman rm test-postgres || true

# destroy and recreate the database
teardown_recreate: postgres_down postgres_up
	sleep 5
	$(MAKE) createdb

# generates the db interface
# db models that are struct
# sql query methods
generate:
	@echo "Generating Go models with sqlc "
	sqlc generate -f $(SQLC_YAML)


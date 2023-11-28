# database name
DB_NAME ?= fullstackdb

# database type
DB_TYPE ?= postgres

# database username
DB_USER ?= victory

# database password
DB_PASS ?= programmingExpert

# ip
IP ?= 127.0.0.1

# psql URL
PSQL_URL ?= $(DB_TYPE)://$(DB_USER):$(DB_PASS)@$(IP)/$(DB_NAME)

.PHONY: db_up migrate create_db check_db generate

# run container db with credentials
db_up:
	docker run --name project-db \
	-e POSTGRES_DB=$(DB_NAME) \
	-e POSTGRES_USER=$(DB_USER) \
	-e POSTGRES_PASSWORD=$(DB_PASS) \
	-v $(PWD):/usr/share/chapter6 \
	-p 5432:5432 \
	-d postgres:11.10-alpine

# migrate db
# create db at the same time
migrate:
	migrate -path ./migrations -database $(PSQL_URL)?sslmode=disable up

# other way to create DB
create_db:
	 docker exec -it project-db psql $(PSQL_URL) -c "\i /usr/share/chapter6/migrations/000001_schema.up.sql"

# check db
check_db:
	docker exec -it project-db psql $(PSQL_URL) -c "\dt"
	
# generating Go models sqlc
generate:
	@echo "Generating Go models with sqlc "
	go generate

# query fail bad password
# curl -H 'Content-Type: application/json' 0.0.0.0:9002/login -d '{"username":"user@user","password":"passwordBad"}' -v

# Login
#curl -H 'Content-Type: application/json' 0.0.0.0:9002/login -d '{"username":"user@user","password":"password"}'  -v

# Get a list of previous workouts and their sets
#curl -X GET -H 'Content-Type: application/json' 0.0.0.0:9002/workout --cookie 'session-name=MTY0NTI2MTczNXxEdi1CQkFFQ180SUFBUkFCRUFBQVJQLUNBQUlHYzNSeWFXNW5EQk1BRVhWelpYSkJkWFJvWlc1MGFXTmhkR1ZrQkdKdmIyd0NBZ0FCQm5OMGNtbHVad3dJQUFaMWMyVnlTVVFGYVc1ME5qUUVBZ0FDfMwSOVjl_-nwIrsRVE1b5Q2ss-kd_RyObfoO-HlrVP0j;'

# Add a new workout (with no entries)
#curl -X POST -H 'Content-Type: application/json' 0.0.0.0:9002/workout --cookie 'session-name=MTY0NTI2MTczNXxEdi1CQkFFQ180SUFBUkFCRUFBQVJQLUNBQUlHYzNSeWFXNW5EQk1BRVhWelpYSkJkWFJvWlc1MGFXTmhkR1ZrQkdKdmIyd0NBZ0FCQm5OMGNtbHVad3dJQUFaMWMyVnlTVVFGYVc1ME5qUUVBZ0FDfMwSOVjl_-nwIrsRVE1b5Q2ss-kd_RyObfoO-HlrVP0j;'

## Delete a workout
#curl -X DELETE -H 'Content-Type: application/json' 0.0.0.0:9002/workout/1 --cookie 'session-name=MTY0NTI2MTczNXxEdi1CQkFFQ180SUFBUkFCRUFBQVJQLUNBQUlHYzNSeWFXNW5EQk1BRVhWelpYSkJkWFJvWlc1MGFXTmhkR1ZrQkdKdmIyd0NBZ0FCQm5OMGNtbHVad3dJQUFaMWMyVnlTVVFGYVc1ME5qUUVBZ0FDfMwSOVjl_-nwIrsRVE1b5Q2ss-kd_RyObfoO-HlrVP0j;'

## Add a set to a workout
#curl -X POST -H 'Content-Type: application/json' 0.0.0.0:9002/workout/5 --cookie 'session-name=MTY0NTI2MTczNXxEdi1CQkFFQ180SUFBUkFCRUFBQVJQLUNBQUlHYzNSeWFXNW5EQk1BRVhWelpYSkJkWFJvWlc1MGFXTmhkR1ZrQkdKdmIyd0NBZ0FCQm5OMGNtbHVad3dJQUFaMWMyVnlTVVFGYVc1ME5qUUVBZ0FDfMwSOVjl_-nwIrsRVE1b5Q2ss-kd_RyObfoO-HlrVP0j;' -d '{"exercise_name": "Barbell Rows", "weight":700'}

# Database
MYSQL_USER ?= notSuperUser
MYSQL_PASSWORD ?= 123456
MYSQL_ADDRESS ?= docker.for.mac.localhost:3306
#MYSQL_ADDRESS ?= localhost:3306
MYSQL_DATABASE ?= go-api

.PHONY: migrate-up
migrate-up:
	@docker exec go-api goose -dir ./migrations/sql mysql \
	 "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)?parseTime=true" up

.PHONY: migrate-down
migrate-down:
	@docker exec go-api goose -dir ./migrations/sql mysql \
	 "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)?parseTime=true" down

#docker exec go-api goose -dir ./migrations/sql create post sql
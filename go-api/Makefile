# Database
MYSQL_USER ?= notSuperUser
MYSQL_PASSWORD ?= 123456
MYSQL_ADDRESS ?= 127.0.0.1:3306
MYSQL_DATABASE ?= go-api

.PHONY: migrate-prepare
migrate-prepare:
	@go get -tags 'mysql' -u github.com/golang-migrate/migrate/v4/cmd/migrate
	@go build -a -o ./bin/migrate -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate

	# @go get -u github.com/golang-migrate/migrate/v4
	# @go build -a -o ./bin/migrate -tags 'mysql' github.com/golang-migrate/migrate/v4/cli

.PHONY: migrate-up
migrate-up:
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=migrations/sql up

.PHONY: migrate-down
migrate-down:
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=migrations/sql down
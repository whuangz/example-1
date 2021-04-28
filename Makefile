# Database
MYSQL_USER ?= notSuperUser
MYSQL_PASSWORD ?= 123456
MYSQL_ADDRESS ?= docker.for.mac.localhost:3306
#MYSQL_ADDRESS ?= localhost:3306
MYSQL_DATABASE ?= go-api

PWD = $(shell pwd)
ACCTPATH = $(PWD)/go-api/config

.PHONY: create-keypair
create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(ACCTPATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(ACCTPATH)/rsa_private_$(ENV).pem -pubout -out $(ACCTPATH)/rsa_public_$(ENV).pem

.PHONY: migrate-up
migrate-up:
	@docker exec go-api goose -dir ./migrations/sql mysql \
	 "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)?parseTime=true" up

.PHONY: migrate-down
migrate-down:
	@docker exec go-api goose -dir ./migrations/sql mysql \
	 "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)?parseTime=true" down

.PHONY: test
test:
	@docker exec go-api go test -v -short -race ./...
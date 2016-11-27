APP ?= mff
ENV_FILE ?= .env

run:
	./bin/$(APP) -env ./custom/$(ENV_FILE)

build:
	go build -o ./bin/$(APP)

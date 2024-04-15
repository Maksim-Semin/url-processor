include ./resource/config.toml

ifeq ($(DB), "postgres")
build:
	docker-compose -f ./Docker/docker-compose.yml build
run:
	docker-compose -f ./Docker/docker-compose.yml up -d
	timeout 2

migrations:
	goose -dir ./pkg/storage/migrations postgres "host=localhost port=5432 user=postgres dbname=urls password=password sslmode=disable" up

else
build:
	docker build -t service -f ./Docker/DockerfileService .
	timeout 1
run:
	docker run -p 50051:50051 service
migrations:
	@echo "No migrations needed for this configuration"
endif

all: build run migrations

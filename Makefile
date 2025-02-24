all: build run

build:
	docker-compose build --no-cache

run:
	docker-compose up -d

inmmem: 
	go run cmd/app/main.go

restart:
	docker-compose down -v
	docker-compose build --no-cache
	docker-compose up -d

test: 
	go test ./... -v

stop:
	docker-compose down -v

APP_NAME=pocket-counter

build:
	go build -o $(APP_NAME) main.go

run:
	go run main.go


build-docker:
	docker build -t pviotti/$(APP_NAME) .

run-docker:
	docker run --env-file .env -p 8080:8080 -v ./data:/app/data pviotti/$(APP_NAME)
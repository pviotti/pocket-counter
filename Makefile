
APP_NAME=pocket-counter

build:
	go build -o $(APP_NAME) main.go

run:
	go run main.go


build-docker:
	docker build -t pviotti/$(APP_NAME) .

run-docker:
	docker run --network host --env-file .env -v ./data:/app/data pviotti/$(APP_NAME)

sh-docker:
	docker run -it --network host --env-file .env -v ./data:/app/data pviotti/$(APP_NAME) /bin/sh
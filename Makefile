APP_NAME=ldt-server

.PHONY: all

all: clean build run

.PHONY: clean

clean:
	@rm -rf ./ldt
	@echo "[✔️] Clean complete!"

.PHONY: build

build: clean
	@go build -o $(APP_NAME) ./cmd
	@echo "[✔️] Build complete!"

.PHONY: test

test:
	@docker-compose up -d
	@echo "[✔️] Docker containers started!"
	@-go test -v ./...
	@echo "[✔️] Tests terminated!"
	@docker-compose down
	@echo "[✔️] Docker containers stopped!"

.PHONY: run

run:
	@./ldt

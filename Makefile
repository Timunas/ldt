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

.PHONY: run

run:
	@./ldt

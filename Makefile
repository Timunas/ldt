APP_NAME=ldt

.PHONY: all

all: clean build run

.PHONY: clean

clean:
	@rm -rf ./ldt
	@echo "[✔️] Clean complete!"

.PHONY: build

build:
	@go build -o $(APP_NAME) ./cmd/app
	@echo "[✔️] Build complete!"

.PHONY: run

run:
	@./ldt

SRC = ./cmd/payr/main.go
OUTPUT = ./build/payr
ENV = ./.env

.PHONY: build run clean

build:
	go build -o $(OUTPUT) $(SRC)

run:
	. $(ENV) && go run $(SRC)

clean:
	rm ./build/*

NAME = payr
SRC = ./cmd/$(NAME)/$(NAME).go
OUTPUT = ./build/$(NAME)

.PHONY: build run clean fmt

build:
	go build -o $(OUTPUT) $(SRC)

run:
	go run $(SRC)

clean:
	rm ./build/*

fmt:
	go fmt $(NAME)/...
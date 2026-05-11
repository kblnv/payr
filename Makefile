PAYREM_PATH = ./cmd/payrem/main.go
ENV_PATH = ./.env

.PHONY: build run clean

build:
	go build -o ./build/payrem $(PAYREM_PATH)

run:
	. $(ENV_PATH) && go run $(PAYREM_PATH)

clean:
	rm ./build/*
.PHONY: build run clean

build:
	go build -o ./build/payrem cmd/payrem/main.go

run:
	go run cmd/payrem/main.go

clean:
	rm ./build/*
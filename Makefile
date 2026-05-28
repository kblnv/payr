NAME = payr
SRC = ./cmd/$(NAME)/$(NAME).go
PLUGINS_SRC = ./plugins
OUTPUT = ./build/$(NAME)
PLUGINS_OUTPUT = ./build/plugins

.PHONY: build build_plugins run clean fmt

build:
	go build -o $(OUTPUT) $(SRC)

build_plugins:
	@for dir in $(PLUGINS_SRC)/*/ ; do \
		name=$$(basename $$dir); \
		go build -buildmode=plugin -o $(PLUGINS_OUTPUT)/$$name.so $$dir; \
	done

run:
	go run $(SRC)

clean:
	rm ./build/*

fmt:
	go fmt $(NAME)/...
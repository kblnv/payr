NAME = payr
SRC = ./cmd/$(NAME)/$(NAME).go
BUILD = ./build
PLUGINS_SRC = ./plugins
OUTPUT = $(BUILD)/$(NAME)
PLUGINS_OUTPUT = $(BUILD)/plugins

.PHONY: core plugins run clean fmt

core:
	go build -o $(OUTPUT) $(SRC)

plugins:
	@for dir in $(PLUGINS_SRC)/*/ ; do \
		name=$$(basename $$dir); \
		go build -buildmode=plugin -o $(PLUGINS_OUTPUT)/$$name.so $$dir; \
	done

clean:
	rm ./build/*

fmt:
	go fmt $(NAME)/...
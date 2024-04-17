# Makefile

.PHONY: build clean install uninstall run

BINARY_NAME=wcgo
INSTALL_PATH=/usr/local/bin

build:
	go build -o $(BINARY_NAME) main.go

clean:
	rm -f $(BINARY_NAME)

install: build
	sudo mv $(BINARY_NAME) $(INSTALL_PATH)

uninstall:
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)

run:
	@$(MAKE) build > /dev/null
	@./$(BINARY_NAME) $(ARGS)
	@$(MAKE) clean > /dev/null
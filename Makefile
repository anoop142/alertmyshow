
BIN="alertmyshow"
BUILD_FLAGS="-s -w"

.PHONY: all
all: build

.PHONY:	build
build:
	CGO_ENABLED=0 go build -o $(BIN) -ldflags $(BUILD_FLAGS) ./cmd/cli

.PHONY: run
run: build
	./$(BIN)

.PHONY:	clean
clean:
	rm -f $(BIN)


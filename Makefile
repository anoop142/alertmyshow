
BIN="alertmyshow"

VERSION_MAJOR ?= 0
VERSION_MINOR ?= 2
VERSION_BUILD ?= 0
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

GO_LDFLAGS :='-s -w
GO_LDFLAGS += -X main.version=$(VERSION)
GO_LDFLAGS +='



.PHONY: all
all: build

.PHONY:	build
build:
	CGO_ENABLED=0 go build -o $(BIN) -ldflags $(GO_LDFLAGS) ./cmd/alertmyshow

.PHONY: run
run: build
	./$(BIN)

.PHONY:	clean
clean:
	rm -f $(BIN)


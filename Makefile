
BIN="alertmyshow"

VERSION_MAJOR ?= 0
VERSION_MINOR ?= 3
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

.PHONY:	cross
cross:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN) -ldflags $(GO_LDFLAGS) ./cmd/alertmyshow
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BIN).exe -ldflags $(GO_LDFLAGS) ./cmd/alertmyshow


.PHONY: run
run: build
	./$(BIN)

.PHONY:	clean
clean:
	rm -f $(BIN)


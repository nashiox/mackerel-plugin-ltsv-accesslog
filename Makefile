REPO = mackerel-plugin-ltsv-accesslog
BIN = $(REPO)

all: clean test build

test:
	go test ./...

build:
	go build -o bin/$(BIN) main.go

clean:
	rm -rf bin

.PHONY: test build clean

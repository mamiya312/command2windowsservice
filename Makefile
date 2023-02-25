VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS := $(LDFLAGS) -X main.commit=$(COMMIT) -X main.branch=$(BRANCH)
ifdef VERSION
	LDFLAGS += -X main.version=$(VERSION)
endif
test:
	go test -v ./...
build:
	GOOS=windows go build --ldflags "$(LDFLAGS)" -o bin/command2windowsservice.exe cmd/command2windowsservice/*.go

BINARY = purestorage-sensor
VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-s -w -X main.version=$(VERSION)"
GO = go

.PHONY: build build-windows build-linux build-all clean test vet

build:
	$(GO) build $(LDFLAGS) -o $(BINARY) ./cmd/purestorage-sensor/

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BINARY).exe ./cmd/purestorage-sensor/

build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BINARY) ./cmd/purestorage-sensor/

build-all: build-windows build-linux

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

clean:
	rm -f $(BINARY) $(BINARY).exe

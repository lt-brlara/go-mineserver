BINARY_NAME = main
GO = go
GOFLAGS = -x -n
GOENTRYPOINT = cmd/main.go

build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(GOENTRYPOINT)

run:
	$(GO) run $(GOENTRYPOINT)

test:
	$(GO) test ./...

lint:
	golangci-lint run

fmt:
	$(GO) fmt ./...

.PHONY : clean
clean :
	rm -f $(EXECUTABLE)

.PHONY: default
default: build

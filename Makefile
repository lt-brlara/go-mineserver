.DEFAULT_GOAL = build

BUILDFLAGS = -x
ENTRYPOINT = cmd/main.go
EXECUTABLE = main

build :
	go build -o $(EXECUTABLE) $(BUILDFLAGS) $(ENTRYPOINT)

run :
	go run $(ENTRYPOINT)

test :
	go test -v ./...

.PHONY : clean
clean :
	rm -f $(EXECUTABLE)

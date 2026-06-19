.PHONY: all build test clean run

all: build test

build:
	go build -o chpass-web

test:
	go test -v -short ./...

clean:
	rm -f chpass-web

run:
	go run .

.PHONY: docker
docker:
	docker build -t chpass-web .
.PHONY: build test lint fmt clean

build:
	go build ./cmd/discover

test:
	go test ./...

lint:
	go vet ./...

fmt:
	gofmt -w .

clean:
	rm -f discover

.PHONY: build test lint fmt clean

build:
	go build ./cmd/docnav

test:
	go test ./...

lint:
	go vet ./...

fmt:
	gofmt -w .

clean:
	rm -f docnav

.PHONY: fmt lint test check

export GOOS=windows

fmt:
	gofmt -w .

lint:
	test -z "$$(gofmt -l .)"
	go vet ./...

test:
	go test ./...

check: lint test

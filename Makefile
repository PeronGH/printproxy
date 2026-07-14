.PHONY: fmt lint test check

fmt:
	gofmt -w .

lint:
	test -z "$$(gofmt -l .)"
	GOOS=windows go vet ./...
	GOOS=darwin go vet ./...

test:
	go test ./...

check: lint test

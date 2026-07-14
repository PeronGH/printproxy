.PHONY: fmt fmt-check vet test check

fmt:
	gofmt -w .

fmt-check:
	test -z "$$(gofmt -l .)"

vet:
	go vet ./...

test:
	go test ./...

check: fmt-check vet test

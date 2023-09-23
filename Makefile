default: build

build: vendor fmt lint
	rm -rf later
	go build -mod vendor -o later ./cmd/later/.

tidy:
	go mod tidy

vendor: tidy
	rm -rf vendor
	go mod vendor

fmt:
	go fmt ./...

lint: fmt
	golangci-lint run

.PHONY: default $(MAKECMDGOALS)
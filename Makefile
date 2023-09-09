default: build

build:
	rm -rf later
	go build -mod vendor -o later ./cmd/later/.

tidy:
	go mod tidy

vendor: tidy
	rm -rf vendor
	go mod vendor

fmt:
	go fmt ./...

lint:
	golangci-lint run

.PHONY: default $(MAKECMDGOALS)
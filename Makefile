default: build

pipeline: vendor lint test build

build: del_binary
	go build -mod vendor -o later ./cmd/later/.

test:
	go test -v ./...

tidy:
	go mod tidy

vendor: tidy
	rm -rf vendor
	go mod vendor

fmt:
	go fmt ./...

lint: fmt
	golangci-lint run

del_binary:
	rm -rf later

.PHONY: default $(MAKECMDGOALS)
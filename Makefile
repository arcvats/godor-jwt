.DEFAULT_GOAL := build

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

vet:
	go vet ./...

audit:
	go list -json -m all | nancy sleuth --exclude-vulnerability-file .nancy-ignore

outdated:
	go list -u -m -json all | go-mod-outdated -update -direct

weight:
	goweight ./...

test:
	go test -race -v ./...

coverage:
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

build: fmt lint vet test
	go build -o bin/ ./...
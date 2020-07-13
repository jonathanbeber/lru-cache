all: test

test: lint format vet
	@echo "running tests"
	@go test -cover -v ./...

lint:
	@command -v golangci-lint >/dev/null; if [ $$? -ne 0 ]; then \
		echo "installing golangci-lint tool"; \
		go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0; \
	fi
	@echo "running golangci-lint"
	@golangci-lint run
	@command -v golint >/dev/null; if [ $$? -ne 0 ]; then \
		echo "installing golint tool"; \
		go get -u golang.org/x/lint/golint; \
	fi
	@echo "running golint on ./examples/factorial"
	@golint -set_exit_status examples/factorial
	@echo "running golint on ./lru/"
	@golint -set_exit_status lru/

format:
	@echo "running gofmt on ./lru"
	@gofmt -w lru
	@echo "running gofmt on ./examples/factorial"
	@gofmt -w examples/factorial

vet:
	@echo "running go vet"
	@go vet ./...

factorial:
	@echo "running factorial example"
	@go run ./examples/factorial/main.go

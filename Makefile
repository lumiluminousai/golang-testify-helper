run-test:
	golangci-lint run
	go test -race -cover ./... -count=1 -failfast

run:
	go run cmd/main.go $(DIR)

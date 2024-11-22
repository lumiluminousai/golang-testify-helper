run-test:
	golangci-lint run
	go test -race -cover ./... -count=1 -failfast

run:
	go run cmd/main.go $(DIR)

#make run DIR=<full path to golang project>
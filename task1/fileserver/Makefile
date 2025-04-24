lint:
	golangci-lint run -E goimports,gocritic -v ./...

tools:
	go install golang.org/x/tools/cmd/goimports@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)


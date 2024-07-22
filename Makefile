install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45

linter: install-linter
	golangci-lint run --out-format=github-actions --config .golangci.toml -v ./...

fix-linter: install-linter
	golangci-lint run --out-format=github-actions --config .golangci.toml -v --fix 
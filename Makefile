build:
	@go build  -o bin/main cli/main.go
run: build
	@ ./bin/mai
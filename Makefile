grpcconfig :
		go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
 		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 \
 		export PATH="$PATH:$(go env GOPATH)/bin"


proto :
		protoc proto/*.proto\
		--go_out=.\
		--go-grpc_out=.\
		--go_opt=paths=source_relative\
		--go-grpc_opt=paths=source_relative\
		--proto_path=.
build:
	@go build  -o bin/main cli/main.go
run: build
	@ ./bin/main



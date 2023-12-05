dockersdk:
		go get github.com/docker/docker/client

grpcconfig :
		go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
 		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 \
 		export PATH="$PATH:$(go env GOPATH)/bin"


grpc :
		protoc proto/*.proto\
		--go_out=.\
		--go-grpc_out=.\
		--go_opt=paths=source_relative\
		--go-grpc_opt=paths=source_relative\
		--proto_path=.
apply: 
	@ go run cli/cmd/main.go apply
drop: 
	@ go run cli/cmd/main.go drop
update: 
	@ go run cli/cmd/main.go update
build:
	@sudo go build  -o bin/core miniCloudCore/cmd/main.go
run: build
	@sudo  ./bin/core




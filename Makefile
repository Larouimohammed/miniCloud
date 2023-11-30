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
buildcli:
	@go build  -o bin/main1 cli/main.go
cli: buildcli
	@ ./bin/main1

build:
	@sudo go build  -o bin/main miniCloudCore/main.go
run: build
	@sudo  ./bin/main
buildpr:
	@sudo go build  -o bin/main2 provisioner/main.go
provision: buildpr
	@sudo  ./bin/main2



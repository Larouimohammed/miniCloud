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
watch: 
	@ go run cli/cmd/main.go watch
build:
	@sudo go build  -o bin/core miniCloudCore/cmd/main.go
run: build
	@sudo  ./bin/core
consul:
    
	@ sudo docker run -d -p 8500:8500 -p 8600:8600/udp --name=server consul:1.15.4 agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
	
	@ sudo docker run --name=client consul:1.15.4 agent -node=client-1 -retry-join=172.17.0.2

dockersdk:

	@ go get github.com/docker/docker/client

grpcconfig :
				
	@ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
 	
	@ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 \
 	
	@ export PATH="$PATH:$(go env GOPATH)/bin"
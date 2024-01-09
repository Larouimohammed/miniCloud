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
consulinstall:
	@ sudo docker run -d --name=consul-server  -p 8600:8600/udp -p 8600:8600/tcp -p 8500:8500/udp -p 8500:8500/tcp -v /home/ubuntu/miniCloud/miniCloud:/consul/config consul:1.15.4 agent -server -bootstrap-expect=1 -ui -node=server-1 -config-file=/consul/config/consul-config.json -client=0.0.0.0
	@ sudo docker run -d --name=client consul:1.15.4 agent -node=client-1 -retry-join=172.17.0.2
	@ go get github.com/hashicorp/consul/api


dockersdk:

	@ go get github.com/docker/docker/client

grpcconfig :
	@ apt install -y protobuf-compiler			
	@ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
 	
	@ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 \
 	
	@ export PATH="$PATH:$(go env GOPATH)/bin"

goconfig:
# @ ech	o "export GOPATH=$HOME/work" >> ~/.profile 
	
# @ echo "export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin" >> ~/.profile
    
	@ source ~/.bashrc
	@ source ~/.profile

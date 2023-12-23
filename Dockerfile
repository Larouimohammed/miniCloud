# builder image
FROM golang:latest
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/core  miniCloudCore/cmd/main.go



# executable
ENTRYPOINT [ "./bin/core" ]
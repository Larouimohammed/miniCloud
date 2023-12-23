# builder image
FROM golang:1.19-alpine3.16 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/core  miniCloudCore/cmd/main.go

# generate clean, final image for end users
FROM alpine:3.16
COPY --from=builder /build/bin/core .

# executable
ENTRYPOINT [ "./bin/core" ]
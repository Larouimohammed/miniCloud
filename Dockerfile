# builder image
FROM golang:1.19-alpine3.16 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -o main  miniCloudCore/cmd/main.go


# generate clean, final image for end users
FROM alpine:3.16
COPY --from=builder /build/laroui-agent .

# executable
ENTRYPOINT [ "./main" ]
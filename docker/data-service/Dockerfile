FROM golang:alpine as builder

LABEL maintainer="Morchev Georgiy <morochev.g@gmail.com>"

WORKDIR /app
COPY data-service/go.mod .
COPY data-service/go.sum .
RUN go mod download
COPY data-service/ .

RUN go install cmd/data-service/main.go

FROM alpine:3.12

COPY --from=builder /go/bin/main /usr/local/bin/data-service

ENTRYPOINT ["/usr/local/bin/data-service"]
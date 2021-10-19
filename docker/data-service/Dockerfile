FROM golang:alpine as builder

LABEL maintainer="Morchev Georgiy <morochev.g@gmail.com>"

COPY data-service/go.mod .
COPY data-service/go.sum .
RUN go mod download
COPY data-service/ .

RUN go install cmd/data-service

FROM alpine:3.12

COPY --from=builder /go/bin/data-service /usr/local/bin/data-service

ENTRYPOINT ["/usr/local/bin/data-service"]
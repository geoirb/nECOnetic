FROM golang:alpine as builder

LABEL maintainer="Morchev Georgiy <morochev.g@gmail.com>"

WORKDIR /app
COPY data-service/go.mod .
COPY data-service/go.sum .
RUN go mod download
COPY . .

WORKDIR /app/data-service/

RUN go install tools/migration/temperature/main.go

FROM alpine:3.12

COPY --from=builder /go/bin/main /usr/local/bin/temperature
COPY --from=builder /app/dataset/temperature /dataset

ENTRYPOINT ["/usr/local/bin/temperature"]
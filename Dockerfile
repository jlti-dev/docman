FROM golang:1.17 as builder
WORKDIR /app
COPY go /app
RUN go mod download && CGO_ENABLED=1 GOOS=linux GOFLAGS=-mod=mod go build -a -installsuffix cgo -o main .

FROM debian:latest
WORKDIR /docker
COPY --from=builder /app/main /docker/main
EXPOSE 8080

copy public /docker/public
CMD ["/docker/main"]


FROM golang:1.17-alpine as builder

WORKDIR /go/src/github.com/etiennecoutaud/backendstore-service
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v

FROM gcr.io/distroless/static-debian10
COPY --from=builder /go/src/github.com/etiennecoutaud/backendstore-service/backendstore-service /
ENTRYPOINT ["/backendstore-service"]
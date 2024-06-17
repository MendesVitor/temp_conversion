FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .
RUN go test -v

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]

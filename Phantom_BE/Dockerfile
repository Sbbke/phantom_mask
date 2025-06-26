FROM golang:1.24.4-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o PhantomBE main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/PhantomBE .
EXPOSE 8080
CMD ["./PhantomBE"]
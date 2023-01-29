
# Build stage
FROM golang:1.18-alpine AS builder
WORKDIR /app
RUN mkdir /images
COPY . .
RUN go build -o main main.go

EXPOSE 8080
CMD ["/app/main"]

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080
CMD ["/app/main"]
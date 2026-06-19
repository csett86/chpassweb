# Build stage
FROM golang:1.20 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o chpass-web

# Run stage
FROM alpine:latest

RUN apk add --no-cache pam pam-dev libgcc

WORKDIR /app
COPY --from=builder /app/chpass-web .
COPY templates/ ./templates/

USER root

EXPOSE 8080

CMD ["./chpass-web"]
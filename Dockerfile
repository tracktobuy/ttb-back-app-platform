# Stage 1: Build
FROM golang:1.26-alpine AS builder
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

# Stage 2: Runtime
FROM gcr.io/distroless/static:nonroot

ENV MONGO_URI=""
ENV MONGO_DB=""
ENV API_SERVER_PORT=8080


WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
ENTRYPOINT ["/app/server"]
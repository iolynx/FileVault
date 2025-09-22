# --- Build Stage ---
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Copy & download dependency files
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

# --- Final Stage ---
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/server .

# Expose app port
EXPOSE 8080

CMD ["/app/server"]

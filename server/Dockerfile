# --- Builder Stage ---
FROM golang:1.23rc2-alpine AS builder

WORKDIR /app

# Install SQLite development libraries
RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o main ./main.go

# --- Runtime Stage ---
FROM alpine:latest

WORKDIR /app

# Install SQLite runtime libraries
# Note: The SQLite runtime library is required for CGO-enabled binaries
# and is not included in the base Alpine image.
# This is necessary for the SQLite database to work properly.
# The SQLite development libraries are not needed in the final image.
RUN apk add --no-cache sqlite-libs
RUN apk add --no-cache sqlite

COPY --from=builder /app/main .
COPY stock_master_data ./stock_master_data

EXPOSE 8080

CMD ["./main"]

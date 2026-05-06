FROM golang:1.22-alpine AS builder

# Install build dependencies for C-based libraries (like webp)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /image-worker main.go

FROM alpine:latest
# Needed for the webp library to function in the final image
RUN apk add --no-cache libwebp
WORKDIR /root/
COPY --from=builder /image-worker .
CMD ["./image-worker"]
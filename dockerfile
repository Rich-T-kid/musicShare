# ---- Build Stage ----
FROM golang:1.22 AS builder
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies (if go.mod or go.sum changed)
RUN go mod tidy && go mod download

# Copy the vendor folder
COPY vendor ./vendor

# Copy the rest of the application code
COPY . .

# Build the application using the vendor folder
RUN go build -mod=vendor -o server

EXPOSE 80
CMD ["./server"]

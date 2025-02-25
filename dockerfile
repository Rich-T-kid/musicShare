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
# python env
FROM python:3.10.12

WORKDIR /app

COPY --from=builder /app/server /app/server

# Copy Python gRPC files
COPY reccommendations/grpc /app/grpc

# Copy the specific requirements file
COPY reccommendations/grpc/requirements.txt /app/reccommendations/grpc/requirements.txt

# Install Python dependencies
RUN pip install --no-cache-dir -r /app/grpc/requirements.txt

EXPOSE 80
CMD ["./server"]

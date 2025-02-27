# ------------------------------------------------
# 1) Builder Stage
# ------------------------------------------------
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for dependency resolution
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy vendor if you have a vendor/ directory
COPY vendor ./vendor

# Copy the rest of your Go source code
COPY . .

# Build the Go application
RUN go build -mod=vendor -o /server .

# ------------------------------------------------
# 2) Final Stage (Python + our Go binary)
# ------------------------------------------------
FROM python:3.10.12

# Set a working directory in the final container
WORKDIR /app

# Copy the built Go server from the builder stage
COPY --from=builder /server /app/server

# Copy only the reccommendations/grpc directory where your Py code and requirements are
COPY reccommendations/grpc /app/reccommendations/grpc

# (Optional) If you want to load your .env for godotenv:
# COPY .env /app/.env

# Create a Python virtual environment at the same path
# that your Go code is referencing in the exec.Command.
# i.e. "reccommendations/grpc/venv/bin/activate"
RUN python3 -m venv /app/reccommendations/grpc/venv \
  && /app/reccommendations/grpc/venv/bin/pip install --no-cache-dir -r /app/reccommendations/grpc/requirements.txt

# Set your environment variables so Go picks them up via `os.Getenv`
ENV MONGO_URI="mongodb+srv://rbb98:cfxARjWMSnojKSjj@cluster0.avlxk.mongodb.net/?retryWrites=true&w=majority"
ENV REDIS_ADDR="redis-17635.c16.us-east-1-3.ec2.redns.redis-cloud.com:17635"
ENV REDIS_PASSWORD="Y3TiIwq5yIk2o7TcnRonae57sWyds6sl"

# Expose port 80
EXPOSE 80

# Finally, run the Go server. The Go code calls the Python server in the background
CMD ["./server"]

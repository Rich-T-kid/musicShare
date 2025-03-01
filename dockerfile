# ------------------------------------------------
# 1) Builder Stage
# ------------------------------------------------
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY vendor ./vendor
COPY . .

RUN go build -mod=vendor -o /server .

# ------------------------------------------------
# 2) Final Stage
# ------------------------------------------------
FROM python:3.10.12

# Install python3-venv to enable virtual environments
RUN apt-get update && apt-get install -y python3-venv

WORKDIR /app

COPY --from=builder /server /app/server
COPY reccommendations/grpc /app/reccommendations/grpc

# Create Python virtual environment & install dependencies
RUN python3 -m venv /app/reccommendations/grpc/venv \
    && /app/reccommendations/grpc/venv/bin/pip install --no-cache-dir -U pip setuptools wheel \
    && /app/reccommendations/grpc/venv/bin/pip install --no-cache-dir -r /app/reccommendations/grpc/requirements.txt

# Define ENV variables directly in Dockerfile (Docker injects them)
ENV MONGO_URI="mongodb+srv://rbb98:cfxARjWMSnojKSjj@cluster0.avlxk.mongodb.net/?retryWrites=true&w=majority"
ENV REDIS_ADDR="redis-17635.c16.us-east-1-3.ec2.redns.redis-cloud.com:17635"
ENV REDIS_PASSWORD="Y3TiIwq5yIk2o7TcnRonae57sWyds6sl"

EXPOSE 8080

CMD ["./server"]

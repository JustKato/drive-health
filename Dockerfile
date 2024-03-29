# === Build Stage ===
FROM debian:bullseye-slim AS builder
ENV IS_DOCKER TRUE

LABEL org.opencontainers.image.source https://github.com/JustKato/drive-health

# Install build dependencies and runtime dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    musl-dev \
    libsqlite3-dev \
    libsqlite3-0 \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Manually install Go 1.21
ENV GOLANG_VERSION 1.21.0
RUN wget https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz -O go.tgz \
    && tar -C /usr/local -xzf go.tgz \
    && rm go.tgz
ENV PATH /usr/local/go/bin:$PATH

# Set the environment variable for Go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
ENV GO111MODULE=on

# Create the directory and set it as the working directory
WORKDIR /app

# Copy the Go files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o drive-health

# Cleanup build dependencies to reduce image size
RUN apt-get purge -y gcc musl-dev libsqlite3-dev wget \
    && apt-get autoremove -y \
    && apt-get clean

# === Final Stage ===
FROM debian:bullseye-slim AS final

# Set the environment variable
ENV IS_DOCKER TRUE

# Create the directory and set it as the working directory
WORKDIR /app

# Copy only the necessary files from the builder stage
COPY --from=builder /app/drive-health .

# Expose the necessary port
EXPOSE 8080

# Volume for external data
VOLUME [ "/data" ]

# Command to run the executable
CMD ["./drive-health"]

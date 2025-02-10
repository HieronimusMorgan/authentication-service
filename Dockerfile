# Use official Golang 1.23.4 image as base
FROM golang:1.23.4-alpine

# Set working directory inside the container
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project including the .env file
COPY . .

# Change working directory to the location of main.go
WORKDIR /app/cmd

# Set environment variables manually (in case .env is missing)

# Debugging: Check if .env file is present
RUN ls -la /app

# Build the Go application
RUN go build -o app .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/cmd/app"]

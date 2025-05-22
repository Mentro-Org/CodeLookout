# Use a minimal Alpine-based Go image
FROM golang:1.24.3-alpine3.21

# Set the working directory
WORKDIR /app

# Copy Go module files first (helps with caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application source code
COPY . .

# Set the working directory to the cmd folder before building
WORKDIR /app/cmd

# Build the application
RUN go build -o /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["/app/main"]

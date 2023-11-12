FROM golang:1.20

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY job-service/go.mod job-service/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY job-service/. .

# Build the Go app
RUN go build -o main main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./main", "service"]

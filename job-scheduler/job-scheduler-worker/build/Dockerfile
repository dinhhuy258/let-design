FROM golang:1.20

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY job-scheduler-worker/go.mod job-scheduler-worker/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY job-scheduler-worker/. .

# Build the Go app
RUN go build -o main main.go

# Command to run the executable
ENTRYPOINT ["./main", "worker"]

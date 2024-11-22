# Start from the official Golang base image
FROM golang:1.23.3

# Copy go mod and sum files
# COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY app /go/src/app

# Set the Current Working Directory inside the container
WORKDIR /go/src/app

# Build the Go app
# RUN go build -o main .
RUN go get github.com/mattn/go-sqlite3
# RUN go run main.go

# Expose port 8080 to the outside world
# EXPOSE 8080

# Command to run the executable
# CMD ["./main"]

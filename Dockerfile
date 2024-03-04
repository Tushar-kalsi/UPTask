# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.17 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o task-reminder .

######## Start a new stage from scratch #######
# Use a small Alpine Linux image to run our application
# Alpine is chosen for its small footprint compared to Ubuntu
FROM alpine:latest  

# Install ca-certificates for HTTPS calls & timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/task-reminder .

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./task-reminder"]


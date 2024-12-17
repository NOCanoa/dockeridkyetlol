# Use the official Go image as the base
FROM golang:alpine

# Set the working directory to /app
WORKDIR /app

# Copy the Go code into the container
COPY . /app

# Set the environment variable for the Go compiler
ENV GOOS=linux
ENV GOARCH=amd64

# Build the Go application
RUN go build -o main main.go

# Expose the port that the application will use
EXPOSE 8080

# Run the command to start the application when the container launches
CMD ["./main"]
# Use the official Go image as the base image
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o dating-be .

# Expose the port the application runs on
EXPOSE 8888

# Command to run the application
CMD ["./dating-be"]

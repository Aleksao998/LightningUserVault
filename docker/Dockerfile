# Use an official Go runtime as the parent image
FROM golang:latest

# Set the working directory in the container
WORKDIR /app

# Copy the local package files to the container's workspace
ADD . /app

# Change directory to where the go.mod is located
WORKDIR /app/docker

# Download the Go modules
RUN go mod download

# Change back to the main app directory
WORKDIR /app

# Build the application
RUN make build

# Expose port 8080 for the app
EXPOSE 8080

# Run the binary on container startup
CMD ["./build/lighting_user_vault", "server"]


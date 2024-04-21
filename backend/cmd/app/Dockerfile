# Use a base image with Golang installed
FROM golang:1.21.0 AS build

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# Copy the source code into the container
COPY ./backend/cmd/app .

# Build the server binary
RUN CGO_ENABLED=1 go build -o backend-server .

# Use a lightweight base image for the final container
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the server binary from the build stage
COPY --from=build /app/backend-server /app/

# Expose the port used by your backend server
EXPOSE 8080

# Command to run the backend server
CMD ["./backend-server"]

# Use a base image with Golang installed
FROM golang:1.21.0 AS build

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY ./calculationServer .

# Build the server binary
RUN CGO_ENABLED=1 go build -o calculation-server ./cmd/server

# Use a lightweight base image for the final container
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the server binary from the build stage
COPY --from=build /app/calculation-server /app/

# Expose the port used by your calculation server
EXPOSE 8050

# Command to run the calculation server
CMD ["./calculation-server"]

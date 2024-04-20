# Use a base image with Golang installed
FROM golang:1.21.0 AS build

# Set the working directory
WORKDIR /src

# Copy the go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# Copy the source code into the container
COPY . .

# Build the server binary
RUN CGO_ENABLED=1 go build -o /bin/calculation-server ./calculationServer/cmd/server

# Use a lightweight base image for the final container
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the server binary from the build stage
COPY --from=build /bin/calculation-server /app/

# Expose the port used by your calculation server
EXPOSE 8050

# Command to run the calculation server
CMD ["/app/calculation-server"]

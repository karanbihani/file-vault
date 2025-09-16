# Stage 1: Build the Go binary using a newer Go version
# UPDATED: We've changed the Go version from 1.22 to 1.24 to match your go.mod file.
FROM golang:1.24-alpine AS builder

# This command ensures our base image has the latest security patches and git.
RUN apk update && apk upgrade && apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy the dependency files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application into a static binary.
# CGO_ENABLED=0 is important for creating a static binary that works in a minimal Alpine image.
RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go


# Stage 2: Create the final, minimal production image
# UPDATED: Pinned to a specific Alpine version for reproducible builds.
FROM alpine:3.19

# This command ensures our final, minimal image also has the latest security patches.
RUN apk update && apk upgrade

WORKDIR /root/

# Copy the compiled binary from the 'builder' stage.
COPY --from=builder /server .

# Expose port 8080 to the outside world.
EXPOSE 8080

# The command to run when the container starts.
CMD ["./server"]
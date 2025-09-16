# Stage 1: The 'builder' stage compiles the Go application.
# It uses a full Go development environment.
FROM golang:1.22-alpine AS builder

# SECURITY: Update and upgrade the base image's packages to patch vulnerabilities.
RUN apk update && apk upgrade && apk add --no-cache git

WORKDIR /app

# Copy dependency files and download them first to leverage Docker's build cache.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# Build the application into a single, static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go


# Stage 2: The 'final' stage creates the minimal production image.
# It starts from a clean slate with a tiny Alpine Linux image.
FROM alpine:latest

# SECURITY: Update and upgrade the base image's packages to patch vulnerabilities.
RUN apk update && apk upgrade

WORKDIR /root/

# Copy only the compiled binary from the 'builder' stage.
# No source code or build tools are included in the final image.
COPY --from=builder /server .

# Expose the port the application will run on.
EXPOSE 8080

# The command to run when the container starts.
CMD ["./server"]
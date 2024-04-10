# Build stage with golang:1.21-alpine AS build
FROM golang:1.21-alpine AS build

# Install build dependencies (e.g., gcc, git, etc.) if necessary
RUN apk add --no-cache git

# Set the working directory outside GOPATH to enable the support for modules.
WORKDIR /src

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app

# Final stage to produce a minimal image
FROM gcr.io/distroless/base

# Copy the binary from the build stage to the final stage
COPY --from=build /bin/app /

# Command to run the binary
CMD ["/app"]

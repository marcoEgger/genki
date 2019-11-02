# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.12
ARG GO_VERSION=1.12

# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates

# Download the grpc health checking client binary
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.3.1/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

WORKDIR /build
COPY . .

# build binary
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor GOOS=linux GOARCH=amd64 go build  -ldflags="-w -s" -o /bin/stringer ./main.go

# runtime container image
FROM alpine

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Load binary from builder
COPY --from=builder /bin/stringer /stringer

COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

# ensure ca-certs
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true


# Ensure we are not running as root
USER nobody:nobody

ENTRYPOINT ["/stringer"]
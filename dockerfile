############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /go/src/cjlapao/security-headers-backend

COPY . .

# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/security-headers-backend

############################
# STEP 2 build a small image
############################
FROM scratch

# Copy our static executable.
COPY --from=builder /go/bin/security-headers-backend /go/bin/security-headers-backend

# Run the project binary.
EXPOSE 5000

ENTRYPOINT ["/go/bin/security-headers-backend"]
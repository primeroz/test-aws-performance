###
FROM k8s.gcr.io/pause:3.1 AS pause

#
FROM golang:1.19-buster as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o downloader

FROM amazon/aws-cli:2.8.5

# Copy the binary to the production image from the builder stage.
# COPY --from=builder /app/server /app/server
COPY --from=pause /pause /pause
COPY --from=builder /app/downloader /downloader


ENTRYPOINT ["/pause"]

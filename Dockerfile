FROM golang:1.16-alpine
VOLUME config downloads

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source files
COPY *.go ./

# Build
RUN go build -o /gogdl-ng

EXPOSE 3200
ENTRYPOINT ["/gogdl-ng"]
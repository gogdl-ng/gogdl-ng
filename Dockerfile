FROM golang:1.16-alpine
VOLUME config downloads

# gcc build base 
RUN apk add build-base

# Build
RUN mkdir -p /build
COPY . /build
WORKDIR /build
RUN go build -o gogdl-ng .

# Create app folder and move binary
WORKDIR /
RUN mkdir -p /app
RUN cp /build/gogdl-ng /app/

# Cleanup
RUN rm -r /build

EXPOSE 3200

CMD ["/app/gogdl-ng"]
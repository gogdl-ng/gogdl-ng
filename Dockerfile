FROM golang:1.16-alpine
VOLUME config downloads

RUN apk add build-base

RUN mkdir -p /build
COPY . /build
WORKDIR /build
RUN go build -o gogdl-ng .

WORKDIR /
RUN mkdir -p /app
RUN cp /build/gogdl-ng /app/

RUN rm -r /build

EXPOSE 3200

CMD ["/app/gogdl-ng"]
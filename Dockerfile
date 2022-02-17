FROM golang:1.16-alpine as build

# gcc build base 
RUN apk add build-base && \ 
    mkdir -p /build

#build
COPY . /build
WORKDIR /build
RUN go build -o gogdl-ng .

FROM golang:1.16-alpine
VOLUME config downloads

# Create app folder and move binary
WORKDIR /
RUN mkdir -p /app

COPY --from=build /build/gogdl-ng ./app/

EXPOSE 3200

CMD ["/app/gogdl-ng"]

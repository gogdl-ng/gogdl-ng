##
## Build
##
FROM golang:1.16-alpine as build-env

COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -o gogdl-ng .

##
## Deploy
##
FROM gcr.io/distroless/static-debian11
VOLUME config downloads

COPY --from=build-env /build/gogdl-ng /

EXPOSE 3200

CMD ["/gogdl-ng"]

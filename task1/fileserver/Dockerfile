FROM golang:1.23 AS build

COPY . /src
RUN cd /src && CGO_ENABLED=0 go build -o /fileserver

FROM alpine:3.20

COPY --from=build /fileserver /fileserver

ENTRYPOINT [ "/fileserver" ]

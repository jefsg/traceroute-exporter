FROM golang:1.15.3 AS build

WORKDIR /src/
COPY . ./
RUN go get
RUN go build -o /bin/traceroute_exporter

FROM busybox:latest

COPY --from=build /bin/traceroute_exporter /bin/traceroute_exporter
EXPOSE 8080
CMD ["/bin/traceroute_exporter"]
FROM golang:1.15.3-alpine AS build

WORKDIR /src/
COPY . /src/ 
RUN go get
RUN go build -o /src/build/traceroute_exporter

FROM scratch
COPY --from=build /src/build/traceroute_exporter /bin/traceroute_exporter
EXPOSE 8080
ENTRYPOINT ["/bin/traceroute_exporter"]
FROM golang:1.15.3-alpine AS build

WORKDIR /src/
COPY . /src/ 
RUN go get
RUN go build -o /bin/traceroute_exporter

EXPOSE 8080
CMD ["/bin/traceroute_exporter"]
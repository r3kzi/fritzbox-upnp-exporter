FROM golang:alpine

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o fritzbox-upnp-exporter .

WORKDIR /dist

RUN cp /build/fritzbox-upnp-exporter .

EXPOSE 9624

ENTRYPOINT ["/dist/fritzbox-upnp-exporter"]
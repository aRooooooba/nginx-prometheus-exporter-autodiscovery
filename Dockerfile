FROM golang:1.18

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o nginx-prometheus-exporter-autodiscovery
EXPOSE 9113
CMD [ "./nginx-prometheus-exporter-autodiscovery" ]

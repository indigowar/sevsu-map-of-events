FROM golang:1.18 as builder
WORKDIR /go/src/github.com/indigowar/sevsu-map-of-events
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go/bin/eventmap ./cmd/eventmap


FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/bin/eventmap .

RUN chmod +x ./eventmap

EXPOSE 8000

ENTRYPOINT ["./eventmap"]
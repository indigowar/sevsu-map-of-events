FROM golang:1.18
WORKDIR /app
COPY . .
RUN go build cmd/eventmap/main.go
EXPOSE 8000
ENTRYPOINT ["./main"]

FROM golang:1.26

COPY . /app

WORKDIR /app

RUN go build subscriptions.go

EXPOSE 8080

CMD ["./subscriptions"]

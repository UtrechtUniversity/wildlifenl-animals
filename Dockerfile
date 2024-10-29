# syntax=docker/dockerfile:1
FROM quay.io/projectquay/golang:1.22

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /app/animals main.go

EXPOSE 8080

CMD ["/app/animals"]
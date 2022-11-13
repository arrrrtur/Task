FROM golang:latest

LABEL mainrainer="Baryspiev Artur"

RUN go build -o main .

RUN mkdir /Task
WORKDIR /Task

RUN go mod download
RUN go build -o Task ./cmd/main/app.go


CMD ["./Task"]
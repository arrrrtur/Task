FROM golang:latest

LABEL mainrainer="oxygen12ium3@gmail.com"

COPY ./ ./
RUN go build -o main .
CMD["./main"]
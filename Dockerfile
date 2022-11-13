FROM golang:alpine
WORKDIR /balance

COPY go.mod ./
COPY go.sum ./

#RUN #apt-get update
#RUN #apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
#RUN chmod +x ./wait-for-postgrs.sh


RUN go mod download
RUN go mod verify

COPY . .
RUN go build -o balance ./cmd/main

CMD ["./balance"]
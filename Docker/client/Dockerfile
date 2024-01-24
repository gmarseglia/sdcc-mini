FROM golang:1.21.6

WORKDIR /usr/src/app

COPY ./go.work ./go.work.sum ./

WORKDIR /usr/src/app/src

COPY ./src .

RUN go build -v -o /usr/local/bin/client ./client/client.go

CMD ["client"]
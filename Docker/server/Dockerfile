FROM golang:1.21.6

RUN apt update -y && apt install iproute2 -y

WORKDIR /usr/src/app

COPY ./go.work ./go.work.sum ./

WORKDIR /usr/src/app/src

COPY ./src .

RUN go build -v -o /usr/local/bin/server ./server/server.go

EXPOSE 55555
EXPOSE 55556

# CMD ["server"]
FROM golang:1.21.6

RUN apt update -y && apt install iproute2 -y

WORKDIR /usr/src/app

COPY ./go.work ./go.work.sum ./

WORKDIR /usr/src/app/src

COPY ./src .

RUN go build -v -o /usr/local/bin/worker ./worker/worker.go

EXPOSE 55557

# CMD ["worker"]
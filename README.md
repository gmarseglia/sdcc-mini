# sdcc-mini
Repo for the mini project for the SDCC exam.

# Usage

## Server

Launch 1 server.
Clients will connect to the front port (default 55555) of the server, and workers will connect to the master port (default 55556) of the server.

```
cd src/server
go build server.go
./server --frontPort xxx --masterPort yyy
```

## Workers

Launch 1+ workers.
Worker will connect to the master address (default 127:0:0:1) and port (default 55556) of the server.
Server will connect to the a runtime chosen port of the worker.

```
cd src/worker
go build worker.go
./worker --masterAddr abcd:efgh:hijk:lmno --masterPort yyy
```

## Client

Launch 1+ clients.
Client will connect to the front port of the server.

```
cd src/client
go build client.go
./client --frontAddr abcd:efgh:hijk:lmno --frontPort xxx
```
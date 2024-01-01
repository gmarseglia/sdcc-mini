# sdcc-mini
Repo for the mini project for the SDCC exam.

# Usage

## Quick intro 

There are 3 main actors:

1. The server, which takes requests from all the the client and forward them to the workers.
2. The clients, which send requests to the server.
3. The workers, which resolve the requests.

The main steps are:

1. The server is launched.
2. The workers notify themselves to the server.
3. The client calls a RPC.
4. The server calls a RPC on a round-robin chosen worker.
5. The worker actually completes the procedure.
6. The reply goes back to the client.

The chosen procedure is a simple procedure which chooses randomly between two options.

## Server

Launch 1 server.
Clients will connect to the front port (default 55555) of the server, and workers will connect to the master port (default 55556) of the server.

Install:

```
cd src/server
go build server.go
```

Run:

```
./server --frontPort xxx --masterPort yyy
```

Runtime:

Send SIGTERM to gracefully stop the server.


## Workers

Launch 1+ workers.
Worker will connect to the master address (default 127:0:0:1) and port (default 55556) of the server.
Server will connect to the a runtime chosen port of the worker.

Install:

```
cd src/worker
go build worker.go
```

Run:

```
./worker --masterAddr abcd:efgh:hijk:lmno --masterPort yyy
```

Runtime:

Send SIGTERM to gracefully stop the worker.


## Client

Launch 1+ clients.
Client will connect to the front port of the server.

Install:

```
cd src/client
go build client.go
```

Run:

```
./client --frontAddr abcd:efgh:hijk:lmno --frontPort xxx
```

# Architecture 

## Server 

The server is composed of two main components:

1. The front: comunicates with the clients.
2. The master: comunicates with the workers.

The master components takes record of which workers notified themselves and their network address.

When a worker is chosen by the server and does not reply, then said worker is removed from the list of active workers.

## Workers 

The workers notify themselves to the server via a gRpc exposed by the server, in particular the master component.

Prior to the notification, the worker has started its own gRpc server which exposes the procedure, this component takes the name of the back. The front component of the server will be the client and the back component of the worker will be the gRpc server.

Periodically the worker will check if the master is alive. If not, then the worker will terminate.

## Client

The client is a simple as it gets and it's complelty obvious of the underlying architecture.

The client only acts as gRpc client toward the front component of the server.

# More 

## Develop branch

The debug branch contains some .sh scripts to run multiple clients and workers.
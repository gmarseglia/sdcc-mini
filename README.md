# sdcc-mini
Repo for the mini project for the SDCC exam.

# Usage

## Quick intro 

There are 3 main actors:

1. The server, which takes requests from all the the client and forward them to the workers.
2. The clients, which send requests to the server.
3. The workers, which resolve the requests.

The main steps are:

1. The server (1) is launched.
2. The workers (1+) notify themselves to the server.
3. The client (1+) calls a RPC.
4. The server calls a RPC on a round-robin chosen worker.
5. The worker actually completes the procedure.
6. The reply goes back to the client.

The chosen procedure chooses randomly between two options but simulates a CPU intensive by scaling linearly the time needed to complete with the number of active requests on the worker.

## Configuration 

All the **fields** indicated in **bold** can be set in multiple ways and the priority is:
1. As a flag when launching, e.g. `-FrontPort`.
2. As an envoironment variable.
3. Default value.



# Server

## Intro

Launch 1 server.
Clients will call RPC to the **front port** (`FrontPort`, default: `55555`) of the server.

Workers will register themselves to the **master port** (`MasterPort`, default: `55556`) of the server.

## Standalone

Install:

```
cd src/server
go build server.go
```

Run:

```
./server -FrontPort 12345 -MasterPort 23456
```

Runtime:

Send SIGTERM to gracefully stop the server.

## Docker 

### Build

Dockerfile is in `Docker/server/Dockerfile`.

On Linux, the image can be build from root directory with the script `./Docker/server/build.sh`.

### Run 

When running a containerized server it's mandatory to:
- Publish a port of the Docker host for the front service.
    - If the port is different from the default, then it must be specified in the envoironment variable `FrontPort`.
- Publish a port of the Docker host for the master service.
    - If the port is different from the default, then it must be specified in the envoironment variable `MasterPort`.

The server application will start automatically at container launch.

On Linux, the container can be launched from root directory with the script `./Docker/server/run.sh`.



# Workers

## Intro

Launch 1+ workers.
The worker will connect to the **master address** (`MasterAddr`, no default value as it's mandatory) and **master port** (`MasterPort` ,default: `55556`) of the server.

The worker will register itself with the pair **host address** (`HostAddr`, default: address of the internet connected interface) and **host port** (`HostPort`, default: `55557`).
This is the address that the server will use to call RPC on the worker, so it should be reachable from the server.
For example, if the worker is launched with Docker, then it should be the address and port of the Docker host.

The worker can also change real port on which it will listen to RPC, with the field **back port** (`BackPort`, default: `55557`).

## Standalone

Install:

```
cd src/worker
go build worker.go
```

Run:

```
./worker -MasterAddr 192.168.1.2 -MasterPort 23456 -HostAddr 192.168.1.3 -HostPort 55557
```

Runtime:

Send SIGTERM to gracefully stop the worker.

## Docker

### Build 

Dockerfile is in `Docker/worker/Dockerfile`.

On Linux, the image can be build from root directory with the script `./Docker/worker/build.sh`.

### Run 

When running a containerized worker it's mandatory to:
- Publish a port of the Docker host for the back service.
    - If the port is different from the default, then it must be specified in the envoironment variable `HostPort`.
- Specify the **master address** in the envoironment variable `MasterAddr`.

The worker application will start automatically at container launch.

On Linux, the container can be launched from root directory with the script `./Docker/worker/run.sh`.



# Client

## Intro

Launch 1+ clients.
Client will call RPC to the pair **front address** (`FrontAddr`, no default value as it's mandatory) and **front port** (`FrontPort`, default: `55555`) of the server.
Client will send **request count** (`RequestCount`, default: `1`) requests in parallel to server.

## Standalone

Install:

```
cd src/client
go build client.go
```

Run:

```
./client -FrontAddr 192.168.1.2 -FrontPort 55555 -RequestCount 2
```

## Docker 

### Build 

Dockerfile is in `Docker/client/Dockerfile`.

On Linux, the image can be build from root directory with the script `./Docker/client/build.sh`.

### Run 

When running a containerized client it's mandatory to:
- Specify the **master address** in the envoironment variable `MasterAddr`.

The client application will start automatically at container launch.

On Linux, the container can be launched from root directory with the script `./Docker/client/run.sh`.



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

The client is complelty obvious of the underlying architecture.

The client only acts as gRpc client toward the front component of the server.

# More 

## Develop branch

The Debug branch contains some .sh scripts to run multiple clients and workers.

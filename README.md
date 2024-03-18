# Mini Kubernetes-like Docker Controller

This project demonstrates a miniature version of Kubernetes' control plane for Docker container management. It's
structured with a simplified API server and controller manager to handle Docker containers, allowing for operations like
container creation, starting, and stopping through HTTP requests.

## Architecture Overview

- **KubeAPIServer**: Mimics the Kubernetes API server. It defines HTTP endpoints to interact with Docker containers,
  handling requests to create, start, and stop containers.
- **KubeControllerManager**: Simulates the Kubernetes controller manager, watching for changes in container states and
  ensuring the desired state is achieved.

## How to Build

Before building the application, ensure that the following prerequisites are met:

- **Go** is installed on your system. For installation instructions, refer to
  the [official Go documentation](https://golang.org/doc/install).
- **Docker** is installed and running. Docker is used to manage containers through the application.
  Visit [Docker's official website](https://docs.docker.com/get-docker/) for installation instructions.
- **etcd** is running as the backend store. The application uses etcd to watch for changes in container states and
  manage configurations. If you don't have etcd installed, you can easily run it as a Docker container using the
  instructions provided below.


## Build Go code
Run `go build` to compile the application. This generates an executable.


### Running etcd Using Docker

To run an etcd instance using Docker, execute the following command in your terminal:

```bash
docker run -d -p 2379:2379 -p 2380:2380 --name etcd quay.io/coreos/etcd:latest /usr/local/bin/etcd --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379
```

Ensure Go and Docker are installed on your system.

1. **Clone the repository** to your local machine.
2. Navigate to the project directory in a terminal.
3. Run `go build` to compile the application. This generates an executable.

## Running the Application

After building the project, start the server by executing the compiled binary. For example:

```bash
./mini-control-plane
```

Create Sample docker container and then start/stop it using it's containerID

```shell
curl -X POST http://localhost:8080/containers/create \
-H "Content-Type: application/json"      \
-d '{"image": "golang:alpine", "cmd": ["sh", "-c", "while true; do sleep 1; done"]}' | jq .
curl -X POST http://localhost:8080/containers/f426836933ad2859bf3e5982a91a0f3391289b23746ac03c2439da218cc2a60d/star
curl -X POST http://localhost:8080/containers/f426836933ad2859bf3e5982a91a0f3391289b23746ac03c2439da218cc2a60d/stop
```


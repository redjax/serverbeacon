# Server Beacon <!-- omit in toc -->

Simple Go app that listens on multiple interfaces for external hosts to check server's online status.

## Table of Contents <!-- omit in toc -->

- [Usage](#usage)
  - [Local Development](#local-development)
- [Docker](#docker)
  - [Build and run locally](#build-and-run-locally)
- [Progress](#progress)

## Usage

```shell
serverbeacon -h
```

Test the connection with:

```shell
curl http[s]://your-ip-or-fqdn:18080/v1/health
```

### Local Development

- Pull packages with `go mod tidy`
- Run `./.scripts/build.sh`
- Run `./.scripts/air-build-run.sh` to start development server with hot reloading
- Test with `curl http://localhost:18080/v1/health`

## Docker

The [Dockerfile](./Dockerfile) for `serverbeacon` builds the binary and runs it in a [distroless layer](https://github.com/GoogleContainerTools/distroless). This significantly reduces the image size, attack surface, and build time.

### Build and run locally

Build the container using the [`./.scripts/docker/build.sh` script](./.scripts/docker/build.sh), or by running:

```shell
docker build -t serverbeacon:latest .
```

Run the container using the [`./.scripts/docker/run.sh` script](./.scripts/docker/run.sh), or by running:

```shell
docker run --rm -d -p 18080:18080 --name serverbeacon serverbeacon:latest
```

## Progress

- [ ] HTTP server/REST API
  - [x] Simple HTTP server
  - [x] `/ping` endpoint (return "pong")
  - [-] `/health` endpoint
    - [x] Return JSON with health status and timestamp
    - [ ] An authenticated `/health` endpoint that returns stats about the underlying host
- [ ] SSH server
  - [ ] Allow SSH connections, immediately terminate
  - [ ] Only allow SSH keys, no user/password auth
- [ ] RPC message
  - [ ] Create a client to 'ping' the server
- [-] Docker container
  - [x] distroless runtime image
  - [ ] Publish to a registry
- [ ] Pipelines
  - [ ] Lint/format
  - [ ] Build and release
    - [ ] Go package
    - [ ] Container registry

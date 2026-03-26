# Server Beacon <!-- omit in toc -->

Simple Go app that listens on multiple interfaces for external hosts to check server's online status.

## Table of Contents <!-- omit in toc -->

- [Usage](#usage)
  - [Local Development](#local-development)
- [Docker](#docker)
  - [Build and run locally](#build-and-run-locally)

## Usage

```shell
serverbeacon -h
```

Test the connection with:

```shell
curl http[s]://your-ip-or-fqdn:18080/v1/health
```

### Local Development

Requirements:

- Go
- Bash or compatible shell

Steps:

- Pull packages with `go mod tidy`
- Run `./.scripts/build.sh`
- Run `./.scripts/air-build-run.sh` to start development server with hot reloading
- Test with `curl http://localhost:18080/v1/health`

## Docker

The release pipeline publishes [Docker containers for each app](https://github.com/redjax?tab=packages&repo_name=serverbeacon).

Available containers:

| Name                                                                                         | Description                                                                                    |
| -------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| [`serverbeacon-api`](https://github.com/redjax/serverbeacon/pkgs/container/serverbeacon-api) | Run the REST API directly.                                                                     |
| [`serverbeacon`](https://github.com/redjax/serverbeacon/pkgs/container/serverbeacon)         | Run the `serverbeacon` CLI app. Start different listeners, i.e. `serverbeacon start rest-api`. |

Run the REST API server:

```shell
docker run --rm -d -p 18080:18080 --name serverbeacon-api serverbeacon-api:latest
```

Run the `serverbeacon` CLI container:

```shell
docker run --rm -d -p 18080:18080 --name serverbeacon-api serverbeacon:latest start rest-api
```

You can also use [the included `compose.yml`](./compose.yml) to run with `docker compose up -d`. Download with:

```shell
curl -o compose.yml https://raw.githubusercontent.com/redjax/serverbeacon/refs/heads/main/compose.yml
```

### Build and run locally

Build the container using the [`./.scripts/docker/build.sh` script](./.scripts/docker/build.sh), or by running:

REST API container:

```shell
docker build --target serverbeacon-api -t serverbeacon-api:latest .
```

`serverbeacon` CLI container:

```shell
docker build --target serverbeacon-api -t serverbeacon:latest .
```

Run the container using the [`./.scripts/docker/run.sh` script](./.scripts/docker/run.sh), or by running:

REST API:

```shell
docker run --rm -d -p 18080:18080 --name serverbeacon-api serverbeacon-apis:latest
```

`serverbeacon` CLI:

```shell
docker run --rm -d -p 18080:18080 --name serverbeacon serverbeacon:latest start rest-api
```

You can also use the [development Docker Compose file](./dev.compose.yml):

```shell
docker compose -f dev.compose.yml up -d
```

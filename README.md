# Golang GRPC Server & Client

## Introduction

This is basic Golang GRPC server & client example. If you want to see cute dog photos in your PC. You can use this repository :)

![Golang GRPC Client-Server](.github/imgs/flow.jpg?raw=true "Flow")

## Requirements

You need to install following packages:
`docker`, `docker-compose`, `make`

## Client

The client sends a request to the server and waits for the response. If the server returns an error, the client prints it.
Otherwise, the client save the image your `/image` path and print the path.

## Server

The server receives a request from the client and sends a request to the dog api. If the api returns an error, the server returns an error to the client. Otherwise, the server saves the image and returns the `[]byte` of the image.

## Running Tests

You can run unit tests with `make unit-tests` command.

```bash
    make unit-tests
```

## Integration Tests
You can run integration tests with `make integration-tests` command. But there is a todo in the code.

```bash
    make integration-tests
```

If you'd like to run integration tests on docker. You can run

```bash
    make integration-tests-docker
```

## TODO

- CERT verification for integration tests

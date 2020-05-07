# Golang microservice template

Exemplary microservice implementation in Golang.
Uses a controller/repository pattern to handle typical CRUD actions.
Exposes the controller methods via HTTP endpoints.


## Dependencies

```bash
go mod download
```

## Mocks

Mocks are generated with (Mockery)[https://github.com/vektra/mockery].
Put the executable under $GOPATH/bin/mockery and create mocks for interfaces with

```bash
%GOPATH%/bin/mockery -all -case=underscore -inpkg
```

## Lint

1. Get golangci-lint from [Github](https://github.com/golangci/golangci-lint).
2. Follow the steps under [Editor Integration](https://github.com/golangci/golangci-lint#editor-integration).
3. Execute with

```bash
golangci-lint run
```

### Running the service

To run the service, follow these steps:

```bash
go run main.go
```

To run the service in a docker container

```bash
docker build -t pizza-service .
```

Once image is built use

```bash
docker run --rm -it pizza-service
```

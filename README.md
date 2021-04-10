# scheduler

[![Test](https://github.com/mv-orchestration/scheduler/actions/workflows/test.yml/badge.svg?branch=develop)](https://github.com/mv-orchestration/scheduler/actions/workflows/test.yml)

### Development

#### Lint
```shell
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.39.0
golangci-lint ./...
```

#### Testing and Coverage
```shell
go test --coverprofile=coverage.out ./...
go tool cover -html=coverage.out 
```

#### Format

```shell
go fmt ./...
```

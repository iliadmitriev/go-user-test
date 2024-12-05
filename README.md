# Simple user db

Created for a demonstartion of a test written in Golang.

## DB Init

```bash
sqlite3 main.db < main.sql
```

## Building

```bash
go generate ./...
go build cmd/main.go
```

## Testing

```bash
go test -coverpkg=./internal/... -coverprofile coverage.out ./...
```

## Cleanup

```bash
rm -rf internal/mocks
```

## References

1. [sqlite3](https://github.com/mattn/go-sqlite3)
2. [mockery](https://github.com/vektra/mockery)
3. [sqlmock](https://github.com/DATA-DOG/go-sqlmock)
4. [testify](https://github.com/stretchr/testify)

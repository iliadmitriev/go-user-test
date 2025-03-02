# Simple user db

Created for a demonstartion of a test written in Golang.

## DB Init

```bash
sqlite3 main.db < main.sql
```

## Building

```bash
go generate ./...
go build -o ./go-user cmd/main.go
```

## Testing

```bash
go test -coverpkg=./internal/... -coverprofile coverage.out ./...
```

## Cleanup

```bash
rm -rf internal/mocks

```

## Check

```bash
xh :8080/user/user5

xh :8080/user/ login=user5 password=secret name=ivan
```

```bash
grpcurl -plaintext localhost:5000 list

echo '{"login":"user5"}' |
   grpcurl -plaintext -d @ localhost:5000 user.v1.UserService/GetByLogin

echo '{"login":"kek","password":"seret","name":"kek"}' |
   grpcurl -plaintext -d @ localhost:5000 user.v1.UserService/Create
```

## References

1. [uber-go/fx](https://github.com/uber-go/fx)
2. [user-go/zap](https://github.com/uber-go/zap)
3. [grpc/grpc-go](https://github.com/grpc/grpc-go)
4. [protobuf-go/types/code](https://github.com/protocolbuffers/protobuf-go)
5. [protobuf-go/types/doc](https://protobuf.dev/reference/protobuf/google.protobuf/)
6. [sqlite3](https://github.com/mattn/go-sqlite3)
7. [mockery](https://github.com/vektra/mockery)
8. [gofakeit](github.com/brianvoe/gofakeit)
9. [sqlmock](https://github.com/DATA-DOG/go-sqlmock)
10. [testify](https://github.com/stretchr/testify)

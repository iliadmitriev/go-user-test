# Simple user db

Created for a demonstartion of a test written in Golang.

db init

```bash
sqlite3 main.db < main.sql
```

building

```bash
go generate ./...
go build cmd/main.go
```

testing

```bash
go test ./...
```

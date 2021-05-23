# Charity backend for SABKAD project

edit .env file to match postgres database info with following info:

```yaml
POSTGRES_URL="host=<host> user=<user> password=<pass> dbname=<dbname> sslmode=disable"
```

run with following command:

```go
go run cmd/server/main.go
```

Or build with:

```go
go build cmd/server/main.go
./main
```

The default configuration runs on port 9091

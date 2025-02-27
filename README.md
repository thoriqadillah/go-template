# App Boilerplate Template
## Requirements
-  Go 
-  [Goose](https://github.com/pressly/goose)
-  Docker (Optional)
-  Postgres (Optional)
-  Redis (Optional)

## Configuration
- Change your app name in the env
- Configure your env. You can create `.env` file with the same key [here](./env/env.go)
- Change your db name in the docker compose

## Build The Project
### With docker (WIP)
```bash
docker compose up --build
```

## Without docker
Make sure you have the following running
- Postgres
- Redis
- Build the app
```bash
go build .
```

### Migrate the database
If the docker volume is new (with docker), or this is the first time you setup the database for this app (with/without docker), then you need to migrate the database first
```bash
goose -dir="./db/migration" postgres "postgresql://postgres@localhost:5432/packform?sslmode=disable" up  
```

## CLI
### Seed the database
```bash
go run ./cmd/main.go seed
```

## Development
### Getting Started
Install the dependencies
```bash
go mod tidy
go run main.go
```

## Logging
To make the logging prettier, run the app with the following
```bash
go run main.go 2>&1 | jq -R 'try fromjson catch .'
```

## Test
```bash
go test ./... -v
```
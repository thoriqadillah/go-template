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
- Install the packages
```bash
go mod tidy
```
- Migrate the database
```bash
goose -dir="./db/migration" postgres "postgresql://postgres@localhost:5432/app?sslmode=disable" up  
```

## Development
Generate the model from database
```bash
cd db 
go run github.com/stephenafamo/bob/gen/bobgen-psql@latest -c ./bobgen.yaml
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

## Build The Project
### With docker (WIP)
```bash
docker compose up
```

## Without docker
Make sure you have the following running
- Postgres
- Redis
- Build the app
```bash
go build .
```

## CLI
You can build the cli or just run it
```bash
go run ./cmd/main.go --help
```
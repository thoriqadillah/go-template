# App
## Requirements
-  Go 
-  [Goose](https://github.com/pressly/goose)
-  [River Queue](https://riverqueue.com/docs#running-migrations)

## Configuration
- Change your app name in the env
- Change your db name in the docker compose
- Migrate the database
```bash
goose postgres "postgresql://postgres@localhost:5432/app?sslmode=disable" up  
```
```bash
river migrate-up --database-url "postgresql://postgres@localhost:5432/app?sslmode=disable"
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
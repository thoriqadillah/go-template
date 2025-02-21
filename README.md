# App
## Requirements
-  Go 
-  [Goose](https://github.com/pressly/goose)
-  [River](https://riverqueue.com/docs#running-migrations)
-  [Bob](https://bob.stephenafamo.com/docs)

## Configuration
- Change your app name in the env
- Change your db name in the docker compose
- Migrate the database
```bash
goose -dir="./db/migration" postgres "postgresql://postgres@localhost:5432/app?sslmode=disable" up  
```
- Generate the model from database
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
# App
## Requirements
-  Go 
-  [Goose](https://github.com/pressly/goose)

## Logging
To make the logging prettier, run the app with the following
```bash
go run main.go 2>&1 | jq -R 'try fromjson catch .'
```


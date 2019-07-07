go test ./... -coverprofile=coverage;  go tool cover -func=coverage; del coverage

#update badge in readme.md with coverage ammount
compile:
	GOOS=windows GOARCH=amd64 go build -o bin/gitswiss-amd64.exe main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/gitswiss-amd64-darwin main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/gitswiss-arm64-darwin main.go
	GOOS=linux GOARCH=amd64 go build -o bin/gitswiss-amd64-linux main.go

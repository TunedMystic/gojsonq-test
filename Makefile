
## @(app) - Run the Go app --watch
run:
	go run main.go


## @(app) - Build the app binary
build:
	go build -ldflags="-s -w" -o app main.go

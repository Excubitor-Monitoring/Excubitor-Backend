install-deps:
	echo "Installing all go dependencies"
	go mod download
build:
	echo "Compiling project for current platform"
	go build -o bin/excubitor-backend cmd/excubitor/main.go
run:
	go run cmd/excubitor/main.go

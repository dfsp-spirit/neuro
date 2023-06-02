
build:
	go build -o bin/neurogo_example example_neurogo.go

run:
	go run example_neurogo.go


compile:
	echo "Compiling example for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/neurogo_example-linux-arm example_neurogo.go
	GOOS=linux GOARCH=arm64 go build -o bin/neurogo_example-linux-arm64 example_neurogo.go
	GOOS=freebsd GOARCH=386 go build -o bin/neurogo_example-freebsd-386 example_neurogo.go

all: build

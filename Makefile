
build:
	go build -o bin/neurogo_example cmd/example1/example_neurogo.go

run:
	go run cmd/example1/example_neurogo.go --meshfile data/lh.white --exportply lhwhite.ply --exportobj lhwhite.obj --exportstl lhwhite.stl


compile:
	echo "Compiling example for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/neurogo_example-linux-arm cmd/example1/example_neurogo.go
	GOOS=linux GOARCH=arm64 go build -o bin/neurogo_example-linux-arm64 cmd/example1/example_neurogo.go
	GOOS=freebsd GOARCH=386 go build -o bin/neurogo_example-freebsd-386 cmd/example1/example_neurogo.go

all: build

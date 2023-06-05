
build:
	go build -o bin/neurogo_example cmd/example1/example_neurogo.go
	go build -o bin/neurogo_example_curv cmd/example_curv/example_curv.go

run:
	go run cmd/example1/example_neurogo.go --meshfile data/lh.white --exportply lhwhite.ply --exportobj lhwhite.obj --exportstl lhwhite.stl

run_surf: run

run_curv:
	go run cmd/example_curv/example_curv.go --curvfile data/lh.thickness --exportjson lhthickness.json

all: build

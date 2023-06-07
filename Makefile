
build:
	go build -o bin/neuro_example_surface cmd/example_surface/example_surface.go
	go build -o bin/neuro_example_curv cmd/example_curv/example_curv.go

run:
	go run cmd/example_surface/example_surface.go --meshfile testdata/lh.white --exportply lhwhite.ply --exportobj lhwhite.obj --exportstl lhwhite.stl

run_surf: run

run_curv:
	go run cmd/example_curv/example_curv.go --curvfile testdata/lh.thickness --exportjson lhthickness.json

all: build

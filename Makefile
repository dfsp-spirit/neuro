
build:
	go build -o bin/neuro_example_surface cmd/example_surface/example_surface.go
	go build -o bin/neuro_example_curv cmd/example_curv/example_curv.go
	go build -o bin/neuro_example_mgh cmd/example_mgh/example_mgh.go

run:
	go run cmd/example_surface/example_surface.go --meshfile testdata/lh.white --exportply lhwhite.ply --exportobj lhwhite.obj --exportstl lhwhite.stl

run_surf: run

run_curv:
	go run cmd/example_curv/example_curv.go --curvfile testdata/lh.thickness --exportjson lhthickness.json

run_mgh:
	go run cmd/example_mgh/example_mgh.go --mghfile testdata/brain.mgh --informat "auto"

run_mgz:
	go run cmd/example_mgh/example_mgh.go --mghfile testdata/brain.mgz --informat "auto"

run_all:
	make run_surf
	make run_curv
	make run_mgh

all: build

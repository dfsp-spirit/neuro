

## `neuro` Development information

For people interested in improving neuro or in trying the [example applications](./cmd/).


### Building and running the brain mesh demo app

To build (but not to run), you will need to have golang installed. The installation is very easy and fast under Linux, MacOS and Windows and explained on the [official Go installation website](https://go.dev/doc/install).

If you have go, clone this repo and change into it:

```shell
git clone https://github.com/dfsp-spirit/neuro
cd neuro
```

Running the demo app as explained below will read the provided demo file `data/lh.white` in FreeSurfer binary surface format and export it to three files in PLY (Stanford), OBJ (Wavefront Object) and STL mesh file format, respectively.


#### Option 1: Building and running manually

Build:

```shell
go build cmd/example_surface/example_surface.go
```

Then run it:

```shell
./example_surface --meshfile data/lh.white --exportply lhwhite.ply --exportobj lhwhite.obj --exportstl lhwhite.stl
```


#### Option 2: Building and running if you have `make`

To build and run, use:

```shell
make run
```

#### Visualizing the exported mesh


If you have a standard mesh viewer like [MeshLab](https://www.meshlab.net/) installed, you can view the exported brain hemisphere mesh:

```shell
meshlab lhwhite.ply
```

You can also try the other file formats (`meshlab lhwhite.obj`, `meshlab lhwhite.stl`) but the meshes look identical.

If you do not have a mesh viewer installed, you can use the web version of MeshLab at [meshlabjs.net](http://www.meshlabjs.net/) directly in your browser.

![Vis](./lhwhite.jpg?raw=true "Visualization of the demo brain mesh.")


### Running the unit tests locally

```shell
go test -v
```

If you want to inspect a detailed HTML coverage report in your browser:

```shell
go test -v -coverprofile cover.out
go tool cover -html=cover.out
```

### Continuous Integration (CI) Results

<!-- badges: start -->
[![Main branch on Github Actions](https://github.com/dfsp-spirit/neuro/actions/workflows/unittests.yml/badge.svg?branch=main)]
<!-- badges: end -->



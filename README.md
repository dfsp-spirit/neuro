# neurogo
Work-in-progress Go module to read structural neuroimaging file formats, ignore.

This repo contains a very early version of a [Go](https://go.dev/) module for reading structural neuroimaging file formats. Currently supported formats include:

* FreeSurfer brain surface format (a triangular mesh file format, used for recon-all output files like `<subject>/lh.white`)
    - Export to PLY, STL, OBJ formats
    - Computation of basic mesh properties (vertex and face count, bounding box, average edge length, total surface area, ...) 


## Usage

Not yet.

The module is far from ready and has not been published yet, so you cannot `go get` it. All information below is developer information, i.e., intended for people who want to try the development version, typically to work on it.


## Building and running the demo app

To build (but not to run), you will need to have golang installed. The installation is very easy and fast under Linux, MacOS and Windows and explained on the [official Go installation website](https://go.dev/doc/install).

If you have go, clone this repo and change into it:

```shell
git clone https://github.com/dfsp-spirit/neurogo
cd neurogo
```

Running the demo app as explained below will read the provided demo file `data/lh.white` in FreeSurfer binary surface format and export it to three files in PLY (Stanford), OBJ (Wavefront Object) and STL mesh file format, respectively.


### Building and running manually

Build:

```shell
go build cmd/example1/example_neurogo.go
```

Then run it:

```shell
./example_neurogo --meshfile data/lh.white --exportply lhwhite.ply --exportobj lhwhite.obj --exportstl lhwhite.stl
```


### Building and running if you have `make`

To build and run, use:

```shell
make run
```

## Visualizing the exported mesh


If you have a standard mesh viewer like [MeshLab](https://www.meshlab.net/) installed, you can view the exported brain hemisphere mesh:

```shell
meshlab lhwhite.ply
meshlab lhwhite.obj
meshlab lhwhite.stl
```

If not, you could use the web version at [meshlabjs.net](http://www.meshlabjs.net/) directly in your browser.

![Vis](./lhwhite.jpg?raw=true "Visualization of the demo brain mesh.")


## Running the unit tests

```shell
go test -v
```


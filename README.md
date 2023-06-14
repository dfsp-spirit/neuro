# neuro
Work-in-progress Go module to read structural neuroimaging file formats, ignore for now.

## About

This repo contains a very early version of a [Go](https://go.dev/) module for reading structural neuroimaging file formats. Currently supported formats include:

* [FreeSurfer](https://freesurfer.net) brain surface format: a triangular mesh file format. Used for recon-all output files like `<subject>/lh.white`
    - Read file format (function `ReadFsSurface`) into `Mesh` data structure.
    - Export `Mesh` to PLY, STL, OBJ formats.
    - Computation of basic `Mesh` properties (vertex and face count, bounding box, average edge length, total surface area, ...).
* FreeSurfer curv format: stores per-vertex data (also known as a brain overlay), e.g., cortical thickness at each vertex of the brain mesh. Typically used for native space data for a single subject, for recon-all output files like `<subject>/lh.thickness`.
    - Read file format (function `ReadFsCurv`)
    - Write file format (function `WriteFsCurv`)
    - Export data to JSON format.
* FreeSurfer MGH and MGZ formats: store 3-dimensional or 4-dimensional (subject/time dimension) magnetic resonance imaging (MRI) scans of the human brain (e.g., `<subject>/mri/brain.mgz). Can also be used to store per-vertex data, including multi-subject data on a common brain template like fsaverage. The MGZ format is just gzip-compressed MGH format.
    - Read MGH format (function `ReadFsMgh`)
    - Read MGZ format (function `ReadFsMgh`), without the need to manually decompress first.
    - Full header information is available, so the image orientation can be reconstructed from the RAS information.


## Usage

Not yet.

The module is far from ready and has not been published yet, so you cannot `go get` it right now.

All information below is developer information, i.e., intended for people who want to try the development version, typically to work on it.


## Development information

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

### Continuous Integration (CI)

<!-- badges: start -->
[![Main branch on Github Actions](https://github.com/dfsp-spirit/neuro/actions/workflows/unittests.yml/badge.svg?branch=main)](https://github.com/dfsp-spirit/neuro/actions/workflows/unittests.yml)
<!-- badges: end -->


## Author, License and Getting Help

The `neuro` module for Go was written by [Tim Schäfer](https://ts.rcmd.org).

It is free software, published under the very permissive [MIT license](./LICENSE).

Note that this library is **not** a part of FreeSurfer, and it is **in no way** endorsed by the FreeSurfer developers. Please do not contact them regarding this library, especially not for support. [Open an issue](https://github.com/dfsp-spirit/neuro/issues) in this repo instead.


### TODO and planned for next releases

[] Add consistent logging
[] Support reading labels (like cortex label in `<subject>/label/lh.cortex.label`)
[] Support reading annots (brain surface parcellations, like Desikan-Killiani in `<subject>/label/lh.aparc.annot`)
[] write support for MGH format

If you need any of these, or something else, urgently, please open an issue. It's no big deal to add it.

### Related packages

* [github.com/okieraised/gonii](https://github.com/okieraised/gonii): Standalone, pure golang NIfTI file parser by Thomas Pham. I have not tried it yet, but it seems to be the most popular Golang NIfTI reader, and the NIfTI format is a more common alternative to MGH/MGZ, used by many neuroimaging software packages. One can, of course, convert between NIfTI and MGH/MGZ on the command line with standard FreeSurfer tools like `mri_convert`. 
* [github.com/dfsp-spirit/libfs](https://github.com/dfsp-spirit/libfs): A portable, header-only, single file, no-dependency, mildly templated, C++11 library for accessing FreeSurfer neuroimaging file formats by Tim Schäfer (me). File format information used in `neuro` comes from `libfs`, and the APIs are similar. Note though that `neuro` is a separate pure Go implementation, not a wrapper around `libfs`. 


# neuro
Go module for reading and writing structural neuroimaging file formats. Supports FreeSurfer MGH, MGZ, and related formats.


<!-- badges: start -->
[![Main branch on Github Actions](https://github.com/dfsp-spirit/neuro/actions/workflows/unittests.yml/badge.svg?branch=main)](https://github.com/dfsp-spirit/neuro/actions/workflows/unittests.yml)
[![GoDoc](https://godoc.org/github.com/dfsp-spirit/neuro?status.svg)](https://godoc.org/github.com/dfsp-spirit/neuro) [![license](https://img.shields.io/github/license/dfsp-spirit/neuro.svg)](https://github.com/dfsp-spirit/neuro/blob/main/LICENSE)
[![DOI](https://zenodo.org/badge/DOI/10.5281/zenodo.8126957.svg)](https://doi.org/10.5281/zenodo.8126957)
<!-- badges: end -->


## About

This repo contains a very early version of a [Go](https://go.dev/) module for reading structural neuroimaging file formats. Currently supported formats include:

* [FreeSurfer](https://freesurfer.net) brain surface format: a triangular mesh file format. Used for recon-all output files like `<subject>/surf/lh.white`.
    - Read file format (function `ReadFsSurface`) into `Mesh` data structure.
    - Export `Mesh` to PLY, STL, OBJ formats.
    - Computation of basic `Mesh` properties (vertex and face count, bounding box, average edge length, total surface area, ...).
* FreeSurfer curv format: stores per-vertex data (also known as a brain overlay), e.g., cortical thickness at each vertex of the brain mesh. Typically used for native space data for a single subject, for recon-all output files like `<subject>/surf/lh.thickness`.
    - Read file format (function `ReadFsCurv`)
    - Write file format (function `WriteFsCurv`)
    - Export data to JSON format.
* FreeSurfer MGH and MGZ formats: store 3-dimensional or 4-dimensional (subject/time dimension) magnetic resonance imaging (MRI) scans of the human brain (e.g., `<subject>/mri/brain.mgz`). Can also be used to store per-vertex data, including multi-subject data on a common brain template like fsaverage (e.g., files like `<subject>/surf/lh.thickness.fwhm5.fsaverage.mgh`). The MGZ format is just gzip-compressed MGH format.
    - Read MGH format (function `ReadFsMgh`)
    - Read MGZ format (function `ReadFsMgh`), without the need to manually decompress first. The function handles both MGH and MGZ.
    - Full header information is available, so the image orientation can be reconstructed from the RAS information.
* FreeSurfer label format: these files store labels, i.e., extra information for a subset of the vertices of a mesh or the voxels of a volume. Sometimes per-vertex or per-voxel data is stored in the labels data field, but in other case the relevant information is simply whether or not a certain element (voxel, vertex) is part of the label. Used for recon-all output files like `<subject>/label/lh.cortex.label`.
    - Read ASCII label format (function `ReadFsLabel`)
    - See also the related utility function `VertexIsPartOfLabel`

![Vis](./lhwhite.jpg?raw=true "Visualization of the demo brain mesh.")

## Usage

### Installation

```shell
go get github.com/dfsp-spirit/neuro
```

### Full Documentation including usage examples for functions

The full documentation can be found on the central go documentation page at [pkg.go.dev](https://pkg.go.dev/github.com/dfsp-spirit/neuro#section-documentation).

It includes the full API documentation and usage examples for the functions.


### Complete demo applications

Demo applications that use `neuro` are available in the [cmd/](./cmd/) directory:

* A command line app that reads a FreeSurfer mesh and prints some mesh information, like total surface area, average edgle length, etc: [example_surface.go](./cmd/example_surface/example_surface.go)
* A command line app that reads per-vertex cortical thickness data from a FreeSurfer curv file and exports it to a JSON file: [example_curv.go](./cmd/example_curv/example_curv.go)
* A command line app that reads a three-dimensional human brain scan (MRI image) from a FreeSurfer MGH file and prints some header data and the value of a voxel: [example_mgh.go](./cmd/example_mgh/example_mgh.go)
* A command line app that reads a label from a FreeSurfer surface label file and optionally exports the label data to JSON format: [example_label.go](./cmd/example_label/example_label.go)


## Developer information

Please see the [Developer information](./README_DEV.md) if you want to compile and run the demo apps, unit tests, and similar things.

## Author, License and Getting Help

The `neuro` module for Go was written by [Tim Schäfer](https://ts.rcmd.org).

It is free software, published under the very permissive [MIT license](./LICENSE).

Note that this library is **not** a part of FreeSurfer, and it is **in no way** endorsed by the FreeSurfer developers. Please do not contact them regarding this library, especially not for support. [Open an issue](https://github.com/dfsp-spirit/neuro/issues) in this repo instead.

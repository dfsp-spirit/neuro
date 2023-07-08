

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

### Contributing

We are happy to accept contributions. To prevent wasted efforts, please get in touch by opening an issue and discussing things first before you start working on larger changes.

We use the typical open source procedure:

* for the repo to your Github account
* checkout from your copy of the repo locally
* change branch to develop, then create a new branch from there
* commit your changes and tests into the local branch
* push to your copy of the repo on Github
* on the Github website, go to the new branch and create a pull request against the develop branch of my repo.
* wait for the tests, and I will comment on the PR.
* if everything is green, I will merge your changes into develop. before the next release, I will merge develop into main, and your changes will be available for everyone.



### For maintainers: publishing a new package version

The process is explained so well and short in the [official Go documentation on publishing a module](https://go.dev/doc/modules/publishing) that I couldn't condense it any further, just read it there.

### Building under MS Windows

This is not officially supported and we cannot help with any issues you encounter, but the installation of Go under Windows is straight-forward. There is an official installer on the golang website, just run it and you should be fine. If your path is setup correctly, you should be able to build in your clone of the repo:

```shell
# in your local copy of the repo
go build cmd\example_surface\example_surface.go
```

The installer mentioned above will not install `make` though. While `make` is optional, it is definitely convenient to have it. If you are building stuff under Windows, chances are you already have conda installed. If so, it is easy to install `make` from MinGW using `conda`:

```shell
# activate your conda environment of choice first. Then:
conda install -c conda-forge m2w64-make
```

Once you have it, you can build and run neuro demo apps using `make` as illustrated for Unix above, just remember that the binary that comes in the `m2w64-make` package is called `mingw32-make` and replace it in the commands:

```shell
mingw32-make
```

or

```shell
mingw32-make run_mgh
```
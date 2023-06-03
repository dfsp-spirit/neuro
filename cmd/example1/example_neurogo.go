// Demo application for the neurogo package. Reads a FreeSurfer mesh and prints basic mesh properties.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dfsp-spirit/neurogo"
)

var (
    meshfile string
    verbosity * int

	err error
)

func init() {
    flag.StringVar(&meshfile, "meshfile", "lh.white", "The mesh file to read, in FreeSurfer curv format.")
    verbosity = flag.Int("verbosity", 1, "Verbosity level: 0 = silent, 1 = info, 2 = debug.")
}

func main() {

	apptag := "[EX1] "
	neurogo.Verbosity = *verbosity

    flag.Parse()
	fmt.Println("=====[ Neuro Example 1: Read a FreeSurfer mesh file ]=====")

	if *verbosity > 0 {
    	fmt.Println(apptag, "meshfile:", meshfile)
    	fmt.Println(apptag, "verbosity:", *verbosity)
	}

	if _, err := os.Stat(meshfile); err != nil {
		fmt.Printf("%sCould not stat file '%s': '%s', exiting.\n.", apptag, meshfile, err)
		return
	}

	mesh, err := neurogo.ReadFreesurferMesh(meshfile)
	if err != nil {
		fmt.Printf("%sFailed to read mesh from file '%s': '%s'.\n", apptag, meshfile, err)
		return
	}

	if *verbosity > 0 {
    	fmt.Printf("%sRead mesh with %d vertices and %d faces from meshfile '%s'.\n", apptag, len(mesh.Vertices), len(mesh.Faces), meshfile)
	}
}
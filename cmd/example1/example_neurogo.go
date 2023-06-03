// Demo application for the neurogo package. Reads a FreeSurfer mesh and prints basic mesh properties.

package main

import (
	"flag"
	"fmt"

	"github.com/dfsp-spirit/neurogo"
)

var (
    meshfile  *string
    verbosity *int

	err error
)

func init() {
    flag.StringVar(meshfile, "meshfile", "lh.white", "The mesh file to read, in FreeSurfer curv format.")
    verbosity = flag.Int("verbosity", 1, "Verbosity level, from 0 = silent to 3 = debug.")
}

func main() {

    flag.Parse()

	if *verbosity > 0 {
    	fmt.Println("meshfile:", *meshfile)
    	fmt.Println("verbosity:", *verbosity)
	}

	err = neurogo.ReadFreesurferMesh(*meshfile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *verbosity > 0 {
    	fmt.Println("Mesh read from meshfile:", *meshfile)
	}
}
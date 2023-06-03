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
    verbosity = flag.Int("verbosity", 1, "Verbosity level, from 0 = silent to 3 = debug.")
}

func main() {

    flag.Parse()

	if *verbosity > 0 {
    	fmt.Println("meshfile:", meshfile)
    	fmt.Println("verbosity:", *verbosity)
	}

	if _, err := os.Stat(meshfile); err != nil {
		fmt.Printf("Could not stat file '%s', file does not exist or is not readable, exiting.\n.", meshfile)
		return
	}

	err = neurogo.ReadFreesurferMesh(meshfile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *verbosity > 0 {
    	fmt.Println("Mesh read from meshfile:", meshfile)
	}
}
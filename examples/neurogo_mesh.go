// Demo application for the neurogo package. Reads a FreeSurfer mesh and prints basic mesh properties.

package main

import (
	"flag"
	"fmt"
)

var (
    meshfile  *string
    verbosity *int
)

func init() {
    meshfile = flag.String("meshfile", "lh.white", "The mesh file to read, in FreeSurfer curv format.")
    verbosity = flag.Int("verbosity", 0, "Verbosity level, from 0 = silent to 3 = debug.")
}

func main() {

    flag.Parse()

	if *verbosity > 0 {
    	fmt.Println("meshfile:", *meshfile)
    	fmt.Println("verbosity:", *verbosity)
	}
}
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
	exportfile_ply string
	exportfile_obj string
	exportfile_stl string
    verbosity * int

	err error
)

func init() {
    flag.StringVar(&meshfile, "meshfile", "lh.white", "The input mesh file to read, in binary FreeSurfer surface format.")
	flag.StringVar(&exportfile_ply, "exportply", "", "Output file to which to export mesh in PLY format. Path to it must exist.")
	flag.StringVar(&exportfile_obj, "exportobj", "", "Output file to which to export mesh in OBJ format. Path to it must exist.")
	flag.StringVar(&exportfile_stl, "exportstl", "", "Output file to which to export mesh in STL format. Path to it must exist.")
    verbosity = flag.Int("verbosity", 2, "Verbosity level: 0 = silent, 1 = info, 2 = debug.")
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
    	fmt.Printf("%sRead mesh with %d vertices and %d faces from meshfile '%s'.\n", apptag, len(mesh.Vertices)/3, len(mesh.Faces)/3, meshfile)
	}

	stats, err := neurogo.MeshStats(mesh)
	if err != nil {
		fmt.Printf("%sFailed to compute mesh statistics: '%s'.\n", apptag, err)
		return
	}

	if *verbosity > 0 {
    	fmt.Printf("%sMesh is in area x: %f to %f, y: %f to %f, z: %f to %f.\n", apptag, stats["min_x"], stats["max_x"], stats["min_y"], stats["max_y"], stats["min_z"], stats["max_z"])
		fmt.Printf("%sMesh has %d vertices with center of mass %f, %f, %f.\n", apptag, int(stats["numVertices"]), stats["mean_x"], stats["mean_y"], stats["mean_z"])
		fmt.Printf("%sMesh has %d edges with average length %f.\n", apptag, int(stats["num_edges"]), stats["avg_edge_length"])
		fmt.Printf("%sMesh has %d faces with average area %f, total mesh area is %f.\n", apptag, int(stats["numFaces"]), stats["avg_face_area"], stats["total_area"])
	}

	if len(exportfile_obj) > 0 {
		
	}


}
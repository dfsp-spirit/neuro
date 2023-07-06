// Demo application for the neurogo package. Reads a FreeSurfer mesh and prints basic mesh properties.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dfsp-spirit/neuro"
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
	//neurogo.Verbosity = *verbosity

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

	mesh, err := neuro.ReadFsSurface(meshfile)
	if err != nil {
		fmt.Printf("%sFailed to read mesh from file '%s': '%s'.\n", apptag, meshfile, err)
		return
	}

	if *verbosity > 0 {
    	fmt.Printf("%sRead mesh with %d vertices and %d faces from meshfile '%s'.\n", apptag, len(mesh.Vertices)/3, len(mesh.Faces)/3, meshfile)
	}

	stats, err := neuro.MeshStats(mesh)
	if err != nil {
		fmt.Printf("%sFailed to compute mesh statistics: '%s'.\n", apptag, err)
		return
	}

	if *verbosity > 0 {
    	fmt.Printf("%sMesh is in area x: %f to %f, y: %f to %f, z: %f to %f.\n", apptag, stats["minX"], stats["maxX"], stats["minY"], stats["maxY"], stats["minZ"], stats["maxZ"])
		fmt.Printf("%sMesh has %d vertices with center of mass %f, %f, %f.\n", apptag, int(stats["numVertices"]), stats["meanX"], stats["meanY"], stats["meanZ"])
		fmt.Printf("%sMesh has %d edges with average length %f.\n", apptag, int(stats["numEdges"]), stats["avgEdgeLength"])
		fmt.Printf("%sMesh has %d faces with average area %f, total mesh area is %f.\n", apptag, int(stats["numFaces"]), stats["avgFaceArea"], stats["totalArea"])
	}

	if len(exportfile_obj) > 0 {
		fmt.Printf("%sExporting mesh in format %s to file '%s'.\n", apptag, "obj", exportfile_obj)
		_, err = neuro.Export(mesh, exportfile_obj, "obj")
		if err != nil {
			fmt.Printf("%sError exporting mesh in OBJ: %s", apptag, err)
		} else {
			fmt.Printf("%sExported mesh to file '%s' in OBJ format.\n", apptag, exportfile_obj)
		}
	}
	if len(exportfile_ply) > 0 {
		fmt.Printf("%sExporting mesh in format %s to file '%s'.\n", apptag, "ply", exportfile_ply)
		_, err = neuro.Export(mesh, exportfile_ply, "ply")
		if err != nil {
			fmt.Printf("%sError exporting mesh in PLY format: %s", apptag, err)
		} else {
			fmt.Printf("%sExported mesh to file '%s' in PLY format.\n", apptag, exportfile_ply)
		}
	}
	if len(exportfile_stl) > 0 {
		fmt.Printf("%sExporting mesh in format %s to file '%s'.\n", apptag, "stl", exportfile_stl)
		_, err = neuro.Export(mesh, exportfile_stl, "stl")
		if err != nil {
			fmt.Printf("%sError exporting mesh in STL format: %s", apptag, err)
		} else {
			fmt.Printf("%sExported mesh to file '%s' in STL format.\n", apptag, exportfile_stl)
		}
	}



}
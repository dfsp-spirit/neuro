package neuro

import (
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
)

func almostEqualF64(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func almostEqualF32(a float32, b float32, tolerance float64) bool {
	return math.Abs(float64(a-b)) <= tolerance
}

func TestNumVertices(t *testing.T) {

	var mymesh Mesh
	mymesh.Vertices = make([]float32, 5*3)

	got := NumVertices(mymesh)
	want := 5

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}

func TestNumVerticesEmpty(t *testing.T) {

	var mymesh Mesh

	got := NumVertices(mymesh)
	want := 0

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}

func TestNumFaces(t *testing.T) {

	var mymesh Mesh
	mymesh.Faces = make([]int32, 5*3)

	got := NumFaces(mymesh)
	want := 5

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}

func TestNumFacesEmpty(t *testing.T) {

	var mymesh Mesh

	got := NumFaces(mymesh)
	want := 0

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}

func TestCubeFaces(t *testing.T) {
	var mycube Mesh = GenerateCube()

	got := NumFaces(mycube)
	want := 12

	if got != want {
		t.Errorf("got NumFaces %d, wanted %d", got, want)
	}
}

func TestCubeVertices(t *testing.T) {
	var mycube Mesh = GenerateCube()

	gotNumVertices := NumVertices(mycube)
	wantNumVertices := 8

	if gotNumVertices != wantNumVertices {
		t.Errorf("got NumVertices %d, wanted %d", gotNumVertices, wantNumVertices)
	}
}

func ExampleGenerateCube() {

	var mycube Mesh = GenerateCube()
	fmt.Printf("Cube mesh has %d vertices and %d faces.\n", NumVertices(mycube), NumFaces(mycube))
	// Output: Cube mesh has 8 vertices and 12 faces.
}

func TestMeshStats(t *testing.T) {
	var mycube Mesh = GenerateCube()

	stats, err := MeshStats(mycube)
	if err != nil {
		t.Errorf("got error %s when computing MeshStats", err)
	}

	var wantNumVertices int = 8
	var wantNumFaces int = 12
	var wantNumEdges int = 36
	var wantAvgEdgeLength float32 = 2.276143
	var wantAvgFaceArea float32 = 2.000000
	var wantTotalArea float32 = 24.000002

	gotNumVertices := int(stats["numVertices"])
	gotNumFaces := int(stats["numFaces"])
	gotNumEdges := int(stats["numEdges"])
	gotAvgEdgeLength := stats["avgEdgeLength"]
	gotAvgFaceArea := stats["avgFaceArea"]
	gotTotalArea := stats["totalArea"]

	if gotNumVertices != wantNumVertices {
		t.Errorf("got NumVertices=%d, wanted %d", gotNumVertices, wantNumVertices)
	}

	if gotNumFaces != wantNumFaces {
		t.Errorf("got NumFaces=%d, wanted %d", gotNumFaces, wantNumFaces)
	}

	if gotNumEdges != wantNumEdges {
		t.Errorf("got NumEdges=%d, wanted %d", gotNumEdges, wantNumEdges)
	}

	if !almostEqualF32(gotAvgEdgeLength, wantAvgEdgeLength, 1e-6) {
		t.Errorf("got AvgEdgeLength=%.18f, wanted %.18f", gotAvgEdgeLength, wantAvgEdgeLength)
	}

	if !almostEqualF32(gotAvgFaceArea, wantAvgFaceArea, 1e-6) {
		t.Errorf("got AvgFaceArea=%.18f, wanted %.18f", gotAvgFaceArea, wantAvgFaceArea)
	}

	if !almostEqualF32(gotTotalArea, wantTotalArea, 1e-6) {
		t.Errorf("got TotalArea=%.18f, wanted %.18f", gotTotalArea, wantTotalArea)
	}
}

func ExampleMesh() {
	var mycube Mesh = GenerateCube()
	nv := NumVertices(mycube)
	nf := NumFaces(mycube)
	fmt.Printf("Cube mesh has %d vertices and %d faces.\n", nv, nf)
	// Output: Cube mesh has 8 vertices and 12 faces.
}

func ExampleNumVertices() {
	var mycube Mesh = GenerateCube()
	nv := NumVertices(mycube)
	nf := NumFaces(mycube)
	fmt.Printf("Cube mesh has %d vertices and %d faces.\n", nv, nf)
	// Output: Cube mesh has 8 vertices and 12 faces.
}

func ExampleNumFaces() {
	var mycube Mesh = GenerateCube()
	nv := NumVertices(mycube)
	nf := NumFaces(mycube)
	fmt.Printf("Cube mesh has %d vertices and %d faces.\n", nv, nf)
	// Output: Cube mesh has 8 vertices and 12 faces.
}

func ExampleToPlyFormat() {
	var mycube Mesh = GenerateCube()
	ply_str, err := ToPlyFormat(mycube)
	if err != nil {
		fmt.Printf("Error getting PLY representation: %s\n", err)
	}
	fmt.Printf("PLY format string has %d lines.\n", strings.Count(ply_str, "\n"))
	// Output: PLY format string has 30 lines.
}

func ExampleToObjFormat() {
	var mycube Mesh = GenerateCube()
	obj_str, err := ToObjFormat(mycube)
	if err != nil {
		fmt.Printf("Error getting OBJ representation: %s\n", err)
	}
	fmt.Printf("OBJ format string has %d lines.\n", strings.Count(obj_str, "\n"))
	// Output: OBJ format string has 21 lines.
}

func ExampleToStlFormat() {
	var mycube Mesh = GenerateCube()
	stl_str, err := ToStlFormat(mycube)
	if err != nil {
		fmt.Printf("Error getting STL representation: %s\n", err)
	}
	fmt.Printf("STL format string has %d lines.\n", strings.Count(stl_str, "\n"))
	// Output: STL format string has 86 lines.
}

func ExampleMesh_fromData() {
	mesh := Mesh{}
	mesh.Vertices = []float32{0.0, 1.0, 2.0, 3.0, 4.0, 5.0} // 2 vertices, 3 dimensions each
	mesh.Faces = []int32{0, 1, 2, 3, 4, 5}                  // 2 faces, 3 vertices each
	nv := NumVertices(mesh)
	nf := NumFaces(mesh)
	fmt.Printf("Mesh has %d vertices and %d faces.\n", nv, nf)
	// Output: Mesh has 2 vertices and 2 faces.
}

func ExampleMesh_fromSurfaceFile() {
	var surfFile string = "testdata/lh.white"
	surf, _ := ReadFsSurface(surfFile)

	nv := NumVertices(surf)
	nf := NumFaces(surf)
	fmt.Printf("Surface has %d vertices and %d faces.\n", nv, nf)
	// Output: Surface has 149244 vertices and 298484 faces.
}

func ExampleMeshStats_fromSurfaceFileVerts() {
	var surfFile string = "testdata/lh.white"
	surf, _ := ReadFsSurface(surfFile)

	stats, _ := MeshStats(surf)
	fmt.Printf("Surface has %d vertices and %d faces.\n", int(stats["numVertices"]), int(stats["numFaces"]))
	// Output: Surface has 149244 vertices and 298484 faces.
}

func TestToPlyFormat(t *testing.T) {

	var myCube Mesh = GenerateCube()
	repr_ply, _ := ToPlyFormat(myCube)

	got := strings.Count(repr_ply, "\n")
	want := 30

	if got != want {
		t.Errorf("got %d PLY lines, wanted %d", got, want)
	}
}

func TestToStlFormat(t *testing.T) {

	var myCube Mesh = GenerateCube()
	repr_stl, _ := ToStlFormat(myCube)

	got := strings.Count(repr_stl, "\n")
	want := 86

	if got != want {
		t.Errorf("got %d STL lines, wanted %d", got, want)
	}
}

func TestToObjFormat(t *testing.T) {

	var mycube Mesh = GenerateCube()
	repr_obj, _ := ToObjFormat(mycube)

	got := strings.Count(repr_obj, "\n")
	want := 21

	if got != want {
		t.Errorf("got %d OBJ lines, wanted %d", got, want)
	}
}

func TestExport(t *testing.T) {
	var myCube Mesh = GenerateCube()
	// get a temp file.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	mesh_out_filename := file.Name()
	defer os.Remove(mesh_out_filename) // clean up
	file.Close()

	// export to tmp file.
	Export(myCube, mesh_out_filename, "ply")
	Export(myCube, mesh_out_filename, "stl")
	Export(myCube, mesh_out_filename, "obj")

	//Export(myCube, "cube.ply", "ply")
}

func TestExport_Sphere(t *testing.T) {
	var mySphere Mesh = GenerateSphere(1.0, 20, 20)
	// get a temp file.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	mesh_out_filename := file.Name()
	defer os.Remove(mesh_out_filename) // clean up
	file.Close()

	// export to tmp file.
	Export(mySphere, mesh_out_filename, "ply")
	Export(mySphere, mesh_out_filename, "stl")
	Export(mySphere, mesh_out_filename, "obj")

	//Export(mySphere, "sphere.ply", "ply")
}

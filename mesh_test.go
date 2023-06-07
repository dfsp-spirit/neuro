package neuro

import (
	"math"
	"testing"
)


func almostEqualF64(a, b, tolerance float64) bool {
    return math.Abs(a - b) <= tolerance
}

func almostEqualF32(a float32, b float32, tolerance float64) bool {
    return math.Abs(float64(a - b)) <= tolerance
}

func TestNumVertices(t *testing.T){

	var mymesh Mesh
	mymesh.Vertices = make([]float32, 5 * 3)

    got := NumVertices(mymesh)
    want := 5

    if got != want {
        t.Errorf("got %d, wanted %d", got, want)
    }
}

func TestNumVerticesEmpty(t *testing.T){

	var mymesh Mesh

    got := NumVertices(mymesh)
    want := 0

    if got != want {
        t.Errorf("got %d, wanted %d", got, want)
    }
}

func TestNumFaces(t *testing.T){

	var mymesh Mesh
	mymesh.Faces = make([]int32, 5 * 3)

    got := NumFaces(mymesh)
    want := 5

    if got != want {
        t.Errorf("got %d, wanted %d", got, want)
    }
}

func TestNumFacesEmpty(t *testing.T){

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

    if gotNumVertices !=  wantNumVertices{
        t.Errorf("got NumVertices=%d, wanted %d", gotNumVertices, wantNumVertices)
    }

    if gotNumFaces !=  wantNumFaces{
        t.Errorf("got NumFaces=%d, wanted %d", gotNumFaces, wantNumFaces)
    }

    if gotNumEdges !=  wantNumEdges{
        t.Errorf("got NumEdges=%d, wanted %d", gotNumEdges, wantNumEdges)
    }

    if ! almostEqualF32(gotAvgEdgeLength, wantAvgEdgeLength, 1e-6) {
        t.Errorf("got AvgEdgeLength=%.18f, wanted %.18f", gotAvgEdgeLength, wantAvgEdgeLength)
    }

    if ! almostEqualF32(gotAvgFaceArea, wantAvgFaceArea, 1e-6) {
        t.Errorf("got AvgFaceArea=%.18f, wanted %.18f", gotAvgFaceArea, wantAvgFaceArea)
    }

    if ! almostEqualF32(gotTotalArea, wantTotalArea, 1e-6) {
        t.Errorf("got TotalArea=%.18f, wanted %.18f", gotTotalArea, wantTotalArea)
    }
}
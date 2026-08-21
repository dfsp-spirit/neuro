package neuro

// https://pkg.go.dev/testing

import (
	"fmt"
	"testing"
)

func TestReadFsAnnot(t *testing.T) {
	var annotFile string = "testdata/lh.aparc.annot"

	annot, err := ReadFsAnnot(annotFile)
	if err != nil {
		t.Errorf("ReadFsAnnot failed: %v", err)
	}

	numVertices, err := annot.NumVertices()
	if err != nil {
		t.Errorf("Annot.NumVertices failed: %v", err)
	}
	if numVertices != 149244 {
		t.Errorf("got %d vertices in annot file, wanted %d", numVertices, 149244)
	}

	if annot.Colortable.NumEntries() != 36 {
		t.Errorf("got %d colortable entries in annot file, wanted %d", annot.Colortable.NumEntries(), 36)
	}

	// The Desikan-Killiany atlas (aparc) contains a set of well-known regions.
	expectedRegions := []string{"unknown", "bankssts", "caudalanteriorcingulate", "lateraloccipital", "superiorfrontal", "insula"}
	for _, regionName := range expectedRegions {
		if annot.Colortable.GetRegionIdxByName(regionName) < 0 {
			t.Errorf("expected region '%s' in colortable of annot file, but it was not found.", regionName)
		}
	}

	// The first vertex is part of the lateral occipital region.
	if annot.VertexLabels[0] != 9182740 {
		t.Errorf("got label %d for vertex 0, wanted %d", annot.VertexLabels[0], 9182740)
	}
}

func TestAnnotVertexColors(t *testing.T) {
	var annotFile string = "testdata/lh.aparc.annot"

	annot, err := ReadFsAnnot(annotFile)
	if err != nil {
		t.Errorf("ReadFsAnnot failed: %v", err)
	}

	// The first vertex is part of the lateral occipital region with RGB (20, 30, 140).
	colors := annot.VertexColors(false)
	if len(colors) != 149244*3 {
		t.Errorf("got %d color values without alpha, wanted %d", len(colors), 149244*3)
	}
	expectedRgb := []uint8{20, 30, 140}
	for i := 0; i < 3; i++ {
		if colors[i] != expectedRgb[i] {
			t.Errorf("got color channel %d for vertex 0, wanted %d", colors[i], expectedRgb[i])
		}
	}

	// With alpha, there are 4 values per vertex.
	colorsAlpha := annot.VertexColors(true)
	if len(colorsAlpha) != 149244*4 {
		t.Errorf("got %d color values with alpha, wanted %d", len(colorsAlpha), 149244*4)
	}
	expectedRgba := []uint8{20, 30, 140, 0}
	for i := 0; i < 4; i++ {
		if colorsAlpha[i] != expectedRgba[i] {
			t.Errorf("got RGBA channel %d for vertex 0, wanted %d", colorsAlpha[i], expectedRgba[i])
		}
	}
}

func TestAnnotRegionVertices(t *testing.T) {
	var annotFile string = "testdata/lh.aparc.annot"

	annot, err := ReadFsAnnot(annotFile)
	if err != nil {
		t.Errorf("ReadFsAnnot failed: %v", err)
	}

	// The superior frontal gyrus is a large region in the Desikan-Killiany atlas.
	superiorFrontalVerts := annot.RegionVerticesByName("superiorfrontal")
	if len(superiorFrontalVerts) != 12569 {
		t.Errorf("got %d vertices in region 'superiorfrontal', wanted %d", len(superiorFrontalVerts), 12569)
	}

	// A region name that does not exist must yield an empty vertex list.
	noSuchRegionVerts := annot.RegionVerticesByName("doesnotexist")
	if len(noSuchRegionVerts) != 0 {
		t.Errorf("got %d vertices for non-existing region, wanted 0", len(noSuchRegionVerts))
	}

	// Region vertices by label must match region vertices by name.
	superiorFrontalLabel := annot.Colortable.Label[annot.Colortable.GetRegionIdxByName("superiorfrontal")]
	superiorFrontalVertsByLabel := annot.RegionVerticesByLabel(superiorFrontalLabel)
	if len(superiorFrontalVertsByLabel) != len(superiorFrontalVerts) {
		t.Errorf("got %d vertices by label, wanted %d by name", len(superiorFrontalVertsByLabel), len(superiorFrontalVerts))
	}
}

func TestAnnotVertexRegionNames(t *testing.T) {
	var annotFile string = "testdata/lh.aparc.annot"

	annot, err := ReadFsAnnot(annotFile)
	if err != nil {
		t.Errorf("ReadFsAnnot failed: %v", err)
	}

	regionNames := annot.VertexRegionNames()
	if len(regionNames) != 149244 {
		t.Errorf("got %d region names, wanted %d", len(regionNames), 149244)
	}

	// The first vertex is part of the lateral occipital region.
	if regionNames[0] != "lateraloccipital" {
		t.Errorf("got region name '%s' for vertex 0, wanted 'lateraloccipital'", regionNames[0])
	}
}

func TestAnnotColortableLookups(t *testing.T) {
	var annotFile string = "testdata/lh.aparc.annot"

	annot, err := ReadFsAnnot(annotFile)
	if err != nil {
		t.Errorf("ReadFsAnnot failed: %v", err)
	}

	// Look up the region index by name.
	idx := annot.Colortable.GetRegionIdxByName("insula")
	if idx < 0 {
		t.Errorf("region 'insula' not found in colortable.")
	}
	if annot.Colortable.Name[idx] != "insula" {
		t.Errorf("got name '%s' at region index %d, wanted 'insula'", annot.Colortable.Name[idx], idx)
	}
	// The insula has RGB (255, 192, 32).
	if annot.Colortable.R[idx] != 255 || annot.Colortable.G[idx] != 192 || annot.Colortable.B[idx] != 32 {
		t.Errorf("got insula color (%d,%d,%d), wanted (255,192,32)", annot.Colortable.R[idx], annot.Colortable.G[idx], annot.Colortable.B[idx])
	}

	// Look up the region index by the label computed from the RGBA values.
	label := annot.Colortable.Label[idx]
	idxByLabel := annot.Colortable.GetRegionIdxByLabel(label)
	if idxByLabel != idx {
		t.Errorf("got region index %d by label, wanted %d", idxByLabel, idx)
	}

	// Unknown names and labels must yield -1.
	if annot.Colortable.GetRegionIdxByName("doesnotexist") != -1 {
		t.Errorf("expected -1 for non-existing region name, but got a valid index.")
	}
	if annot.Colortable.GetRegionIdxByLabel(int32(-12345)) != -1 {
		t.Errorf("expected -1 for non-existing region label, but got a valid index.")
	}
}

func ExampleReadFsAnnot() {
	var annotFile string = "testdata/lh.aparc.annot"

	// Read the annotation file
	annot, _ := ReadFsAnnot(annotFile)

	numVertices, _ := annot.NumVertices()
	fmt.Printf("Read annot with %d vertices and %d regions from annot file '%s'.\n",
		numVertices, annot.Colortable.NumEntries(), annotFile)
	// Output: Read annot with 149244 vertices and 36 regions from annot file 'testdata/lh.aparc.annot'.
}

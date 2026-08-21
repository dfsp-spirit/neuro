package neuro

// Related software: libfs for C++, see:
// https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L3206 for the annot
// data structures and https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L4343
// for the annot file format.

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// Maximum number of entries (regions) allowed in an annotation colortable.
const maxColortableEntries = 10000

// Maximum length (in bytes) allowed for a region name in an annotation colortable.
const maxColortableRegionNameLength = 256

// Colortable models the color table from an annotation file, typically used for
// parcellations and integer labels. Each index (in all fields) describes a brain region.
type Colortable struct {
	ID    []int32  // The internal region index (in practice, these go from 0 to N-1).
	Name  []string // The region name, e.g., 'superiorfrontal' for the Desikan-Killiany atlas.
	R     []int32  // The red channel of the RGBA color, range 0-255.
	G     []int32  // The green channel of the RGBA color, range 0-255.
	B     []int32  // The blue channel of the RGBA color, range 0-255.
	A     []int32  // The alpha channel of the RGBA color, range 0-255.
	Label []int32  // The label integer computed from the RGBA values (r + g*256 + b*65536 + a*16777216). Maps to the Annot.VertexLabels field.
}

// NumEntries returns the number of entries (regions) in the Colortable.
//
// Returns:
//   - int: the number of entries in the Colortable.
func (ct Colortable) NumEntries() int {
	return len(ct.ID)
}

// GetRegionIdxByName returns the index of a region in the Colortable by region name.
// Returns -1 if the region is not found.
//
// Parameters:
//   - queryName: the region name to look up, e.g., 'superiorfrontal'.
//
// Returns:
//   - int: the index of the region, or -1 if it was not found.
func (ct Colortable) GetRegionIdxByName(queryName string) int {
	for i := 0; i < ct.NumEntries(); i++ {
		if ct.Name[i] == queryName {
			return i
		}
	}
	return -1
}

// GetRegionIdxByLabel returns the index of a region in the Colortable by label integer.
// Returns -1 if the region is not found.
//
// Parameters:
//   - queryLabel: the label integer to look up.
//
// Returns:
//   - int: the index of the region, or -1 if it was not found.
func (ct Colortable) GetRegionIdxByLabel(queryLabel int32) int {
	for i := 0; i < ct.NumEntries(); i++ {
		if ct.Label[i] == queryLabel {
			return i
		}
	}
	return -1
}

// Annot models a FreeSurfer annotation, also known as a brain surface parcellation.
// It assigns to each vertex of a brain surface a region (identified by a label integer),
// and contains a Colortable that maps the region labels to region names and colors.
type Annot struct {
	VertexIndices []int32    // The indices of the vertices, these always go from 0 to N-1 (where N is the number of vertices in the respective surface/annotation). Not really needed, but stored in the file.
	VertexLabels  []int32    // The label code for each vertex, defining the region it belongs to. Check in the Colortable for a region that has this label.
	Colortable    Colortable // A Colortable defining the regions (most importantly, the region name and visualization color).
}

// NumVertices returns the number of vertices of this parcellation (or the associated surface).
//
// Returns:
//   - int: the number of vertices.
//   - error: an error if the annot is inconsistent (vertex indices and labels have different lengths).
func (annot Annot) NumVertices() (int, error) {
	nv := len(annot.VertexIndices)
	if len(annot.VertexLabels) != nv {
		return 0, fmt.Errorf("Annot.NumVertices: inconsistent annot, number of vertex indices (%d) and labels (%d) does not match.", nv, len(annot.VertexLabels))
	}
	return nv, nil
}

// RegionVerticesByName returns all vertices of a region given by name in the brain
// surface parcellation.
//
// Parameters:
//   - regionName: the name of the region, e.g., 'superiorfrontal'.
//
// Returns:
//   - []int32: the vertex indices belonging to the region. An empty slice if the region name was not found.
func (annot Annot) RegionVerticesByName(regionName string) []int32 {
	regionIdx := annot.Colortable.GetRegionIdxByName(regionName)
	if regionIdx < 0 {
		return []int32{}
	}
	return annot.RegionVerticesByLabel(annot.Colortable.Label[regionIdx])
}

// RegionVerticesByLabel returns all vertices of a region given by label integer in the
// brain surface parcellation.
//
// Parameters:
//   - regionLabel: the label integer of the region.
//
// Returns:
//   - []int32: the vertex indices belonging to the region.
func (annot Annot) RegionVerticesByLabel(regionLabel int32) []int32 {
	regVerts := []int32{}
	for i := 0; i < len(annot.VertexLabels); i++ {
		if annot.VertexLabels[i] == regionLabel {
			regVerts = append(regVerts, int32(i))
		}
	}
	return regVerts
}

// VertexRegions computes the region indices in the Colortable for all vertices in this
// brain surface parcellation. With the region indices, it becomes very easy to obtain all
// region names, labels, and color channel values from the Colortable. Vertices whose label
// is not contained in the Colortable are assigned region index 0 (typically the 'unknown'
// region).
//
// Returns:
//   - []int: the region index in the Colortable for each vertex.
func (annot Annot) VertexRegions() []int {
	numVertices := len(annot.VertexLabels)
	vertReg := make([]int, numVertices) // init with zeros (region 0, typically 'unknown').
	for regionIdx := 0; regionIdx < annot.Colortable.NumEntries(); regionIdx++ {
		regVertices := annot.RegionVerticesByLabel(annot.Colortable.Label[regionIdx])
		for _, regVertexIdx := range regVertices {
			vertReg[regVertexIdx] = regionIdx
		}
	}
	return vertReg
}

// VertexColors returns the vertex colors as a slice of uint8 values, where 3 consecutive
// values are the red, green and blue channel values for a single vertex. If alpha is true,
// 4 consecutive values are returned per vertex (RGBA).
//
// Parameters:
//   - alpha: whether to include the alpha channel and return 4 values per vertex instead of 3.
//
// Returns:
//   - []uint8: the per-vertex colors, of length (3 or 4) * numVertices.
func (annot Annot) VertexColors(alpha bool) []uint8 {
	numChannels := 3
	if alpha {
		numChannels = 4
	}
	numVertices := len(annot.VertexLabels)
	col := make([]uint8, 0, numVertices*numChannels)
	vertexRegionIndices := annot.VertexRegions()
	for i := 0; i < numVertices; i++ {
		regionIdx := vertexRegionIndices[i]
		col = append(col, uint8(annot.Colortable.R[regionIdx]))
		col = append(col, uint8(annot.Colortable.G[regionIdx]))
		col = append(col, uint8(annot.Colortable.B[regionIdx]))
		if alpha {
			col = append(col, uint8(annot.Colortable.A[regionIdx]))
		}
	}
	return col
}

// VertexRegionNames returns the region names in the Colortable for all vertices in this
// brain surface parcellation.
//
// Returns:
//   - []string: the region name for each vertex.
func (annot Annot) VertexRegionNames() []string {
	regionNames := make([]string, len(annot.VertexLabels))
	vertexRegionIndices := annot.VertexRegions()
	for i := 0; i < len(annot.VertexLabels); i++ {
		regionNames[i] = annot.Colortable.Name[vertexRegionIndices[i]]
	}
	return regionNames
}

// ReadFsAnnot reads a FreeSurfer annotation file (also known as a brain surface
// parcellation) and returns it as an Annot struct.
//
// The file contains a region table (the Colortable) and assigns to each vertex of a
// surface a region. Note that only annotation files in the new format (version 2) that
// contain a color table are supported, which is the case for standard FreeSurfer output
// files like '<subject>/label/lh.aparc.annot' (Desikan-Killiany atlas) or
// '<subject>/label/lh.aparc.a2009s.annot' (Destrieux atlas).
//
// Parameters:
//   - filepath: path to the annotation file, e.g. '<subject>/label/lh.aparc.annot'.
//
// Returns:
//   - Annot: an Annot struct containing the vertex labels and the Colortable.
//   - error: an error if one occurred
func ReadFsAnnot(filepath string) (Annot, error) {

	endian := binary.BigEndian
	annot := Annot{}

	if _, err := os.Stat(filepath); err != nil {
		return annot, fmt.Errorf("ReadFsAnnot: could not stat annot file '%s': %s", filepath, err)
	}

	file, err := os.Open(filepath)
	if err != nil {
		return annot, fmt.Errorf("ReadFsAnnot: could not open annot file '%s': %s", filepath, err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return annot, fmt.Errorf("ReadFsAnnot: could not stat annot file '%s': %s", filepath, err)
	}

	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		return annot, fmt.Errorf("ReadFsAnnot: could not read annot file '%s': %s", filepath, err)
	}

	r := bytes.NewReader(bs)

	// Read the number of vertices.
	var numVertices int32
	if err := binary.Read(r, endian, &numVertices); err != nil {
		return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on number of vertices: %s", err)
	}
	if numVertices <= 0 {
		return annot, fmt.Errorf("ReadFsAnnot: annot file '%s' has invalid num_vertices=%d.", filepath, numVertices)
	}

	if Verbosity > 0 {
		fmt.Printf("ReadFsAnnot: Annot file '%s' has %d vertices.\n", filepath, numVertices)
	}

	// Read the interleaved vertex indices and labels.
	annot.VertexIndices = make([]int32, numVertices)
	annot.VertexLabels = make([]int32, numVertices)
	for i := 0; i < int(numVertices); i++ {
		var vertexIdx int32
		var vertexLabel int32
		if err := binary.Read(r, endian, &vertexIdx); err != nil {
			return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on vertex index %d: %s", i, err)
		}
		if err := binary.Read(r, endian, &vertexLabel); err != nil {
			return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on vertex label %d: %s", i, err)
		}
		annot.VertexIndices[i] = vertexIdx
		annot.VertexLabels[i] = vertexLabel
	}

	// Read the colortable presence flag.
	var hasColortable int32
	if err := binary.Read(r, endian, &hasColortable); err != nil {
		return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on colortable presence flag: %s", err)
	}
	if hasColortable != 1 {
		return annot, fmt.Errorf("ReadFsAnnot: annot file '%s' does not contain a colortable, reading annotations without colortable is not supported.", filepath)
	}

	// Read the colortable format version.
	var colortableFormatVersion int32
	if err := binary.Read(r, endian, &colortableFormatVersion); err != nil {
		return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on colortable format version: %s", err)
	}
	if colortableFormatVersion > 0 {
		return annot, fmt.Errorf("ReadFsAnnot: annot file '%s' uses the old colortable format (version %d), which is not supported.", filepath, colortableFormatVersion)
	}
	formatVersion := -colortableFormatVersion // Negative value means new format, absolute value is the version.
	if formatVersion != 2 {
		return annot, fmt.Errorf("ReadFsAnnot: annot file '%s' uses new format version %d, only version 2 is supported.", filepath, formatVersion)
	}

	// Read the number of colortable entries.
	var numColortableEntries int32
	if err := binary.Read(r, endian, &numColortableEntries); err != nil {
		return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on number of colortable entries: %s", err)
	}
	if numColortableEntries < 0 || numColortableEntries > maxColortableEntries {
		return annot, fmt.Errorf("ReadFsAnnot: annot file '%s' has invalid number of colortable entries (%d), exceeding the maximum of %d.", filepath, numColortableEntries, maxColortableEntries)
	}

	// Read the length of the original filename metadata and skip it.
	var origFilenameLength int32
	if err := binary.Read(r, endian, &origFilenameLength); err != nil {
		return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on original filename length: %s", err)
	}
	if origFilenameLength < 0 || origFilenameLength > maxColortableRegionNameLength {
		return annot, fmt.Errorf("ReadFsAnnot: annot file '%s' has invalid original filename length (%d).", filepath, origFilenameLength)
	}
	origFilenameBytes := make([]byte, origFilenameLength)
	if err := binary.Read(r, endian, &origFilenameBytes); err != nil {
		return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on original filename: %s", err)
	}

	// Read the duplicated number of colortable entries (the format stores it twice).
	var numColortableEntriesDuplicated int32
	if err := binary.Read(r, endian, &numColortableEntriesDuplicated); err != nil {
		return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on duplicated number of colortable entries: %s", err)
	}
	if numColortableEntriesDuplicated != numColortableEntries {
		if Verbosity > 0 {
			fmt.Printf("ReadFsAnnot: Warning: the two num_entries header fields of annot file '%s' do not match (%d vs %d). Use with care.\n",
				filepath, numColortableEntries, numColortableEntriesDuplicated)
		}
	}

	// Read the colortable entries.
	ct := annot.Colortable
	ct.ID = make([]int32, numColortableEntries)
	ct.Name = make([]string, numColortableEntries)
	ct.R = make([]int32, numColortableEntries)
	ct.G = make([]int32, numColortableEntries)
	ct.B = make([]int32, numColortableEntries)
	ct.A = make([]int32, numColortableEntries)
	ct.Label = make([]int32, numColortableEntries)

	for i := 0; i < int(numColortableEntries); i++ {
		var id int32
		if err := binary.Read(r, endian, &id); err != nil {
			return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on colortable entry %d id: %s", i, err)
		}
		var nameLength int32
		if err := binary.Read(r, endian, &nameLength); err != nil {
			return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on colortable entry %d name length: %s", i, err)
		}
		if nameLength <= 0 || nameLength > maxColortableRegionNameLength {
			return annot, fmt.Errorf("ReadFsAnnot: annot file '%s' has invalid region name length (%d) in colortable entry %d.", filepath, nameLength, i)
		}
		nameBytes := make([]byte, nameLength)
		if err := binary.Read(r, endian, &nameBytes); err != nil {
			return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on colortable entry %d name: %s", i, err)
		}
		var rVal int32
		var gVal int32
		var bVal int32
		var aVal int32
		if err := binary.Read(r, endian, &rVal); err != nil {
			return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on colortable entry %d red channel: %s", i, err)
		}
		if err := binary.Read(r, endian, &gVal); err != nil {
			return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on colortable entry %d green channel: %s", i, err)
		}
		if err := binary.Read(r, endian, &bVal); err != nil {
			return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on colortable entry %d blue channel: %s", i, err)
		}
		if err := binary.Read(r, endian, &aVal); err != nil {
			return annot, fmt.Errorf("ReadFsAnnot: binary.Read failed on colortable entry %d alpha channel: %s", i, err)
		}

		ct.ID[i] = id
		ct.Name[i] = string(nameBytes[:nameLength-1]) // Strip the trailing null byte.
		ct.R[i] = rVal
		ct.G[i] = gVal
		ct.B[i] = bVal
		ct.A[i] = aVal
		ct.Label[i] = annotLabelFromRgba(rVal, gVal, bVal, aVal)
	}
	annot.Colortable = ct

	if Verbosity > 0 {
		fmt.Printf("ReadFsAnnot: Annot file '%s' contains %d regions in the colortable.\n", filepath, annot.Colortable.NumEntries())
	}

	return annot, nil
}

// annotLabelFromRgba computes the FreeSurfer annotation label integer from the RGBA
// color channel values, using the formula label = r + g*256 + b*65536 + a*16777216.
func annotLabelFromRgba(rVal int32, gVal int32, bVal int32, aVal int32) int32 {
	label := uint32(uint8(rVal)) |
		uint32(uint8(gVal))<<8 |
		uint32(uint8(bVal))<<16 |
		uint32(uint8(aVal))<<24
	return int32(label)
}

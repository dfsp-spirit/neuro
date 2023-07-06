package neuro

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// Struct modelling a FreeSurfer label.
// A label contains information on a subset of the voxels or vertices only,
// i.e., the number of entries is typically less than the number of voxels
// or vertices in the volume or mesh. Sometimes per-vertex or per-voxel
// data is stored in the labels data field, but sometimes the relevant
// information is simply whether or not a certain element
// (voxel, vertex) is part of the label (e.g., for a cortex label),
// and the per-element data stored is not relevant (and typically set to 0.0).
type FsLabel struct {
	ElementIndex []int32   // The index of the vertex or voxel in the volume or mesh. The first element is 0.
	CoordX   []float32   // The first coordinate of the vertex or voxel in the volume or mesh.
	CoordY   []float32   // The first coordinate of the vertex or voxel in the volume or mesh.
	CoordZ   []float32   // The first coordinate of the vertex or voxel in the volume or mesh.
	Value []float32 // The per-element data.
}


// Check for all vertices in the mesh whether they are part of the label.
//
// Parameters:
//  - label: the label to check
//  - meshNumVertices: the number of vertices in the mesh
//
// Returns:
//  - is_part_of_label: a bool array of length meshNumVertices, where each element is true if the vertex is part of the label, and false otherwise.
//  - error: an error if one occurred, e.g., the number of vertices in the mesh is less than the number of elements in the label.
func vertexIsPartOfLabel(label FsLabel, meshNumVertices int32) ([]bool, error) {
	if meshNumVertices  < int32(len(label.ElementIndex)) {
		err := fmt.Errorf("vertexIsPartOfLabel: number of vertices in mesh (%d) is less than number of elements in label (%d), label invalid for this mesh.", meshNumVertices, len(label.ElementIndex))
		return nil, err
	}
	is_part_of_label := make([]bool, meshNumVertices)  // default value is false
	for _, element_index := range label.ElementIndex {
		is_part_of_label[element_index] = true
	}
	return is_part_of_label, nil
}


// Read an file in FreeSurfer label format.
//
// A label file is a text file representing vertices or voxels in a label. A label contains information on a subset of the voxels or vertices only, i.e., the number of entries is typically less than the number of voxels or vertices in the volume or mesh. Sometimes per-vertex or per-voxel data is stored in the labels data field, but sometimes the real information is whether or not a certain element (voxel, vertex) is part of the label (e.g., for a cortex label), and the per-element data stored is not relevant.
//
// Parameters:
//  - filepath: the path to the file, must be a FreeSurfer label file from recon-all output, like subject/label/lh.cortex.label.
//
// Returns:
//  - pervertex_data: float32 array of per-vertex descriptor values (e.g. cortical thickness)
//  - error: an error if one occurred
func ReadFsLabel(filepath string) (FsLabel, error) {

	var label FsLabel

	// The label format is CSV-like, but with 2 lines at the top that we need to skip.
	// The first line is a comment meant to identify the file type, the second line is a file header field (sort of) and contains the number of elements in the label.
	// There is no header columns that gives the column names, like one would expect in a CSV file.

	// It seems the current CSVReader package is not flexible enough to allow this, so for now,
	// we have to manually skip the first two lines.
	lines, err := readLines(filepath)
	if err != nil {
		return label, err
	}
	if len(lines) <= 2 {
		err = fmt.Errorf("readFsLabel: label file '%s' contains %d lines, but at least 3 required.", filepath, len(lines))
		return label, err
	}

	// Get header field for number of elements in label and check it versus data in file.
	num_rows, err := strconv.Atoi(strings.TrimSpace(lines[1]))
    if err != nil {
		err = fmt.Errorf("readFsLabel: could not convert number of rows (from line 2) in label file '%s' to integer: '%s'.", filepath, err)
		return label, err
    }
	if num_rows != len(lines) -2 {
		err = fmt.Errorf("readFsLabel: number of rows (from line 2) in label file '%s' is %d, but number of lines is %d.", filepath, num_rows, len(lines))
		return label, err
	}

	var linesStartingAtThird strings.Builder
	for _, line := range lines[2:] {
		linesStartingAtThird.WriteString(strings.TrimSpace(line))
		linesStartingAtThird.WriteString("\n")
	}

	r := csv.NewReader(strings.NewReader(linesStartingAtThird.String()))
	r.Comma = ' '
	r.Comment = '#'

	records, err := r.ReadAll()
	if err != nil {
		return label, err
	}

	// Make sure the actual number of records is correct (no comments or empty lines in combination with less records).
	if len(records) != num_rows {
		err = fmt.Errorf("readFsLabel: number of rows (from line 2) in label file '%s' is %d, but number of records is %d.", filepath, num_rows, len(records))
		return label, err
	}

	// Prepare the label struct: set correct length of slices.
	label.ElementIndex = make([]int32, num_rows)
	label.CoordX = make([]float32, num_rows)
	label.CoordY = make([]float32, num_rows)
	label.CoordZ = make([]float32, num_rows)
	label.Value = make([]float32, num_rows)

	var tmpElementIndex int
	var tmpfloat float64

	// Read the data and fill in the label struct.
	for idx, record := range records {

		// The FreeSurfer label files often contain double spaces, which the CSV reader does not like.
		// So we remove all whitespace-fields from the record.
		record_no_whitespace := make([]string, 0)
		for _, field := range record {
			if len(strings.TrimSpace(field)) > 0 {
				record_no_whitespace = append(record_no_whitespace, field)
			}
		}

		if len(record_no_whitespace) != 5 {
			err = fmt.Errorf("readFsLabel: number of columns in label file '%s' record %d is %d, but should be 5: %s.", filepath, idx, len(record_no_whitespace), record_no_whitespace)
			return label, err
		}

		tmpElementIndex, err = strconv.Atoi(record_no_whitespace[0])
		if err != nil {
			err = fmt.Errorf("readFsLabel: could not convert element index in label file '%s' record %d to integer: '%s'.", filepath, idx, err)
			return label, err
		} else {
			label.ElementIndex[idx] = int32(tmpElementIndex)
		}

		tmpfloat, err = strconv.ParseFloat(record_no_whitespace[1], 32)
		if err != nil {
			err = fmt.Errorf("readFsLabel: could not convert X coordinate in label file '%s' record %d to float32: '%s'.", filepath, idx, err)
			return label, err
		} else {
			label.CoordX[idx] = float32(tmpfloat)
		}

		tmpfloat, err = strconv.ParseFloat(record_no_whitespace[2], 32)
		if err != nil {
			err = fmt.Errorf("readFsLabel: could not convert Y coordinate in label file '%s' record %d to float32: '%s'.", filepath, idx, err)
			return label, err
		} else {
			label.CoordY[idx] = float32(tmpfloat)
		}

		tmpfloat, err = strconv.ParseFloat(record_no_whitespace[3], 32)
		if err != nil {
			err = fmt.Errorf("readFsLabel: could not convert Z coordinate in label file '%s' record %d to float32: '%s'.", filepath, idx, err)
			return label, err
		} else {
			label.CoordZ[idx] = float32(tmpfloat)
		}

		tmpfloat, err = strconv.ParseFloat(record_no_whitespace[4], 32)
		if err != nil {
			err = fmt.Errorf("readFsLabel: could not convert value in label file '%s' record %d to float32: '%s'.", filepath, idx, err)
			return label, err
		} else {
			label.Value[idx] = float32(tmpfloat)
		}
	}

	return label, nil
}

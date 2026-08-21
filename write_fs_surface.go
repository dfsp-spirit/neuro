package neuro

// Related software: libfs for C++, see:
// https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L5959 for the fs surface file format.

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

// WriteFsSurface writes a FreeSurfer surface file (a triangular brain mesh) to disk.
//
// The FreeSurfer surface file format is a binary file containing the reconstructed
// surface of a brain hemisphere. The file starts with the magic bytes 255 255 254,
// followed by two newline-terminated ASCII lines (the "created" and "comment" lines),
// the number of vertices and faces (int32 each), and the mesh data (vertices as
// float32 xyz triplets, faces as int32 vertex index triplets). All multi-byte values
// are stored in big endian byte order.
//
// Parameters:
//   - filename: path to the output file, e.g. '<subject>/surf/lh.white'
//   - mesh: the Mesh to write, with Vertices and Faces slices of matching length
//
// Returns:
//   - error: an error if one occurred, or nil on success
func WriteFsSurface(filename string, mesh Mesh) error {

	numVertices := len(mesh.Vertices) / 3
	numFaces := len(mesh.Faces) / 3

	// Basic sanity checks to avoid writing corrupt files.
	if len(mesh.Vertices) == 0 || len(mesh.Vertices)%3 != 0 {
		return fmt.Errorf("WriteFsSurface: mesh has %d vertices (must be a multiple of 3, 3 values per vertex).", len(mesh.Vertices))
	}
	if len(mesh.Faces) == 0 || len(mesh.Faces)%3 != 0 {
		return fmt.Errorf("WriteFsSurface: mesh has %d face indices (must be a multiple of 3, 3 vertex indices per face).", len(mesh.Faces))
	}

	if Verbosity >= 1 {
		fmt.Printf("WriteFsSurface: Writing surface with %d vertices and %d faces to file '%s'.\n", numVertices, numFaces, filename)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("WriteFsSurface: could not create surface file '%s': %s", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	endian := binary.BigEndian

	// Write the magic bytes 255 255 254 (as 3 individual bytes, not an int32).
	magic := [3]byte{255, 255, 254}
	if _, err := writer.Write(magic[:]); err != nil {
		return fmt.Errorf("WriteFsSurface: could not write magic bytes to file '%s': %s", filename, err)
	}

	// Write the created and comment lines, each newline-terminated.
	createdLine := "Created by neurogo"
	commentLine := "FreeSurfer surface file written by the neuro Go module."
	if _, err := writer.WriteString(createdLine + "\n"); err != nil {
		return fmt.Errorf("WriteFsSurface: could not write created line to file '%s': %s", filename, err)
	}
	if _, err := writer.WriteString(commentLine + "\n"); err != nil {
		return fmt.Errorf("WriteFsSurface: could not write comment line to file '%s': %s", filename, err)
	}

	// Write the header with the number of vertices and faces.
	if err := binary.Write(writer, endian, int32(numVertices)); err != nil {
		return fmt.Errorf("WriteFsSurface: could not write number of vertices to file '%s': %s", filename, err)
	}
	if err := binary.Write(writer, endian, int32(numFaces)); err != nil {
		return fmt.Errorf("WriteFsSurface: could not write number of faces to file '%s': %s", filename, err)
	}

	// Write the vertex coordinates (float32) and face indices (int32).
	for _, x := range mesh.Vertices {
		if err := binary.Write(writer, endian, math.Float32bits(x)); err != nil {
			return fmt.Errorf("WriteFsSurface: could not write vertex coordinate to file '%s': %s", filename, err)
		}
	}
	for _, f := range mesh.Faces {
		if err := binary.Write(writer, endian, f); err != nil {
			return fmt.Errorf("WriteFsSurface: could not write face index to file '%s': %s", filename, err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("WriteFsSurface: could not flush surface file '%s': %s", filename, err)
	}

	if Verbosity >= 1 {
		fmt.Printf("WriteFsSurface: Wrote %d vertices and %d faces to file '%s'.\n", numVertices, numFaces, filename)
	}

	return nil
}

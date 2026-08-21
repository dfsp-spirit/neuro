package neuro

// Related software: libfs for C++, see:
// https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L4692 for the annot
// file format.

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

// WriteFsAnnot writes a FreeSurfer annotation (brain surface parcellation) to a file.
//
// The annotation file format (new format, version 2) stores the number of vertices, the
// interleaved vertex indices and labels, and the Colortable containing the region
// metadata (region index, name, and RGBA color channels).
//
// Parameters:
//   - filename: the path to the output file, e.g. '<subject>/label/lh.aparc.annot'.
//   - annot: the Annot struct to write, containing the vertex labels and the Colortable.
//
// Returns:
//   - error: an error if one occurred, or nil on success
func WriteFsAnnot(filename string, annot Annot) error {

	numVertices, err := annot.NumVertices()
	if err != nil {
		return fmt.Errorf("WriteFsAnnot: cannot write annot: %s", err)
	}

	if Verbosity >= 1 {
		fmt.Printf("WriteFsAnnot: Writing annot with %d vertices and %d colortable entries to file '%s'.\n",
			numVertices, annot.Colortable.NumEntries(), filename)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("WriteFsAnnot: could not create annot file '%s': %s", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	endian := binary.BigEndian

	// Write the number of vertices.
	if err := binary.Write(writer, endian, int32(numVertices)); err != nil {
		return fmt.Errorf("WriteFsAnnot: could not write number of vertices to file '%s': %s", filename, err)
	}

	// Write the interleaved vertex indices and labels.
	for i := 0; i < numVertices; i++ {
		if err := binary.Write(writer, endian, annot.VertexIndices[i]); err != nil {
			return fmt.Errorf("WriteFsAnnot: could not write vertex index %d to file '%s': %s", i, filename, err)
		}
		if err := binary.Write(writer, endian, annot.VertexLabels[i]); err != nil {
			return fmt.Errorf("WriteFsAnnot: could not write vertex label %d to file '%s': %s", i, filename, err)
		}
	}

	// Write the colortable presence flag and the format version tag.
	if err := binary.Write(writer, endian, int32(1)); err != nil { // has_colortable
		return fmt.Errorf("WriteFsAnnot: could not write colortable presence flag to file '%s': %s", filename, err)
	}
	if err := binary.Write(writer, endian, int32(-2)); err != nil { // version tag: negative means new format, absolute value is the version.
		return fmt.Errorf("WriteFsAnnot: could not write colortable format version to file '%s': %s", filename, err)
	}

	// Write the number of colortable entries.
	numEntries := annot.Colortable.NumEntries()
	if err := binary.Write(writer, endian, int32(numEntries)); err != nil {
		return fmt.Errorf("WriteFsAnnot: could not write number of colortable entries to file '%s': %s", filename, err)
	}

	// Write the original filename metadata (not meaningful when writing, "unknown" is used as placeholder).
	origFilename := "unknown"
	if err := binary.Write(writer, endian, int32(len(origFilename))); err != nil {
		return fmt.Errorf("WriteFsAnnot: could not write original filename length to file '%s': %s", filename, err)
	}
	if _, err := writer.WriteString(origFilename); err != nil {
		return fmt.Errorf("WriteFsAnnot: could not write original filename to file '%s': %s", filename, err)
	}

	// Write the duplicated number of colortable entries (yes, the format stores it twice).
	if err := binary.Write(writer, endian, int32(numEntries)); err != nil {
		return fmt.Errorf("WriteFsAnnot: could not write duplicated number of colortable entries to file '%s': %s", filename, err)
	}

	// Write the colortable entries.
	for i := 0; i < numEntries; i++ {
		if err := binary.Write(writer, endian, annot.Colortable.ID[i]); err != nil {
			return fmt.Errorf("WriteFsAnnot: could not write colortable entry %d id to file '%s': %s", i, filename, err)
		}
		// Name length: len(name) + 1 for the trailing null byte.
		nameWithNull := annot.Colortable.Name[i] + "\x00"
		if err := binary.Write(writer, endian, int32(len(nameWithNull))); err != nil {
			return fmt.Errorf("WriteFsAnnot: could not write colortable entry %d name length to file '%s': %s", i, filename, err)
		}
		if _, err := writer.WriteString(nameWithNull); err != nil {
			return fmt.Errorf("WriteFsAnnot: could not write colortable entry %d name to file '%s': %s", i, filename, err)
		}
		if err := binary.Write(writer, endian, annot.Colortable.R[i]); err != nil {
			return fmt.Errorf("WriteFsAnnot: could not write colortable entry %d red channel to file '%s': %s", i, filename, err)
		}
		if err := binary.Write(writer, endian, annot.Colortable.G[i]); err != nil {
			return fmt.Errorf("WriteFsAnnot: could not write colortable entry %d green channel to file '%s': %s", i, filename, err)
		}
		if err := binary.Write(writer, endian, annot.Colortable.B[i]); err != nil {
			return fmt.Errorf("WriteFsAnnot: could not write colortable entry %d blue channel to file '%s': %s", i, filename, err)
		}
		if err := binary.Write(writer, endian, annot.Colortable.A[i]); err != nil {
			return fmt.Errorf("WriteFsAnnot: could not write colortable entry %d alpha channel to file '%s': %s", i, filename, err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("WriteFsAnnot: could not flush annot file '%s': %s", filename, err)
	}

	if Verbosity >= 1 {
		fmt.Printf("WriteFsAnnot: Wrote annot with %d vertices and %d colortable entries to file '%s'.\n",
			numVertices, annot.Colortable.NumEntries(), filename)
	}

	return nil
}

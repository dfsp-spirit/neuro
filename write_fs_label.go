package neuro

// Related software: libfs for C++, see:
// https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L6119 for the fs label file format.

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// formatFloatExact formats a float32 so that it can be read back exactly
// (i.e., without precision loss) by a parser using bitSize 32.
func formatFloatExact(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

// WriteFsLabel writes a FreeSurfer label to a file in ASCII label format.
//
// A label file is a text file representing vertices or voxels in a label. The file
// starts with a comment line identifying the file type, a header line with the number
// of entries, and one line per element with the format '<vertex_index> <x> <y> <z> <value>'.
//
// Parameters:
//   - filename: the path to the output file, e.g. 'subject/label/lh.cortex.label'.
//   - label: the label to write.
//
// Returns:
//   - error: an error if one occurred, or nil on success
func WriteFsLabel(filename string, label FsLabel) error {

	numEntries := len(label.ElementIndex)
	if numEntries != len(label.CoordX) || numEntries != len(label.CoordY) ||
		numEntries != len(label.CoordZ) || numEntries != len(label.Value) {
		return fmt.Errorf("WriteFsLabel: label slices have inconsistent lengths (ElementIndex=%d, CoordX=%d, CoordY=%d, CoordZ=%d, Value=%d).",
			numEntries, len(label.CoordX), len(label.CoordY), len(label.CoordZ), len(label.Value))
	}

	if Verbosity >= 1 {
		fmt.Printf("WriteFsLabel: Writing label containing %d vertices to file '%s'.\n", numEntries, filename)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("WriteFsLabel: could not create label file '%s': %s", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	if _, err := writer.WriteString("#!ascii label from subject anonymous\n"); err != nil {
		return fmt.Errorf("WriteFsLabel: could not write comment line to file '%s': %s", filename, err)
	}
	if _, err := writer.WriteString(fmt.Sprintf("%d\n", numEntries)); err != nil {
		return fmt.Errorf("WriteFsLabel: could not write header line to file '%s': %s", filename, err)
	}

	for i := 0; i < numEntries; i++ {
		line := fmt.Sprintf("%d %s %s %s %s\n",
			label.ElementIndex[i],
			formatFloatExact(label.CoordX[i]),
			formatFloatExact(label.CoordY[i]),
			formatFloatExact(label.CoordZ[i]),
			formatFloatExact(label.Value[i]))
		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("WriteFsLabel: could not write label entry %d to file '%s': %s", i, filename, err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("WriteFsLabel: could not flush label file '%s': %s", filename, err)
	}

	if Verbosity >= 1 {
		fmt.Printf("WriteFsLabel: Wrote label containing %d vertices to file '%s'.\n", numEntries, filename)
	}

	return nil
}

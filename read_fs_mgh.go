package neuro

// Related software: libfs for C++, see:
// https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L2023 for the fs surface file format

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)


// MghHeader models the header section of an MGH file.
struct MghHeader {
	MghVersion int32
	Dim1Length int32
	Dim2Length int32
	Dim3Length int32
	Dim4Length int32
	MghDataType int32
	DoF int32
	RasGoodFlag int16

}

// ReadFsSurface reads a FreeSurfer surface file and returns a Mesh struct.
//
// A surface file is a binary file containing the reconstructed surface of a brain hemisphere.
//
// Parameters:
//  - filepath: path to the FreeSurfer mesh file, e.g. '<subject>/surf/lh.white'
//
// Returns:
//  - Mesh: a Mesh struct containing the mesh data
//  - error: an error if one occurred
func ReadFsMghHeader(filepath string) (MghHeader, error) {

}
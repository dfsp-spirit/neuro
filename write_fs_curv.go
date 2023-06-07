package neuro

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

// curvStruct is a struct representing a FreeSurfer curv file.
type curvStruct struct {
	MagicB1 uint8
	MagicB2 uint8
	MagicB3 uint8
	NumVertices int32
	NumFaces int32
	NumValuesPerVertex int32
	Data []float32
}

// curvHeaderStruct is a struct representing a FreeSurfer curv file header, without the data.
type curvHeaderStruct struct {
	MagicB1 uint8
	MagicB2 uint8
	MagicB3 uint8
	NumVertices int32
	NumFaces int32
	NumValuesPerVertex int32
}

// getCurvStruct wraps a CurvStruct around a slice of float32 values.
func getCurvStruct(data[]float32) curvStruct {
	curv := curvStruct{
		MagicB1: 255,
		MagicB2: 255,
		MagicB3: 255,
		NumVertices: int32(len(data)),
		NumFaces: 0,
		NumValuesPerVertex: 1,
		Data: data,
	}
	return curv
}

// getCurvHeaderStructForData creates a CurvStruct around a slice of float32 values.
func getCurvHeaderStruct(data[]float32) curvHeaderStruct {
	curvHdr := curvHeaderStruct{
		MagicB1: 255,
		MagicB2: 255,
		MagicB3: 255,
		NumVertices: int32(len(data)),
		NumFaces: 0,
		NumValuesPerVertex: 1,
	}
	return curvHdr
}

// float32ToByte converts a float32 to a byte array, using big endian byte order.
func float32ToByte(f float32) []byte {

	endian := binary.BigEndian

	var buf [4]byte
	endian.PutUint32(buf[:], math.Float32bits(f))
	return buf[:]
 }

// WriteFsCurv writes a FreeSurfer curv file.
//
// Parameters:
//  - filename: the name of the file to write. Path to it must exist.
//  - data: the slice of float32 values. Must not be empty.
//
// Returns:
//  - error: an error if one occurred, e.g., the slice was empty. Or nil otherwise.
func WriteFsCurv(filename string, data[]float32) error {

	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	curvHdr := getCurvHeaderStruct(data)

	if Verbosity >= 1 {
		fmt.Printf("WriteFsCurv: curvHdr.NumVertices: %d\n", curvHdr.NumVertices)
		fmt.Printf("WriteFsCurv: curvHdr.NumFaces: %d\n", curvHdr.NumFaces)
		fmt.Printf("WriteFsCurv: curvHdr.NumValuesPerVertex: %d\n", curvHdr.NumValuesPerVertex)
	}

	err = binary.Write(file, binary.BigEndian, &curvHdr)
    if err != nil {
        return err
    }

	writer := bufio.NewWriter(file)

	if Verbosity >= 1 {
		fmt.Printf("WriteFsCurv: Writing %d per-vertex descriptor values to file '%s'.\n", len(data), filename)
	}

	numBytesWrittenTotal := 0
	for _, x := range data {
		numBytesWritten, err := writer.Write(float32ToByte(x))
		if err != nil {
			return err
		}
		numBytesWrittenTotal += numBytesWritten
	}

	if Verbosity >= 1 {
		fmt.Printf("WriteFsCurv: Wrote %d bytes to file '%s'.\n", numBytesWrittenTotal, filename)
	}

	writer.Flush()

	return nil
}
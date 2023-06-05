package neurogo

import (
	"encoding/binary"
	"os"
)

// CurvStruct is a struct representing a FreeSurfer curv file.
type CurvStruct struct {
	MagicB1 uint8
	MagicB2 uint8
	MagicB3 uint8
	NumVertices int32
	NumFaces int32
	NumValuesPerVertex int32
	Data []float32
}

// getCurvStruct wraps a CurvStruct around a slice of float32 values.
func getCurvStruct(data[]float32) CurvStruct {
	curv := CurvStruct{
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

	curv := getCurvStruct(data)

	err = binary.Write(file, binary.LittleEndian, curv)
    if err != nil {
        return err
    }

	return nil
}
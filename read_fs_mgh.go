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
// MGH stands for Massachusetts General Hospital, and the MGH format is a binary format for 
// storing 3-dimensional or 4-dimensional structural MRI images of the human brain. The MGZ
// file extension is used for GZIP-compressed files in MGH format. 
type MghHeader struct {
	MghVersion int32
	Dim1Length int32
	Dim2Length int32
	Dim3Length int32
	Dim4Length int32
	MghDataType int32
	DoF int32
	RasGoodFlag int16
}

// getMghDataTypCode translates an MRI data type code (int32) into the respective name (string).
func getMghDataTypeName(dtCode int32) (string, error) {
	// MRI data type representing an 8 bit unsigned integer.
	const MRI_UCHAR int32 = 0
  
	// MRI data type representing a 32 bit signed integer.
  	const MRI_INT int32 = 1

  	// MRI data type representing a 32 bit float.
  	const MRI_FLOAT int32 = 3

  	// MRI data type representing a 16 bit signed integer.
  	const MRI_SHORT int32 = 4

	switch dt := dtCode; dt {
	case MRI_UCHAR:
		return "MRI_UCHAR", nil
	case MRI_INT:
		return "MRI_INT", nil
	case MRI_FLOAT:
		return "MRI_FLOAT", nil
	case MRI_SHORT:
		return "MRI_SHORT", nil
	default:
		err := fmt.Errorf("Invalid or unsupported MGH data type code '%d'. Supported are 0=uchar/uint8, 1=int32, 3=float32, 4=uint16.\n", dtCode)
		return "", err
	}
}

// getMghDataTypCode translates an MRI data type name (string) into the respective header code (int32).
func getMghDataTypCode(dtName string) (int32, error) {
	// MRI data type representing an 8 bit unsigned integer.
	const MRI_UCHAR int32 = 0
  
	// MRI data type representing a 32 bit signed integer.
  	const MRI_INT int32 = 1

  	// MRI data type representing a 32 bit float.
  	const MRI_FLOAT int32 = 3

  	// MRI data type representing a 16 bit signed integer.
  	const MRI_SHORT int32 = 4

	switch dt := dtName; dt {
	case "MRI_UCHAR":
		return MRI_UCHAR, nil
	case "MRI_INT":
		return MRI_INT, nil
	case "MRI_FLOAT":
		return MRI_FLOAT, nil
	case "MRI_SHORT":
		return MRI_SHORT, nil
	default:
		err := fmt.Errorf("Invalid or unsupported MGH data type name '%s'. Supported are 'MRI_UCHAR' (code 0,uchar/uint8), 'MRI_INT'(code 1=int32), 'MRI_FLOAT' (code 3, float32), 'MRI_SHORT' (code 4, uint16).\n", dtName)
		return -1, err
	}
}

// ReadFsMghHeader reads a FreeSurfer MGH file and returns the header as an MghHeader struct.
//
// Parameters:
//  - filepath: path to the FreeSurfer MGH file, e.g. '<subject>/mri/brain.mgh'. Note that gzipped MGH files (file extension .mgz) are currently not supported.
//
// Returns:
//  - MghHeader: an MghHeader struct containing the header data
//  - error: an error if one occurred
func ReadFsMghHeader(filepath string) (MghHeader, error) {
	endian := binary.BigEndian

	hdr := MghHeader{}

	if _, err := os.Stat(filepath); err != nil {
		fmt.Printf("Could not stat MGH file '%s'\n.", filepath)
		return hdr, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		err := fmt.Errorf("Could not open file MGH file '%s' for reading: %s\n", filepath, err)
   		return hdr, err
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
	   fmt.Println(err)
	   return hdr, err
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
	   fmt.Println(err)
	   return hdr, err
	}

	// Read the byte slice
	r := bytes.NewReader(bs)

	if err := binary.Read(r, endian, &hdr); err != nil {
		fmt.Println("ReadFsMghHeader: Read failed on MGH header:", err)
		return hdr, err
	}

	if Verbosity > 0 {
		fmt.Printf("ReadFsMghHeader: Mgh version=%d, dimensions: %d %d %d %d.\n", hdr.MghVersion, hdr.Dim1Length, hdr.Dim2Length, hdr.Dim3Length, hdr.Dim4Length)
		fmt.Printf("ReadFsMghHeader: Mgh data type=%d, DoF=%d, RAS good=%d.\n", hdr.MghDataType, hdr.DoF, hdr.RasGoodFlag)
	}

	if hdr.MghVersion != 1 {
		err := fmt.Errorf("MGH file '%s' is not a valid MGH file or has unsupported file format version (%d), while only version 1 is supported.\n", filepath, hdr.MghVersion)
		return hdr, err
	}

	return hdr, nil

}
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
	"strings"
	"compress/gzip"
)


// MRI data type representing an 8 bit unsigned integer.
const MRI_UCHAR int32 = 0
  
// MRI data type representing a 32 bit signed integer.
const MRI_INT int32 = 1

// MRI data type representing a 32 bit float.
const MRI_FLOAT int32 = 3

// MRI data type representing a 16 bit signed integer.
const MRI_SHORT int32 = 4

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
	RasGoodFlag int16  // flag (1=yes, everything else=no) indicating whether the file contains valid RAS info. 

	// All fields below are so-called RAS info fields. They should be ignored (assumed to contain random data) unless RasGoodFlag is 1.
	XSize float32  // size of voxels in x direction (mm)
	YSize float32  // size of voxels in y direction (mm)
	ZSize float32  // size of voxels in z direction (mm)
	
	Mdc [9]float32 // 9 float values, the 3x3 Mdc matrix that contains image orientation information. The interpretation order is row wise (row1-column1, row1-column2, row1-column3, row2-column1, ...)
	Pxyz_c [3]float32  // 3 float values, the xyz coordinates of the central voxel. Think of the name 'Pxyz_c' as 'Point (x,y,z) coordinates of the center'. 

	// There are 194 more (currently unused) bytes reserved for the header before the data part starts.
	Reserved [194]uint8  // Reached end of header after 284 bytes.
}

// Struct modelling the data part of an MGH file. Only the data in the field identified by MghDataType is valid.
type MghData struct {
	DataMriUchar []uint8
	DataMriInt []int32
	DataMriFloat []float32
	DataMriShort []int16
	MghDataType int32
}

// getMghDataTypeName translates an MRI data type code (int32) into the respective name (string).
func getMghDataTypeName(dtCode int32) (string, error) {

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

// getMghDataTypeCode translates an MRI data type name (string) into the respective header code (int32).
func getMghDataTypeCode(dtName string) (int32, error) {

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


// Determine or guess from filepath and is_gzipped whether the file at filepath is in gzip format. 
func getIsGzipped(filepath string, is_gzipped string) (bool) {
	switch is_gzipped {
	case "yes":
		return true
	case "no":
		return false
	case "auto":
		fpl := strings.ToLower(filepath)
		if strings.HasSuffix(fpl, ".mgz") || strings.HasSuffix(fpl, ".gz") {
			return true
		} else {
			return false
		}
	default:
		panic("is_gzipped must be one of 'yes', 'no', or 'auto'.")
	}
}

func readFileIntoByteSlice(filepath string, treatGzipped bool) ([]byte, error) {
	
	bs := make([]byte, 0)

	if _, err := os.Stat(filepath); err != nil {
		fmt.Printf("Could not stat MGH file '%s'\n.", filepath)
		return bs, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		err := fmt.Errorf("Could not open file MGH file '%s' for reading: %s\n", filepath, err)
   		return bs, err
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
	   fmt.Println (err)
	   return bs, err
	}
	
	// Read the file into a byte slice, whether gzipped or not.
	
	bs = make([]byte, stat.Size())
	if treatGzipped {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return bs, err
		}
		gzipReader.Read(bs)
		defer gzipReader.Close()
	} else {
		_, err := bufio.NewReader(file).Read(bs)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return bs, err
		}
	}
	return bs, nil
}

//gzipReader, err := gzip.NewReader(gzippedFile)
//    defer gzipReader.Close()

// ReadFsMghHeader reads a FreeSurfer MGH file and returns the header as an MghHeader struct.
//
// Parameters:
//  - filepath: path to the FreeSurfer MGH file, e.g. '<subject>/mri/brain.mgh'. Note that gzipped MGH files (file extension .mgz) are currently not supported.
//
// Returns:
//  - MghHeader: an MghHeader struct containing the header data
//  - error: an error if one occurred
func ReadFsMghHeader(filepath string, isGzipped string) (MghHeader, error) {
	endian := binary.BigEndian

	hdr := MghHeader{}

	treatGzipped := getIsGzipped(filepath, isGzipped)
	bs, err := readFileIntoByteSlice(filepath, treatGzipped)
	if err != nil {
		err = fmt.Errorf("Could not read file '%s' into byte slice: %s", filepath, err)
		return hdr, err
	}
	r := bytes.NewReader(bs)

	if err := binary.Read(r, endian, &hdr); err != nil {
		fmt.Println("ReadFsMghHeader: Read failed on MGH header:", err)
		return hdr, err
	}

	dataTypeName, err := getMghDataTypeName(hdr.MghDataType)
	if err != nil {
		return hdr, err
	}

	if Verbosity > 0 {
		fmt.Printf("ReadFsMghHeader: Mgh version=%d, dimensions: %d %d %d %d.\n", hdr.MghVersion, hdr.Dim1Length, hdr.Dim2Length, hdr.Dim3Length, hdr.Dim4Length)		 
		fmt.Printf("ReadFsMghHeader: Mgh data type=%d (%s), DoF=%d, RAS good=%d.\n", hdr.MghDataType, dataTypeName, hdr.DoF, hdr.RasGoodFlag)
	}

	if hdr.MghVersion != 1 {
		err := fmt.Errorf("MGH file '%s' is not a valid MGH file or has unsupported file format version (%d), while only version 1 is supported.\n", filepath, hdr.MghVersion)
		return hdr, err
	}

	if Verbosity > 0 && hdr.RasGoodFlag == 1 {
		fmt.Printf("ReadFsMghHeader: Mgh voxel size x y z: %f, %f, %f.\n", hdr.XSize, hdr.YSize, hdr.ZSize)
		fmt.Printf("ReadFsMghHeader: Mgh Mdc: row0=%f, %f, %f. row1=%f, %f, %f. row2=%f, %f, %f.\n", hdr.Mdc[0], hdr.Mdc[1], hdr.Mdc[2], hdr.Mdc[3], hdr.Mdc[4], hdr.Mdc[5], hdr.Mdc[6], hdr.Mdc[7], hdr.Mdc[8])
		fmt.Printf("ReadFsMghHeader: Mgh Pxyz_c: %f, %f, %f.\n", hdr.Pxyz_c[0], hdr.Pxyz_c[1], hdr.Pxyz_c[2])
	}
	
	return hdr, nil
}

// Mgh models a full MGH format file, including the MghHeader and the MghData.
// See the separate documention for MghHeader and MghData for details on accessing fields.
type Mgh struct {
	header MghHeader
	data MghData
}

func ReadFsMgh(filepath string, isGzipped string) (Mgh, error) {
	var mgh Mgh
	hdr, err := ReadFsMghHeader(filepath, isGzipped)
	if err != nil {
		err := fmt.Errorf("Failed to read MGH header: %s.", err)
		return mgh, err
	}
	mgh.header = hdr
	data, err := ReadFsMghData(filepath, hdr, isGzipped)
	if err != nil {
		err := fmt.Errorf("Failed to read MGH data: %s.", err)
		return mgh, err
	}
	mgh.data = data
	return mgh, nil
}

func ReadFsMghData(filepath string, hdr MghHeader, isGzipped string) (MghData, error) {

	readMghData := MghData{}
	readMghData.MghDataType = -1
	treatGzipped := getIsGzipped(filepath, isGzipped)
	
	switch dt := hdr.MghDataType; dt {
	case MRI_INT:
		arr, err := readFsMghDataMriInt(filepath, hdr, treatGzipped)
		if err != nil {
			err := fmt.Errorf("Failed to read MRI_INT data from MGH file '%s': %s.\n", filepath, err)
			return readMghData, err
		}
		readMghData.DataMriInt = arr
	case MRI_FLOAT:
		arr, err := readFsMghDataMriFloat(filepath, hdr, treatGzipped)
		if err != nil {
			err := fmt.Errorf("Failed to read MRI_FLOAT data from MGH file '%s': %s.\n", filepath, err)
			return readMghData, err
		}
		readMghData.DataMriFloat = arr
	case MRI_UCHAR:
		arr, err := readFsMghDataMriUchar(filepath, hdr, treatGzipped)
		if err != nil {
			err := fmt.Errorf("Failed to read MRI_UCHAR data from MGH file '%s': %s.\n", filepath, err)
			return readMghData, err
		}
		readMghData.DataMriUchar = arr
	case MRI_SHORT:
		arr, err := readFsMghDataMriShort(filepath, hdr, treatGzipped)
		if err != nil {
			err := fmt.Errorf("Failed to read MRI_SHORT data from MGH file '%s': %s.\n", filepath, err)
			return readMghData, err
		}
		readMghData.DataMriShort = arr
	default:
		return readMghData, fmt.Errorf("Header of MGH file '%s' declares unsupported MGH data type code: %d.\n", filepath, dt)
	}
	readMghData.MghDataType = hdr.MghDataType
	return readMghData, nil
}


// readFsMghDataMriInt reads the MRI_INT data array part of an MGH format file.
func readFsMghDataMriInt(filepath string, hdr MghHeader, treatGzipped bool) ([]int32, error) {

	var numValues int64 = (int64)(hdr.Dim1Length * hdr.Dim2Length * hdr.Dim3Length * hdr.Dim4Length)
	dataArr := make([]int32, numValues)
	var mghDataType string = "MRI_INT"

	endian := binary.BigEndian

	bs, err := readFileIntoByteSlice(filepath, treatGzipped)
	if err != nil {
		err = fmt.Errorf("Could not read file '%s' into byte slice: %s", filepath, err)
		return dataArr, err
	}

	// Read the byte slice
	r := bytes.NewReader(bs)

	// Skip the header
	numBytesHeader := 284
	skippedHeaderData := make([]uint8, numBytesHeader)
	if err := binary.Read(r, endian, &skippedHeaderData); err != nil {
		err := fmt.Errorf("ReadFsMghData: binary.Read failed on header part of MGH file %s: %s", filepath, err)
		return dataArr, err
	}

	// Read the data
	if err := binary.Read(r, endian, &dataArr); err != nil {
		err := fmt.Errorf("ReadFsMghData: binary.Read failed on %d values of %s data: %s", numValues, mghDataType, err)
		return dataArr, err
	}

	return dataArr, nil
}

// readFsMghDataMriFloat reads the MRI_FLOAT data array part of an MGH format file.
func readFsMghDataMriFloat(filepath string, hdr MghHeader, treatGzipped bool) ([]float32, error) {

	var numValues int64 = (int64)(hdr.Dim1Length * hdr.Dim2Length * hdr.Dim3Length * hdr.Dim4Length)
	dataArr := make([]float32, numValues)
	var mghDataType string = "MRI_FLOAT"

	endian := binary.BigEndian
	bs, err := readFileIntoByteSlice(filepath, treatGzipped)
	if err != nil {
		err = fmt.Errorf("Could not read file '%s' into byte slice: %s", filepath, err)
		return dataArr, err
	}
	// Read the byte slice
	r := bytes.NewReader(bs)

	// Skip the header
	numBytesHeader := 284
	skippedHeaderData := make([]uint8, numBytesHeader)
	if err := binary.Read(r, endian, &skippedHeaderData); err != nil {
		err := fmt.Errorf("ReadFsMghData: binary.Read failed on header part of MGH file %s: %s", filepath, err)
		return dataArr, err
	}

	// Read the data
	if err := binary.Read(r, endian, &dataArr); err != nil {
		err := fmt.Errorf("ReadFsMghData: binary.Read failed on %d values of %s data: %s", numValues, mghDataType, err)
		return dataArr, err
	}

	return dataArr, nil
}


// readFsMghDataMriuchar reads the MRI_UCHAR data array part of an MGH format file.
func readFsMghDataMriUchar(filepath string, hdr MghHeader, treatGzipped bool) ([]uint8, error) {

	var numValues int64 = (int64)(hdr.Dim1Length * hdr.Dim2Length * hdr.Dim3Length * hdr.Dim4Length)
	dataArr := make([]uint8, numValues)
	var mghDataType string = "MRI_UCHAR"

	endian := binary.BigEndian
	bs, err := readFileIntoByteSlice(filepath, treatGzipped)
	if err != nil {
		err = fmt.Errorf("Could not read file '%s' into byte slice: %s", filepath, err)
		return dataArr, err
	}
	// Read the byte slice
	r := bytes.NewReader(bs)

	// Skip the header
	numBytesHeader := 284
	skippedHeaderData := make([]uint8, numBytesHeader)
	if err := binary.Read(r, endian, &skippedHeaderData); err != nil {
		err := fmt.Errorf("ReadFsMghData: binary.Read failed on header part of MGH file %s: %s", filepath, err)
		return dataArr, err
	}

	// Read the data
	if err := binary.Read(r, endian, &dataArr); err != nil {
		err := fmt.Errorf("ReadFsMghData: binary.Read failed on %d values of %s data: %s", numValues, mghDataType, err)
		return dataArr, err
	}

	return dataArr, nil
}


// readFsMghDataMrishort reads the MRI_SHORT data array part of an MGH format file.
func readFsMghDataMriShort(filepath string, hdr MghHeader, treatGzipped bool) ([]int16, error) {

	var numValues int64 = (int64)(hdr.Dim1Length * hdr.Dim2Length * hdr.Dim3Length * hdr.Dim4Length)
	dataArr := make([]int16, numValues)
	var mghDataType string = "MRI_SHORT"

	endian := binary.BigEndian

	bs, err := readFileIntoByteSlice(filepath, treatGzipped)
	if err != nil {
		err = fmt.Errorf("Could not read file '%s' into byte slice: %s", filepath, err)
		return dataArr, err
	}

	// Read the byte slice
	r := bytes.NewReader(bs)

	// Skip the header
	numBytesHeader := 284
	skippedHeaderData := make([]uint8, numBytesHeader)
	if err := binary.Read(r, endian, &skippedHeaderData); err != nil {
		err := fmt.Errorf("ReadFsMghDataMriShort: binary.Read failed on header part of MGH file %s: %s", filepath, err)
		return dataArr, err
	}

	// Read the data
	if err := binary.Read(r, endian, &dataArr); err != nil {
		err := fmt.Errorf("ReadFsMghDataMriShort: binary.Read failed on %d values of %s data: %s", numValues, mghDataType, err)
		return dataArr, err
	}

	return dataArr, nil
}
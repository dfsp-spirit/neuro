package neuro

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

// MRI data type representing an 8 bit unsigned integer. Used by MGH format, see MghHeader struct.
const MRI_UCHAR int32 = 0

// MRI data type representing a 32 bit signed integer. Used by MGH format, see MghHeader struct.
const MRI_INT int32 = 1

// MRI data type representing a 32 bit float. Used by MGH format, see MghHeader struct.
const MRI_FLOAT int32 = 3

// MRI data type representing a 16 bit signed integer. Used by MGH format, see MghHeader struct.
const MRI_SHORT int32 = 4

// MghHeader models the header section of an MGH file.
// MGH stands for Massachusetts General Hospital, and the MGH format is a binary format for
// storing 3-dimensional or 4-dimensional structural MRI images of the human brain. The MGZ
// file extension is used for GZIP-compressed files in MGH format.
type MghHeader struct {
	MghVersion  int32 // version of the MGH file format. Currently, this is always 1.
	Dim1Length  int32 // number of voxels in x direction
	Dim2Length  int32 // number of voxels in y direction
	Dim3Length  int32 // number of voxels in z direction
	Dim4Length  int32 // number of voxels in 4th dimension (typically time or subject index)
	MghDataType int32 // MRI data type code. See MRI_UCHAR, MRI_INT, MRI_FLOAT, MRI_SHORT constants in this package.
	DoF         int32
	RasGoodFlag int16 // flag (1=yes, everything else=no) indicating whether the file contains valid RAS info.

	// All fields below are so-called RAS info fields. They should be ignored (assumed to contain random data) unless RasGoodFlag is 1.
	XSize float32 // size of voxels in x direction (mm)
	YSize float32 // size of voxels in y direction (mm)
	ZSize float32 // size of voxels in z direction (mm)

	Mdc    [9]float32 // 9 float values, the 3x3 Mdc matrix that contains image orientation information. The interpretation order is row wise (row1-column1, row1-column2, row1-column3, row2-column1, ...)
	Pxyz_c [3]float32 // 3 float values, the xyz coordinates of the central voxel. Think of the name 'Pxyz_c' as 'Point (x,y,z) coordinates of the center'.

	// There are 194 more (currently unused) bytes reserved for the header before the data part starts.
	Reserved [194]uint8 // Reached end of header after 284 bytes. Data in here should be considered random.
}

// Struct modelling the data part of an MGH file. Only the data in the field identified by MghDataType is valid.
type MghData struct {
	DataMriUchar []uint8   // The data, if MghDataType is MRI_UCHAR. Otherwise this field contains random data.
	DataMriInt   []int32   // The data, if MghDataType is MRI_INT. Otherwise this field contains random data.
	DataMriFloat []float32 // The data, if MghDataType is MRI_FLOAT. Otherwise this field contains random data.
	DataMriShort []int16   // The data, if MghDataType is MRI_SHORT. Otherwise this field contains random data.
	MghDataType  int32     // The MRI data type code. See MRI_UCHAR, MRI_INT, MRI_FLOAT, MRI_SHORT. Use this to determine which of the data fields above is valid.
}

// getMghDataTypeName translates an MRI data type code (int32) into the respective name (string).
//
// Parameters:
//   - dtCode: The MRI data type code, e.g., integer constants MRI_UCHAR, MRI_INT, MRI_FLOAT, MRI_SHORT.
//
// Returns:
//   - string: The MRI data type name, e.g., "MRI_UCHAR", "MRI_INT", "MRI_FLOAT", "MRI_SHORT".
//   - error: Error if any, e.g., if the data type code is invalid or unsupported.
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
//
// Parameters:
//   - dtName: Name of the MRI data type, e.g., "MRI_UCHAR", "MRI_INT", "MRI_FLOAT", "MRI_SHORT".
//
// Returns:
//   - int32: The MRI data type code, e.g., integer constants MRI_UCHAR, MRI_INT, MRI_FLOAT, MRI_SHORT.
//   - error: Error if any, e.g., on invalid dtName.
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
//
// Parameters:
//   - filepath: Path to the file to be read.
//   - isGzipped: Whether to treat the file as gzip-compressed.
//
// Returns:
//   - bool: Whether the file is in gzip format.
func getIsGzipped(filepath string, isGzipped string) bool {
	switch isGzipped {
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

// getIsGzippedMgh translates special is_gzipped values "mgh" and "mgz" into "no" and "yes", respectively.
//
// Parameters:
//   - isGzipped: Whether to treat the file as gzip-compressed. If it is one of "mgh" or "mgz", it will be translated to "no" and "yes", respectively.
//
// Returns:
//   - string: "yes" if isGzipped is "mgz", "no" if isGzipped is "mgh", otherwise isGzipped.
func getIsGzippedMgh(isGzipped string) string {
	if isGzipped == "mgh" {
		return "no"
	}
	if isGzipped == "mgz" {
		return "yes"
	}
	return isGzipped
}

// readFileIntoByteSlice reads a file from disk into a byte slice. Supports uncompressed or gzip format.
//
// Parameters:
// filepath: Path to the file to be read.
// treatGzipped: Whether to treat the file as gzip-compressed.
//
// Returns:
// bs: The file contents as a byte slice.
// err: An error, if any.
func readFileIntoByteSlice(filepath string, treatGzipped bool) ([]byte, error) {

	bs := make([]byte, 0)

	if _, err := os.Stat(filepath); err != nil {
		fmt.Printf("Could not stat file '%s'\n.", filepath)
		return bs, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		err := fmt.Errorf("Could not open file file '%s' for reading: %s\n", filepath, err)
		return bs, err
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return bs, err
	}

	// Read the file into a byte slice, whether gzipped or not.
	bs = make([]byte, stat.Size())
	if treatGzipped {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return bs, err
		}
		defer gzipReader.Close()
		//gzipReader.Read(bs)
		bs, err = io.ReadAll(gzipReader)
		if err != nil {
			return bs, err
		}
	} else {
		_, err := bufio.NewReader(file).Read(bs)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return bs, err
		}
	}
	return bs, nil
}

// ReadFsMghHeader reads a FreeSurfer MGH file and returns the header as an MghHeader struct.
//
// See the documentation of the MghHeader struct for details on the header fields.
//
// Parameters:
//   - filepath: path to the FreeSurfer MGH file, e.g. '<subject>/mri/brain.mgh'. Note that gzipped MGH files (file extension .mgz) are currently not supported.
//   - isGzipped: Whether to treat the file as gzip-compressed. If "auto", the file extension is used to determine whether the file is gzip-compressed. If not "auto", it has to be "yes"/"mgz" or "no"/"mgh" to force MGZ or MGH format, respectively.
//
// Returns:
//   - MghHeader: an MghHeader struct containing the header data
//   - error: an error if one occurred
func ReadFsMghHeader(filepath string, isGzipped string) (MghHeader, error) {
	endian := binary.BigEndian

	hdr := MghHeader{}

	isGzipped = getIsGzippedMgh(isGzipped)
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
	Header MghHeader
	Data   MghData
}

// ReadFsMgh reads a FreeSurfer MGH file and returns it as an Mgh struct. The Mgh struct contains the MghHeader and MghData.
//
// See the documentation for Mgh, MghHeader and MghData for details on accessing fields.
//
// Parameters:
//   - filepath: path to the FreeSurfer MGH file, e.g. '<subject>/mri/brain.mgh'. Note that gzipped MGH files (file extension .mgz) are currently not supported.
//   - isGzipped: Whether to treat the file as gzip-compressed. If "auto", the file extension is used to determine whether the file is gzip-compressed. If not "auto", it has to be "yes"/"mgz" or "no"/"mgh" to force MGZ or MGH format, respectively.
//
// Returns:
//   - Mgh: an Mgh struct containing the MghHeader and MghData
func ReadFsMgh(filepath string, isGzipped string) (Mgh, error) {
	var mgh Mgh
	hdr, err := ReadFsMghHeader(filepath, isGzipped)
	if err != nil {
		err := fmt.Errorf("Failed to read MGH header: %s.", err)
		return mgh, err
	}
	mgh.Header = hdr
	data, err := ReadFsMghData(filepath, hdr, isGzipped)
	if err != nil {
		err := fmt.Errorf("Failed to read MGH data: %s.", err)
		return mgh, err
	}
	mgh.Data = data
	return mgh, nil
}

// ReadFsMghData reads the data part of an MGH or MGZ format file into an MghData struct.
//
// See the documentation for MghData for details on accessing fields.
//
// Parameters:
//   - filepath: path to readable input file in MGH or MGZ format
//   - hdr: MghHeader struct containing the header data
//   - isGzipped: string indicating whether the input file is gzipped or not. If set to 'auto', the function will try to determine this automatically.
//
// Returns:
//   - MghData: an MghData struct containing the data
//   - error: an error if one occurred, nil otherwise
func ReadFsMghData(filepath string, hdr MghHeader, isGzipped string) (MghData, error) {

	readMghData := MghData{}
	readMghData.MghDataType = -1
	isGzipped = getIsGzippedMgh(isGzipped)
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
//
// Parameters:
//   - filepath: path to readable input file in MGH or MGZ format
//   - hdr: MghHeader struct containing the header data
//   - treatGzipped: boolean indicating whether the input file is gzipped or not
//
// Returns:
//   - []int32: an array of int32 values. You will have to reshape this 1D array to a 4D array with the dimensions given in the MghHeader.
//   - error: an error if one occurred, nil otherwise
func readFsMghDataMriInt(filepath string, hdr MghHeader, treatGzipped bool) ([]int32, error) {

	var numValues int64 = (int64)(hdr.Dim1Length * hdr.Dim2Length * hdr.Dim3Length * hdr.Dim4Length)
	dataArr := make([]int32, numValues)
	var mghDataType string = "MRI_INT"

	if Verbosity >= 1 {
		fmt.Printf("Reading %d values of type %s from MGH file '%s', treatGzipped=%t\n", numValues, mghDataType, filepath, treatGzipped)
	}

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
//
// Parameters:
//   - filepath: path to readable input file in MGH or MGZ format
//   - hdr: MghHeader struct containing the header data
//   - treatGzipped: boolean indicating whether the input file is gzipped or not
//
// Returns:
//   - []float32: an array of float32 values. You will have to reshape this 1D array to a 4D array with the dimensions given in the MghHeader.
//   - error: an error if one occurred, nil otherwise
func readFsMghDataMriFloat(filepath string, hdr MghHeader, treatGzipped bool) ([]float32, error) {

	var numValues int64 = (int64)(hdr.Dim1Length * hdr.Dim2Length * hdr.Dim3Length * hdr.Dim4Length)
	dataArr := make([]float32, numValues)
	var mghDataType string = "MRI_FLOAT"

	if Verbosity >= 1 {
		fmt.Printf("Reading %d values of type %s from MGH file '%s', treatGzipped=%t\n", numValues, mghDataType, filepath, treatGzipped)
	}

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
//
// Parameters:
//   - filepath: path to readable input file in MGH or MGZ format
//   - hdr: MghHeader struct containing the header data
//   - treatGzipped: boolean indicating whether the input file is gzipped or not
//
// Returns:
//   - []uint8: an array of uint8 values. You will have to reshape this 1D array to a 4D array with the dimensions given in the MghHeader.
//   - error: an error if one occurred, nil otherwise
func readFsMghDataMriUchar(filepath string, hdr MghHeader, treatGzipped bool) ([]uint8, error) {

	var numValues int64 = (int64)(hdr.Dim1Length * hdr.Dim2Length * hdr.Dim3Length * hdr.Dim4Length)
	dataArr := make([]uint8, numValues)
	var mghDataType string = "MRI_UCHAR"

	if Verbosity >= 1 {
		fmt.Printf("Reading %d values of type %s from MGH file '%s', treatGzipped=%t\n", numValues, mghDataType, filepath, treatGzipped)
	}

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

// readFsMghDataMriShort reads the MRI_SHORT data array part of an MGH format file.
//
// Parameters:
//   - filepath: path to readable input file in MGH or MGZ format
//   - hdr: MghHeader struct containing the header data
//   - treatGzipped: boolean indicating whether the input file is gzipped or not
//
// Returns:
//   - []int16: an array of int16 values. You will have to reshape this 1D array to a 4D array with the dimensions given in the MghHeader.
//   - error: an error if one occurred, nil otherwise
func readFsMghDataMriShort(filepath string, hdr MghHeader, treatGzipped bool) ([]int16, error) {

	var numValues int64 = (int64)(hdr.Dim1Length * hdr.Dim2Length * hdr.Dim3Length * hdr.Dim4Length)
	dataArr := make([]int16, numValues)
	var mghDataType string = "MRI_SHORT"

	if Verbosity >= 1 {
		fmt.Printf("Reading %d values of type %s from MGH file '%s', treatGzipped=%t\n", numValues, mghDataType, filepath, treatGzipped)
	}

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

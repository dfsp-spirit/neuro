package neuro

// Related software: libfs for C++, see:
// https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L4815 for the MGH file format.

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"os"
)

// writeByteSliceToFile writes a byte slice to a file, optionally gzip-compressed.
//
// Parameters:
//   - bs: the byte slice to write
//   - filepath: the path to the output file
//   - gzipIt: whether to gzip-compress the data before writing it to the file
//
// Returns:
//   - error: an error if one occurred, or nil on success
func writeByteSliceToFile(bs []byte, filepath string, gzipIt bool) error {

	if Verbosity >= 1 {
		fmt.Printf("writeByteSliceToFile: Writing %d bytes to file '%s', gzip=%t.\n", len(bs), filepath, gzipIt)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("writeByteSliceToFile: could not create file '%s': %s", filepath, err)
	}
	defer file.Close()

	if gzipIt {
		gzipWriter := gzip.NewWriter(file)
		_, err = gzipWriter.Write(bs)
		if err != nil {
			return fmt.Errorf("writeByteSliceToFile: could not write gzip-compressed data to file '%s': %s", filepath, err)
		}
		if err := gzipWriter.Close(); err != nil {
			return fmt.Errorf("writeByteSliceToFile: could not finish gzip compression for file '%s': %s", filepath, err)
		}
	} else {
		_, err = file.Write(bs)
		if err != nil {
			return fmt.Errorf("writeByteSliceToFile: could not write data to file '%s': %s", filepath, err)
		}
	}

	return nil
}

// mghToByteSlice serializes an Mgh struct into the binary MGH format (not compressed).
//
// Parameters:
//   - mgh: the Mgh struct to serialize, containing the MghHeader and MghData.
//
// Returns:
//   - []byte: the byte representation of the MGH file content.
//   - error: an error if one occurred, e.g., on unsupported data types or data/header size mismatches.
func mghToByteSlice(mgh Mgh) ([]byte, error) {

	endian := binary.BigEndian
	var buf bytes.Buffer

	// Header part 1: version, dimensions, data type, dof.
	headerFields := []int32{
		1, // MGH file format version, currently always 1.
		mgh.Header.Dim1Length,
		mgh.Header.Dim2Length,
		mgh.Header.Dim3Length,
		mgh.Header.Dim4Length,
		mgh.Header.MghDataType,
		mgh.Header.DoF,
	}
	for _, val := range headerFields {
		if err := binary.Write(&buf, endian, val); err != nil {
			return nil, fmt.Errorf("mghToByteSlice: could not write MGH header field: %s", err)
		}
	}

	// RAS good flag.
	if err := binary.Write(&buf, endian, mgh.Header.RasGoodFlag); err != nil {
		return nil, fmt.Errorf("mghToByteSlice: could not write MGH RAS good flag: %s", err)
	}

	// RAS info part of the header (60 bytes), only valid if the RAS good flag is set.
	if mgh.Header.RasGoodFlag == 1 {
		rasFloatFields := []float32{
			mgh.Header.XSize,
			mgh.Header.YSize,
			mgh.Header.ZSize,
		}
		for _, val := range rasFloatFields {
			if err := binary.Write(&buf, endian, val); err != nil {
				return nil, fmt.Errorf("mghToByteSlice: could not write MGH RAS size field: %s", err)
			}
		}
		for i := 0; i < 9; i++ {
			if err := binary.Write(&buf, endian, mgh.Header.Mdc[i]); err != nil {
				return nil, fmt.Errorf("mghToByteSlice: could not write MGH Mdc value %d: %s", i, err)
			}
		}
		for i := 0; i < 3; i++ {
			if err := binary.Write(&buf, endian, mgh.Header.Pxyz_c[i]); err != nil {
				return nil, fmt.Errorf("mghToByteSlice: could not write MGH Pxyz_c value %d: %s", i, err)
			}
		}
	} else {
		// If the RAS info is not valid, write 60 zero bytes (3 sizes + 9 Mdc + 3 Pxyz_c values).
		zeroRas := make([]byte, 60)
		if _, err := buf.Write(zeroRas); err != nil {
			return nil, fmt.Errorf("mghToByteSlice: could not write MGH RAS placeholder bytes: %s", err)
		}
	}

	// Remaining reserved header bytes (194), so the header is 284 bytes in total.
	reservedBytes := make([]byte, 194)
	if _, err := buf.Write(reservedBytes); err != nil {
		return nil, fmt.Errorf("mghToByteSlice: could not write MGH reserved header bytes: %s", err)
	}

	// Validate the data size against the header and write the data part.
	numValues := int64(mgh.Header.Dim1Length) * int64(mgh.Header.Dim2Length) * int64(mgh.Header.Dim3Length) * int64(mgh.Header.Dim4Length)

	switch dt := mgh.Header.MghDataType; dt {
	case MRI_UCHAR:
		if int64(len(mgh.Data.DataMriUchar)) != numValues {
			return nil, fmt.Errorf("mghToByteSlice: MRI_UCHAR data has %d values, but MGH header declares %d.", len(mgh.Data.DataMriUchar), numValues)
		}
		for _, val := range mgh.Data.DataMriUchar {
			if err := binary.Write(&buf, endian, val); err != nil {
				return nil, fmt.Errorf("mghToByteSlice: could not write MRI_UCHAR data value: %s", err)
			}
		}
	case MRI_INT:
		if int64(len(mgh.Data.DataMriInt)) != numValues {
			return nil, fmt.Errorf("mghToByteSlice: MRI_INT data has %d values, but MGH header declares %d.", len(mgh.Data.DataMriInt), numValues)
		}
		for _, val := range mgh.Data.DataMriInt {
			if err := binary.Write(&buf, endian, val); err != nil {
				return nil, fmt.Errorf("mghToByteSlice: could not write MRI_INT data value: %s", err)
			}
		}
	case MRI_FLOAT:
		if int64(len(mgh.Data.DataMriFloat)) != numValues {
			return nil, fmt.Errorf("mghToByteSlice: MRI_FLOAT data has %d values, but MGH header declares %d.", len(mgh.Data.DataMriFloat), numValues)
		}
		for _, val := range mgh.Data.DataMriFloat {
			if err := binary.Write(&buf, endian, val); err != nil {
				return nil, fmt.Errorf("mghToByteSlice: could not write MRI_FLOAT data value: %s", err)
			}
		}
	case MRI_SHORT:
		if int64(len(mgh.Data.DataMriShort)) != numValues {
			return nil, fmt.Errorf("mghToByteSlice: MRI_SHORT data has %d values, but MGH header declares %d.", len(mgh.Data.DataMriShort), numValues)
		}
		for _, val := range mgh.Data.DataMriShort {
			if err := binary.Write(&buf, endian, val); err != nil {
				return nil, fmt.Errorf("mghToByteSlice: could not write MRI_SHORT data value: %s", err)
			}
		}
	default:
		return nil, fmt.Errorf("mghToByteSlice: unsupported MGH data type code %d, cannot write data.", mgh.Header.MghDataType)
	}

	return buf.Bytes(), nil
}

// WriteFsMgh writes an Mgh struct to a file in FreeSurfer MGH or MGZ format.
//
// The MGH format is a binary, big-endian FreeSurfer file format for storing 4D data.
// Several data types are supported (see MRI_UCHAR, MRI_INT, MRI_FLOAT, MRI_SHORT constants),
// and one has to check the header to see which one is contained in a file. The MGZ format
// is just a gzip-compressed MGH file.
//
// Parameters:
//   - filename: path to the output file, e.g. '<subject>/mri/brain.mgh' or '<subject>/mri/brain.mgz'
//   - mgh: the Mgh struct to write, containing the MghHeader and MghData
//   - isGzipped: whether to write the file gzip-compressed (MGZ format). If "auto", the file
//     extension is used to determine whether the file should be gzip-compressed. If not "auto",
//     it has to be "yes"/"mgz" or "no"/"mgh" to force MGZ or MGH format, respectively.
//
// Returns:
//   - error: an error if one occurred, or nil on success
func WriteFsMgh(filename string, mgh Mgh, isGzipped string) error {

	bs, err := mghToByteSlice(mgh)
	if err != nil {
		return fmt.Errorf("WriteFsMgh: could not serialize Mgh struct to MGH format: %s", err)
	}

	isGzipped = getIsGzippedMgh(isGzipped)
	gzipIt := getIsGzipped(filename, isGzipped)

	if Verbosity >= 1 {
		fmt.Printf("WriteFsMgh: Writing %d bytes of MGH data to file '%s', gzip=%t.\n", len(bs), filename, gzipIt)
	}

	err = writeByteSliceToFile(bs, filename, gzipIt)
	if err != nil {
		return fmt.Errorf("WriteFsMgh: could not write MGH data to file '%s': %s", filename, err)
	}

	return nil
}

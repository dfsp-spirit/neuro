package neuro

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// Read a binary file in FreeSurfer curv format.
//
// Curv files are used to store per-vertex descriptors like cortical thickness in native space (i.e., for a single subject, not mapped to a group template).
//
// Parameters:
//  - filepath: the path to the file, must be a FreeSurfer curv file from recon-all output, like subject/surf/lh.thickness.
//
// Returns:
//  - pervertex_data: float32 array of per-vertex descriptor values (e.g. cortical thickness)
//  - error: an error if one occurred
func ReadFsCurv(filepath string) ([]float32, error) {
	//b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40, 0xff, 0x01, 0x02, 0x03, 0xbe, 0xef}
	//r := bytes.NewReader(b)

	endian := binary.BigEndian
	pervertex_data := []float32{}

	if _, err := os.Stat(filepath); err != nil {
		fmt.Printf("Could not stat file '%s'\n.", filepath)
		return pervertex_data, err
	}

	file, err := os.Open(filepath)
	if err != nil {
   		panic(err)
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
	   fmt.Println(err)
	   return pervertex_data, err
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
	   fmt.Println(err)
	   return pervertex_data, err
	}

	// Read the byte slice
	r := bytes.NewReader(bs)

	type curvHeaderPart1 struct {
		MagicB1 uint8
		MagicB2 uint8
		MagicB3 uint8
	}

	hdr1 := curvHeaderPart1{}


	if err := binary.Read(r, endian, &hdr1); err != nil {
		fmt.Println("ReadFsCurv: binary.Read failed on curv header part 1:", err)
		return pervertex_data, err
	}

	if Verbosity > 0 {
		fmt.Printf("ReadFsCurv: Curv header magic bytes: %d %d %d.\n", hdr1.MagicB1, hdr1.MagicB2, hdr1.MagicB3)
	}


	if ! (hdr1.MagicB1 == 255 && hdr1.MagicB2 == 255 && hdr1.MagicB3 == 255) {
		fmt.Println("ReadFsCurv: Error: Curv magic bytes are not 255 255 255, this is not a FreeSurfer curv file. Provide a recon-all output file like '<subject>/surf/lh.thickness'.")
		return pervertex_data, err
	}

	type curvHeaderPart2 struct {
		NumVertices int32
		NumFaces int32
		NumValuesPerVertex int32
	}

	hdr2 := curvHeaderPart2{}


	if err := binary.Read(r, endian, &hdr2); err != nil {
		fmt.Println("ReadFsCurv: binary.Read failed on curv header part 2:", err)
		return pervertex_data, err
	}

	if Verbosity > 0 {
		fmt.Println("ReadFsCurv: NumVertices:", hdr2.NumVertices)
		fmt.Println("ReadFsCurv: NumFaces:", hdr2.NumFaces)
		fmt.Println("ReadFsCurv: NumValuesPerVertex:", hdr2.NumValuesPerVertex)
	}

	// read per-vertex data
	pervertex_data = make([]float32, hdr2.NumVertices) 	// one descriptor value per vertex
	if err := binary.Read(r, endian, &pervertex_data); err != nil {
		fmt.Println("ReadFsCurv: binary.Read failed on per-vertex descriptor slice:", err)
		return pervertex_data, err
	}

	return pervertex_data, nil
}


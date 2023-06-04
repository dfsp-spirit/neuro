package neurogo

// Related packages and documentation:
// https://pkg.go.dev/github.com/oschwald/maxminddb-golang#example-Reader.Lookup-Interface
// https://pkg.go.dev/encoding/binary#example-Read-Multi
// maybe https://www.jonathan-petitcolas.com/2014/09/25/parsing-binary-files-in-go.html, but it's old
//
// https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L2023 for the fs surface file format

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Read a newline-terminated string from a bytes.Reader.
func readNewlineTerminatedString(r *bytes.Reader, endian binary.ByteOrder, do_strip_newline bool) (string, error) {

	//endian = binary.BigEndian // TODO: make this a parameter

	//return "not implemented yet", nil
	var line string = ""
	char := make([]byte, 1)
	var char_num int = 0
	for string(char[:]) != "\n" {
		if err := binary.Read(r, endian, &char); err != nil {
			fmt.Printf("binary.Read failed on character %d of newline-terminated string: %s\n", char_num, err)
			return "", err
		} else {
			if !(do_strip_newline && string(char[:]) == "\n") {
				line = line + string(char[:])
			}
		}
		char_num++
	}
	return line, nil
}

func _magicByte3ToInt(magic []byte) (int) {
	int1, _ := strconv.Atoi(string(magic[:])) // FIXME: incorrect
	int1 = ((int1 >> 8) & 0xffffff);
	return int1
}


// ReadFreesurferMesh reads a FreeSurfer mesh file and returns a Mesh struct.
func ReadFreesurferMesh(filepath string) (Mesh, error) {
	//b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40, 0xff, 0x01, 0x02, 0x03, 0xbe, 0xef}
	//r := bytes.NewReader(b)

	endian := binary.BigEndian
	surface := Mesh{}

	if _, err := os.Stat(filepath); err != nil {
		fmt.Printf("Could not stat file '%s'\n.", filepath)
		return surface, err
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
	   return surface, err
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
	   fmt.Println(err)
	   return surface, err
	}

	// Read the byte slice
	r := bytes.NewReader(bs)

	type header_part1 struct {
		MagicB1 uint8
		MagicB2 uint8
		MagicB3 uint8
		//Mine [3]byte
	}

	hdr1 := header_part1{}


	if err := binary.Read(r, endian, &hdr1); err != nil {
		fmt.Println("binary.Read failed on first part of fs surface header:", err)
		return surface, err
	}


	if Verbosity > 0 {
		fmt.Println("hdr1.Magic_b1:", hdr1.MagicB1)
		fmt.Println("hdr1.Magic_b2:", hdr1.MagicB2)
		fmt.Println("hdr1.Magic_b3:", hdr1.MagicB3)
	}

	magic := make([]byte, 3)
	magic[0] = hdr1.MagicB1
	magic[1] = hdr1.MagicB2
	magic[2] = hdr1.MagicB3
	int1 := _magicByte3ToInt(magic)
	fmt.Printf("Header magic bytes %d %d %d gives: %d.\n", hdr1.MagicB1, hdr1.MagicB2, hdr1.MagicB3, int1)

	createdLine, err := readNewlineTerminatedString(r, endian, true);
    commentLine, err := readNewlineTerminatedString(r, endian, true);

	if Verbosity > 0 {
		fmt.Printf("createdLine: '%s'\n", createdLine)
		fmt.Printf("commentLine: '%s'\n", commentLine)
	}

	type header_part2 struct {
		NumVerts int32
		NumFaces int32
	}

	hdr2 := header_part2{}

	if err := binary.Read(r, endian, &hdr2); err != nil {
		fmt.Println("binary.Read failed on second part of fs surface header:", err)
		return surface, err
	}

	if Verbosity > 0 {
		fmt.Println("NumVerts:", hdr2.NumVerts)
		fmt.Println("NumFaces:", hdr2.NumFaces)
	}

	// read mesh data
	surface.Vertices = make([]float32, hdr2.NumVerts * 3) // x,y,z coordinates for each vertex
	surface.Faces = make([]int32, hdr2.NumFaces * 3)  // vertex 1, 2, 3 for each face

	// read vertices
	if err := binary.Read(r, endian, &surface.Vertices); err != nil {
		fmt.Println("binary.Read failed on mesh vertices array:", err)
		return surface, err
	}

	// read faces
	if err := binary.Read(r, endian, &surface.Faces); err != nil {
		fmt.Println("binary.Read failed on mesh faces array:", err)
		return surface, err
	}

	if Verbosity >= 2 {
		var numToPrint int = 5
		if hdr2.NumVerts >= int32(numToPrint) {
			// print first 5 vertices
			for i := 0; i < numToPrint; i++ {
				for j := 0; j < 3; j++ {
					fmt.Println("surface.Vertices[", i, "][", j, "]:", surface.Vertices[i*3+j])
				}
			}
		}
		// print first 5 faces
		if hdr2.NumFaces >= int32(numToPrint) {
			for i := 0; i < numToPrint; i++ {
				for j := 0; j < 3; j++ {
					fmt.Println("surface.Faces[", i, "][", j, "]:", surface.Faces[i*3+j])
				}
			}
		}
	}

	return surface, nil
}
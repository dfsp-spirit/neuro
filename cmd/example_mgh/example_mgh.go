// Demo application for the neurogo package. Reads a FreeSurfer curv file containing a slice of per-vertex data.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dfsp-spirit/neuro"
)

var (
	mghfile   string
	verbosity *int
	informat  string

	err error
)

func init() {
	flag.StringVar(&mghfile, "mghfile", "brain.mgz", "The input MGH or MGZ file to read.")
	flag.StringVar(&informat, "informat", "auto", "Input file format. One of 'mgh', 'mgz', or 'auto' (default) to auto-detect the format based on the file extension.")
	verbosity = flag.Int("verbosity", 2, "Verbosity level: 0 = silent, 1 = info, 2 = debug.")
}

func main() {

	apptag := "[EX3] "
	//neurogo.Verbosity = *verbosity

	flag.Parse()
	fmt.Println("=====[ Neuro Example 3: Read a FreeSurfer MGH file ]=====")

	if !(informat == "mgh" || informat == "mgz" || informat == "auto") {
		fmt.Printf("%sInvalid value '%s' for parameter 'informat': must be one of 'mgh', 'mgz' or 'auto'.\n", apptag, informat)
		return
	}
	var isGzipped string = informat

	if *verbosity > 0 {
		fmt.Println(apptag, "mghfile:", mghfile)
		fmt.Println(apptag, "informat:", informat)
		fmt.Println(apptag, "verbosity:", *verbosity)
	}

	if _, err := os.Stat(mghfile); err != nil {
		fmt.Printf("%sCould not stat file '%s': '%s', exiting.\n.", apptag, mghfile, err)
		return
	}

	mgh, err := neuro.ReadFsMgh(mghfile, isGzipped)
	if err != nil {
		fmt.Printf("%sFailed to Mgh data from file '%s': '%s'.\n", apptag, mghfile, err)
		return
	}

	h := mgh.Header

	fmt.Printf("%sRead Mgh with data dimensions (%d, %d, %d, %d) from file '%s'.\n", apptag, int(h.Dim1Length), int(h.Dim2Length), int(h.Dim3Length), int(h.Dim4Length), mghfile)

	switch mgh.Data.MghDataType {
	case neuro.MRI_UCHAR:
		fmt.Printf("%sThe first value of the MGH data is %d.\n", apptag, mgh.Data.DataMriUchar[0])
	case neuro.MRI_INT:
		fmt.Printf("%sThe first value of the MGH data is %d.\n", apptag, mgh.Data.DataMriInt[0])
	case neuro.MRI_FLOAT:
		fmt.Printf("%sThe first value of the MGH data is %f.\n", apptag, mgh.Data.DataMriFloat[0])
	case neuro.MRI_SHORT:
		fmt.Printf("%sThe first value of the MGH data is %d.\n", apptag, mgh.Data.DataMriShort[0])
	default:
		fmt.Printf("%sUnknown data type %d in MGH file.\n", apptag, mgh.Data.MghDataType)
	}

}

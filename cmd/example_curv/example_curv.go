// Demo application for the neurogo package. Reads a FreeSurfer curv file containing a slice of per-vertex data.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dfsp-spirit/neurogo"
)

var (
    curvfile string
	exportfile_json string
    verbosity * int

	err error
)

func init() {
    flag.StringVar(&curvfile, "curvfile", "lh.thickness", "The input per-vertex data curv file to read, in binary FreeSurfer curv format.")
	flag.StringVar(&exportfile_json, "exportjson", "", "Output file to which to export per-vertex data in JSON format. Path to it must exist.")
    verbosity = flag.Int("verbosity", 2, "Verbosity level: 0 = silent, 1 = info, 2 = debug.")
}

func main() {

	apptag := "[EX2] "
	//neurogo.Verbosity = *verbosity

    flag.Parse()
	fmt.Println("=====[ Neuro Example 2: Read a FreeSurfer curv file ]=====")

	if *verbosity > 0 {
    	fmt.Println(apptag, "curvfile:", curvfile)
    	fmt.Println(apptag, "verbosity:", *verbosity)
	}

	if _, err := os.Stat(curvfile); err != nil {
		fmt.Printf("%sCould not stat file '%s': '%s', exiting.\n.", apptag, curvfile, err)
		return
	}

	pervertex_data, err := neurogo.ReadFsCurv(curvfile)
	if err != nil {
		fmt.Printf("%sFailed to read per-vertex data from file '%s': '%s'.\n", apptag, curvfile, err)
		return
	}

	if *verbosity > 0 {
    	fmt.Printf("%sRead per-vertex data overlay for %d vertices from meshfile '%s'.\n", apptag, len(pervertex_data), curvfile)
	}


	if len(exportfile_json) > 0 {
		fmt.Printf("%sExporting per-vertex data in format %s to file '%s'.\n", apptag, "JSON", exportfile_json)
		file, err := json.MarshalIndent(pervertex_data, "", " ")
		if err != nil {
			fmt.Printf("%sError converting per-vertex data to JSON: %s", apptag, err)
		}
		err = ioutil.WriteFile(exportfile_json, file, 0644)
	}
	if err != nil {
		fmt.Printf("%sError exporting per-vertex data: %s", apptag, err)
	} else {
		fmt.Printf("%sExported per-vertex data.\n", apptag)
	}


}
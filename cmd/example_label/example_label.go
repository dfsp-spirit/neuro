// Demo application for the neurogo package. Reads a FreeSurfer surface label file containing all vertices of the cortex. Used to mask out medial wall.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dfsp-spirit/neuro"
)

var (
    labelfile string
	exportfile_json string
    verbosity * int

	err error
)

func init() {
    flag.StringVar(&labelfile, "labelfile", "lh.cortex.label", "The input label file to read, in ASCII FreeSurfer label format.")
	flag.StringVar(&exportfile_json, "exportjson", "", "Output file to which to export label data in JSON format. Path to it must exist.")
    verbosity = flag.Int("verbosity", 2, "Verbosity level: 0 = silent, 1 = info, 2 = debug.")
}

func main() {

	apptag := "[EX4] "
	//neurogo.Verbosity = *verbosity

    flag.Parse()
	fmt.Println("=====[ Neuro Example 4: Read a FreeSurfer label file ]=====")

	if *verbosity > 0 {
    	fmt.Println(apptag, "labelfile:", labelfile)
    	fmt.Println(apptag, "verbosity:", *verbosity)
	}

	if _, err := os.Stat(labelfile); err != nil {
		fmt.Printf("%sCould not stat file '%s': '%s', exiting.\n.", apptag, labelfile, err)
		return
	}

	label, err := neuro.ReadFsLabel(labelfile)
	if err != nil {
		fmt.Printf("%sFailed to read label data from file '%s': '%s'.\n", apptag, labelfile, err)
		return
	}

	if *verbosity > 0 {
    	fmt.Printf("%sRead label containing %d vertices from file '%s'.\n", apptag, len(label.Value), labelfile)
	}


	if len(exportfile_json) > 0 {
		fmt.Printf("%sExporting label data in format %s to file '%s'.\n", apptag, "JSON", exportfile_json)
		file, err := json.MarshalIndent(label, "", " ")
		if err != nil {
			fmt.Printf("%sError converting label data to JSON: %s", apptag, err)
		}
		err = ioutil.WriteFile(exportfile_json, file, 0644)
		if err != nil {
			fmt.Printf("%sError exporting label data: %s", apptag, err)
		} else {
			fmt.Printf("%sExported label data to file '%s'.\n", apptag, exportfile_json)
		}
	}



}
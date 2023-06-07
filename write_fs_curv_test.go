package neuro

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteRereadCurv(t *testing.T){

	data := []float32{1.0, 2.0, 3.0, 4.0, 5.0}

	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	curv_file_name := file.Name()
	file.Close()

    err = WriteFsCurv(curv_file_name, data)
	if err != nil {
		t.Errorf("WriteFsCurv failed: %v", err)
	}
	data_reread, err := ReadFsCurv(curv_file_name)
	if err != nil {
		t.Errorf("ReadFsCurv failed: %v", err)
	}


    if diff := cmp.Diff(data, data_reread); diff != "" {
        t.Error(diff)
    }
}
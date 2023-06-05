package neurogo

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteRereadCurv(t *testing.T){

	data := []float32{1.0, 2.0, 3.0, 4.0, 5.0}

    err := WriteFsCurv("test.curv", data)
	if err != nil {
		t.Errorf("WriteFsCurv failed: %v", err)
	}
	data_reread, err := ReadFsCurv("test.curv")
	if err != nil {
		t.Errorf("ReadFsCurv failed: %v", err)
	}


    if diff := cmp.Diff(data, data_reread); diff != "" {
        t.Error(diff)
    }
}
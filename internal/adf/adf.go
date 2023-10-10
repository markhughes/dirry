package adf

import (
	"fmt"

	"github.com/markhughes/dirry/internal/libadf"
	"github.com/markhughes/dirry/internal/libmrf"
)

func Dump(filePath string) {
	// note: this isn't working yet.. for now users should use a manual way to dump the mrf

	var adf, err = libadf.UnpackAdfFromFile(filePath)
	if err != nil {

		fmt.Printf("Error: %v", err)
	}

	for k, v := range adf {
		byt, err := libmrf.FromBytes(v)
		if err != nil {
			fmt.Printf("Error: %v", err)
		} else {
			fmt.Printf("%v\n", k)
			fmt.Printf("bytes: %v", byt)

		}
	}

}

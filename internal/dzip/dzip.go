package dzip

import (
	"fmt"
	"path/filepath"

	"github.com/markhughes/dirry/internal/shockwave"
	"github.com/markhughes/dirry/internal/utils"
)

func DZip(filePath string, pkg string) error {
	var err error

	var shockwave shockwave.Shockwave
	shockwave.PkgName = pkg
	expanded, err := shockwave.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", filePath, err)
		return fmt.Errorf("error opening file: %s", err)
	}

	if len(expanded) > 0 {
		fmt.Printf("Expanded %s to %d files\n", filePath, len(expanded))

		for i := range expanded {
			fmt.Printf("Dumping %s\n", expanded[i])
			var err = DZip(expanded[i].Path, filepath.Base(filePath))
			if err != nil {
				fmt.Printf("Error dumping file %s: %s\n", expanded[i], err)
				return fmt.Errorf("error dumping file: %s", err)
			}
		}

		return nil
	} else {
		utils.InfoMsg("dzip", "Creating zip file for %s\n", filePath)
		return shockwave.Zip()
	}
}

package consts

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// enum style  const of output paths

var PathZips = ""
var PathDump = ""
var LogsDir = ""
var PalettesDir = ""
var PatternsDir = ""

func GetDirryRoot() (string, error) {
	if strings.Contains(os.Args[0], os.TempDir()) {
		// Running through `go run`
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return dir, nil
	} else {
		// Running as a built binary
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "dirry"), nil
	}
}

func init() {
	dirryRoot, err := GetDirryRoot()
	if err != nil {
		fmt.Println("NOTE: could not get dirry root")
	} else {

		fmt.Printf("dirryRoot: %s\n", dirryRoot)
		os.MkdirAll(path.Join(dirryRoot), os.ModePerm)

		PathZips = path.Join(dirryRoot, "out", "zips")
		PathDump = path.Join(dirryRoot, "out", "dump")
		LogsDir = path.Join(dirryRoot, "logs")
		PalettesDir = path.Join(dirryRoot, "resources", "palettes")
		PatternsDir = path.Join(dirryRoot, "resources", "patterns")
	}
}

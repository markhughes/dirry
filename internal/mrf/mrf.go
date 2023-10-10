package mrf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/markhughes/dirry/internal/consts"
	"github.com/markhughes/dirry/internal/libmrf"
)

func Dump(filePath string) {

	var resourceFork, err = libmrf.FromFile(filePath)
	if err != nil {
		panic(err)
	}

	var codeResources []libmrf.Resource

	for _, resource := range resourceFork.Resources {
		outputFolder := filepath.Join(consts.PathDump, filepath.Base(filePath), "mrf", "binary", resource.Type)
		os.MkdirAll(outputFolder, os.ModePerm)

		// dump data into a file dump/<file name>/<resource type>/<resource name>
		var fileName = filepath.Join(outputFolder, resource.Name)
		var file, err = os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = file.Write(resource.Data)
		if err != nil {
			panic(err)
		}

		if resource.Type == "CODE" || resource.Type == "CSND" {
			codeResources = append(codeResources, resource)
		}

		if resource.Type == "TEXT" {
			var text = string(resource.Data)
			outputFolder := filepath.Join(consts.PathDump, filepath.Base(filePath), "mrf", "text")
			os.MkdirAll(outputFolder, os.ModePerm)

			var fileName = filepath.Join(outputFolder, resource.Name+".txt")
			var file, err = os.Create(fileName)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			fmt.Fprintf(file, "%s", text)

		}

	}

}

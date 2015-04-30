package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	var name, src, dest string

	src = os.Args[1]

	if len(os.Args) > 2 {
		dest = os.Args[2]
	} else {
		dest = src
	}

	files, err := ioutil.ReadDir(src)

	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range files {
		name = fi.Name()

		if !fi.IsDir() && path.Ext(name) == ".csv" {
			os.Rename(path.Join(src, name), path.Join(dest, strings.Replace(name, "-Table 1", "", -1)))
		}
	}
}

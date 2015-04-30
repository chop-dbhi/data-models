package main

import (
	"io/ioutil"
	"os"
)

func main() {
	data, _ := ioutil.ReadAll(os.Stdin)

	for i, b := range data {
		if b == '\r' {
			data[i] = '\n'
		}
	}

	os.Stdout.Write(data)
}

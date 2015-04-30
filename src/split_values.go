package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

var choiceRe = regexp.MustCompile(`([A-Z0-9_]+)=([^=]+)+`)

func main() {
	data, _ := ioutil.ReadAll(os.Stdin)

	var (
		// equals has been seen
		equals bool
		// index of the last space seen
		space int

		offset int

		choices = make([][]byte, 0)
	)

	for {
	NEXT:
		space = 0
		equals = false
		choice := make([]byte, 0)

		for i, b := range data[offset:] {
			// Keep track of the spaces
			if b == ' ' || b == '\r' || b == '\n' {
				b = ' '
				space = i

				// Equals has been seen
			} else if b == '=' {
				// Equals has already been seen, then is the next
				// choice. Trim choice to the last space.
				if equals {
					choice = choice[:space]
					choices = append(choices, choice)
					offset += space + 1
					goto NEXT
				}

				// Flag as seen
				equals = true
			}

			choice = append(choice, b)
		}

		offset += len(choice)

		// Ensure we do not have a newline
		if len(bytes.Trim(choice, " ")) > 0 {
			choices = append(choices, choice)
		}

		if offset >= len(data) {
			break
		}
	}

	for _, c := range choices {
		b := bytes.Replace(c, []byte{'='}, []byte{' ', '=', ' '}, -1)
		fmt.Println(string(b))
	}
}

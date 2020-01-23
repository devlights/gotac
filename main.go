package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	filePath string
)

// REF: https://github.com/coreutils/coreutils/blob/master/src/tac.c
func main() {
	flag.Parse()

	logErr := func(e error) {
		if e != nil {
			log.Fatal(e)
		}
	}

	args := flag.Args()
	if len(flag.Args()) == 0 {
		args = append(args, "stdin")
	}

	for _, filePath := range args {
		var f *os.File
		if filePath == "stdin" {
			f = os.Stdin
			tempfile, err := ioutil.TempFile("", "gotac")
			logErr(err)

			defer os.Remove(tempfile.Name())

			d, err := ioutil.ReadAll(f)
			logErr(err)

			ioutil.WriteFile(tempfile.Name(), d, 0644)

			filePath = tempfile.Name()
		}

		f, err := os.Open(filePath)
		logErr(err)
		defer f.Close()

		offset, err := f.Seek(-1, os.SEEK_END)
		logErr(err)

		firstTime := true
		data := make([]byte, 0, 1024)
		for {
			buf := make([]byte, 1)

			_, err := f.ReadAt(buf, offset)
			if err != nil {
				data = data[:len(data)]
				for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
					data[i], data[j] = data[j], data[i]
				}

				fmt.Println(string(data))
				break
			}

			b := buf[0]
			if b == '\n' {
				if firstTime {
					// drop first newline
					firstTime = false
					offset--
					continue
				}

				data = data[:len(data)]
				for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
					data[i], data[j] = data[j], data[i]
				}

				fmt.Println(string(data))
				data = make([]byte, 0, 1024)
			} else {
				data = append(data, b)
				firstTime = false
			}

			offset--
		}
	}
}

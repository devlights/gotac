package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func pl(w io.Writer, s string) {
	_, err := fmt.Fprintln(w, s)
	logErr(err, false)
}

func safeClose(c io.Closer) {
	err := c.Close()
	logErr(err, false)
}

func safeRemove(f string) {
	err := os.Remove(f)
	logErr(err, false)
}

func logErr(err error, fatal bool) {
	if err != nil {
		if fatal {
			log.Fatal(err)
		} else {
			log.Println(err)
		}
	}
}

func copyToTmpFile(src io.Reader) (string, error) {
	tempFile, err := ioutil.TempFile("", "gotac")
	if err != nil {
		return "", err
	}

	defer safeClose(tempFile)

	d, err := ioutil.ReadAll(src)
	if err != nil {
		return "", err
	}

	name := tempFile.Name()
	err = ioutil.WriteFile(name, d, 0644)
	if err != nil {
		return name, err
	}

	return name, nil
}

package main

import (
	"log"
	"os"
)

// REF: https://github.com/coreutils/coreutils/blob/master/src/tac.c
func main() {
	args := NewArgs()

	err := args.Parse()
	if err != nil {
		log.Fatal(err)
	}

	defer args.RemoveTempFileAtExit()

	for _, filePath := range args.TargetFiles {
		command := NewTacCommand()
		commandArgs := NewTacArgs(filePath, os.Stdout)

		err := command.Run(commandArgs)
		logErr(err, true)
	}
}

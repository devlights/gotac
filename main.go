package main

import (
	"os"
)

// REF: https://github.com/coreutils/coreutils/blob/master/src/tac.c
func main() {
	args := NewArgs()

	err := args.Parse()
	logErr(err, true)

	defer args.RemoveTempFile()

	for _, filePath := range args.TargetFiles {
		command := NewTacCommand()
		commandArgs := NewTacArgs(filePath, os.Stdout)

		err := command.Run(commandArgs)
		logErr(err, true)
	}
}

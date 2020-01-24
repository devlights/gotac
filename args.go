package main

import (
	"flag"
	"os"
)

type Args struct {
	TargetFiles []string
	IsStdin     bool
}

func NewArgs() *Args {
	args := new(Args)
	return args
}

func (a *Args) Parse() error {
	flag.Parse()

	a.TargetFiles = flag.Args()
	if len(a.TargetFiles) == 0 {
		// ファイルが一つも指定されていない場合は標準入力から読み取る
		// 標準入力のままだと seek 出来ないため、一時ファイルにコピー
		a.IsStdin = true

		fname, err := copyToTmpFile(os.Stdin)
		if err != nil {
			return err
		}

		a.TargetFiles = append(a.TargetFiles, fname)
	}

	return nil
}

func (a *Args) RemoveTempFile() {
	if a.IsStdin {
		if len(a.TargetFiles) > 0 {
			f := a.TargetFiles[0]
			safeRemove(f)
		}
	}
}

package main

import (
	"io"
	"os"
)

const (
	BufferSize    = 1024 * 8
	InitialOffset = -1
)

type (
	TacArgs struct {
		TargetFile string
		Writer     io.Writer
	}

	TacCommand struct {
		offset    int64
		firstTime bool
	}
)

func NewTacCommand() *TacCommand {
	c := new(TacCommand)
	c.firstTime = true
	return c
}

func NewTacArgs(file string, writer io.Writer) *TacArgs {
	a := new(TacArgs)
	a.TargetFile = file
	a.Writer = writer

	return a
}

//noinspection GoRedundantSecondIndexInSlices
func (t *TacCommand) Run(args *TacArgs) error {

	f, err := os.Open(args.TargetFile)
	if err != nil {
		return err
	}

	defer safeClose(f)

	o, err := f.Seek(InitialOffset, io.SeekEnd)
	if err != nil {
		return err
	}

	t.offset = o

	data := t.newBuffer()
	for {
		buf := make([]byte, 1)

		_, err := f.ReadAt(buf, t.offset)
		if err != nil {
			data = data[:len(data)]
			t.reverse(data)

			pl(os.Stdout, string(data))
			break
		}

		b := buf[0]

		if b == '\n' {
			if t.firstTime {
				// drop first newline
				t.firstTime = false
				t.backOffset()
				continue
			}

			data = data[:len(data)]
			t.reverse(data)

			pl(os.Stdout, string(data))
			data = t.newBuffer()
		} else {
			data = append(data, b)
			t.firstTime = false
		}

		t.backOffset()
	}

	return nil
}

func (t *TacCommand) newBuffer() []byte {
	return make([]byte, 0, BufferSize)
}

func (t *TacCommand) reverse(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func (t *TacCommand) backOffset() {
	t.offset--
}

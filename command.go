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
		buffer    []byte
	}
)

func NewTacCommand() *TacCommand {
	c := new(TacCommand)
	c.setFirstTime(true)
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

	t.newBuffer()
	for {
		buf := make([]byte, 1)

		_, err := f.ReadAt(buf, t.offset)
		if err != nil {
			t.adjustBuffer()
			t.reverse()

			pl(os.Stdout, string(t.buffer))
			break
		}

		b := buf[0]

		if b == '\n' {
			if t.firstTime {
				// drop first newline
				t.setFirstTime(false)
				t.backOffset()
				continue
			}

			t.adjustBuffer()
			t.reverse()

			pl(os.Stdout, string(t.buffer))
			t.newBuffer()
		} else {
			t.addData(b)
			t.setFirstTime(false)
		}

		t.backOffset()
	}

	return nil
}

func (t *TacCommand) setFirstTime(isFirstTime bool) {
	t.firstTime = isFirstTime
}

func (t *TacCommand) newBuffer() {
	t.buffer = make([]byte, 0, BufferSize)
}

func (t *TacCommand) adjustBuffer() {
	t.buffer = t.buffer[:len(t.buffer)]
}

func (t *TacCommand) addData(b byte) {
	t.buffer = append(t.buffer, b)
}

func (t *TacCommand) reverse() {
	for i, j := 0, len(t.buffer)-1; i < j; i, j = i+1, j-1 {
		t.buffer[i], t.buffer[j] = t.buffer[j], t.buffer[i]
	}
}

func (t *TacCommand) backOffset() {
	t.offset--
}

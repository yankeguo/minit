package mlog

import (
	"fmt"
	"github.com/guoyk93/gg"
	"io"
	"os"
)

type ProcLoggerOptions struct {
	*RotatingFileOptions

	ConsoleOut io.Writer
	ConsoleErr io.Writer

	ConsolePrefix string
	FilePrefix    string
}

type ProcLogger interface {
	Print(items ...interface{})
	Printf(layout string, items ...interface{})
	Error(items ...interface{})
	Errorf(layout string, items ...interface{})

	ProcOutput
}

type procLogger struct {
	out Output
	err Output
}

func NewProcLogger(opts ProcLoggerOptions) (pl ProcLogger, err error) {
	if opts.ConsoleOut == nil {
		opts.ConsoleOut = os.Stdout
	}
	if opts.ConsoleErr == nil {
		opts.ConsoleErr = os.Stderr
	}
	if opts.MaxFileSize == 0 {
		opts.MaxFileSize = 128 * 1024 * 1024
	}
	if opts.MaxFileCount == 0 {
		opts.MaxFileCount = 5
	}

	if opts.RotatingFileOptions == nil {
		pl = &procLogger{
			out: NewWriterOutput(opts.ConsoleOut, []byte(opts.ConsolePrefix), nil),
			err: NewWriterOutput(opts.ConsoleErr, []byte(opts.ConsolePrefix), nil),
		}
		return
	}

	var fileOut io.WriteCloser
	if fileOut, err = NewRotatingFile(RotatingFileOptions{
		Dir:          opts.Dir,
		Filename:     opts.Filename + ".out",
		MaxFileSize:  opts.MaxFileSize,
		MaxFileCount: opts.MaxFileCount,
	}); err != nil {
		return
	}

	var fileErr io.WriteCloser
	if fileErr, err = NewRotatingFile(RotatingFileOptions{
		Dir:          opts.Dir,
		Filename:     opts.Filename + ".err",
		MaxFileSize:  opts.MaxFileSize,
		MaxFileCount: opts.MaxFileCount,
	}); err != nil {
		return
	}

	pl = &procLogger{
		out: MultiOutput(
			NewWriterOutput(fileOut, []byte(opts.FilePrefix), nil),
			NewWriterOutput(opts.ConsoleOut, []byte(opts.ConsolePrefix), nil),
		),
		err: MultiOutput(
			NewWriterOutput(fileErr, []byte(opts.FilePrefix), nil),
			NewWriterOutput(opts.ConsoleErr, []byte(opts.ConsolePrefix), nil),
		),
	}
	return
}

func (pl *procLogger) Close() error {
	eg := gg.NewErrorGroup()
	eg.Add(pl.out.Close())
	eg.Add(pl.err.Close())
	return eg.Unwrap()
}

func (pl *procLogger) Print(items ...interface{}) {
	_, _ = pl.out.Write(append([]byte(fmt.Sprint(items...)), '\n'))
}

func (pl *procLogger) Error(items ...interface{}) {
	_, _ = pl.err.Write(append([]byte(fmt.Sprint(items...)), '\n'))
}

func (pl *procLogger) Printf(pattern string, items ...interface{}) {
	_, _ = pl.out.Write(append([]byte(fmt.Sprintf(pattern, items...)), '\n'))
}

func (pl *procLogger) Errorf(pattern string, items ...interface{}) {
	_, _ = pl.err.Write(append([]byte(fmt.Sprintf(pattern, items...)), '\n'))
}

func (pl *procLogger) Out() Output {
	return pl.out
}

func (pl *procLogger) Err() Output {
	return pl.err
}

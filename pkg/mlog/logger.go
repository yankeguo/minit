package mlog

import (
	"fmt"
	"github.com/guoyk93/minit/pkg/merrs"
	"io"
	"os"
)

type ProcLoggerOptions struct {
	ConsoleOut    io.Writer
	ConsoleErr    io.Writer
	ConsolePrefix string

	FilePrefix  string
	FileOptions *RotatingFileOptions
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

	if opts.FileOptions == nil {
		pl = &procLogger{
			out: NewWriterOutput(opts.ConsoleOut, []byte(opts.ConsolePrefix), nil),
			err: NewWriterOutput(opts.ConsoleErr, []byte(opts.ConsolePrefix), nil),
		}
		return
	}

	if opts.FileOptions.MaxFileSize == 0 {
		opts.FileOptions.MaxFileSize = 128 * 1024 * 1024
	}
	if opts.FileOptions.MaxFileCount == 0 {
		opts.FileOptions.MaxFileCount = 5
	}

	var fileOut io.WriteCloser
	if fileOut, err = NewRotatingFile(RotatingFileOptions{
		Dir:          opts.FileOptions.Dir,
		Filename:     opts.FileOptions.Filename + ".out",
		MaxFileSize:  opts.FileOptions.MaxFileSize,
		MaxFileCount: opts.FileOptions.MaxFileCount,
	}); err != nil {
		return
	}

	var fileErr io.WriteCloser
	if fileErr, err = NewRotatingFile(RotatingFileOptions{
		Dir:          opts.FileOptions.Dir,
		Filename:     opts.FileOptions.Filename + ".err",
		MaxFileSize:  opts.FileOptions.MaxFileSize,
		MaxFileCount: opts.FileOptions.MaxFileCount,
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
	eg := merrs.NewErrorGroup()
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

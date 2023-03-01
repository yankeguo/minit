package mlog

import (
	"bufio"
	"github.com/guoyk93/minit/pkg/merrs"
	"io"
	"sync"
)

// Output interface for single stream log output
type Output interface {
	// WriteCloser is for single line writing
	io.WriteCloser

	// ReaderFrom is for streaming
	io.ReaderFrom
}

// ProcOutput interface for process
type ProcOutput interface {
	// Out stdout
	Out() Output
	// Err stderr
	Err() Output
}

type writerOutput struct {
	pfx []byte
	sfx []byte
	w   io.Writer
}

func (w *writerOutput) Write(p []byte) (n int, err error) {
	if len(w.pfx) == 0 && len(w.sfx) == 0 {
		n, err = w.w.Write(p)
		return
	}
	if n, err = w.w.Write(
		append(
			append(w.pfx, p...),
			w.sfx...,
		),
	); err != nil {
		return
	}

	n = len(p)

	return
}

func (w *writerOutput) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func (w *writerOutput) ReadFrom(r io.Reader) (n int64, err error) {
	br := bufio.NewReader(r)
	for {
		var line []byte
		if line, err = br.ReadBytes('\n'); err == nil {
			_, _ = w.Write(line)
			n += int64(len(line))
		} else {
			if err == io.EOF {
				err = nil
			}
			if len(line) != 0 {
				_, _ = w.Write(append(line, '\n'))
				n += int64(len(line))
			}
			break
		}
	}
	return
}

// NewWriterOutput wrap a writer as a Output, with optional line Prefix and Suffix
func NewWriterOutput(w io.Writer, pfx, sfx []byte) Output {
	return &writerOutput{w: w, pfx: pfx, sfx: sfx}
}

type multiOutput struct {
	outputs []Output
}

// MultiOutput create a new Output for proc logging
func MultiOutput(outputs ...Output) Output {
	return &multiOutput{outputs: outputs}
}

func (pc *multiOutput) Close() error {
	eg := merrs.NewErrorGroup()
	for _, output := range pc.outputs {
		eg.Add(output.Close())
	}
	return eg.Unwrap()
}

// Write this method is used to write a single line of log
func (pc *multiOutput) Write(buf []byte) (n int, err error) {
	for _, output := range pc.outputs {
		if n, err = output.Write(buf); err != nil {
			return
		}
	}
	n = len(buf)
	return
}

// ReadFrom implements ReaderFrom
func (pc *multiOutput) ReadFrom(r io.Reader) (n int64, err error) {
	eg := merrs.NewErrorGroup()
	wg := &sync.WaitGroup{}

	var (
		cs []io.Closer
		ws []io.Writer
	)

	for _, _out := range pc.outputs {
		out := _out

		childR, childW := io.Pipe()
		cs, ws = append(cs, childW), append(ws, childW)

		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := out.ReadFrom(childR)
			if err == io.EOF {
				err = nil
			}
			eg.Add(err)
		}()
	}

	_, err = io.Copy(io.MultiWriter(ws...), r)
	if err == io.EOF {
		err = nil
	}
	for _, c := range cs {
		_ = c.Close()
	}

	wg.Wait()

	if err == nil {
		err = eg.Unwrap()
	}
	return
}

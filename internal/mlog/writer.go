package mlog

import (
	"bytes"
	"io"
	"log"
	"sync"
)

type loggerWriter struct {
	logger *log.Logger
	buf    *bytes.Buffer
	pfx    string

	lock sync.Locker
}

// NewLoggerWriter create a new io.WriteCloser that append each line to log.procLogger
func NewLoggerWriter(logger *log.Logger, prefix string) io.WriteCloser {
	return &loggerWriter{
		logger: logger,
		buf:    &bytes.Buffer{},
		pfx:    prefix,

		lock: &sync.Mutex{},
	}
}

func (w *loggerWriter) finish(force bool) (err error) {
	var line string

	for {
		// read till new line
		if line, err = w.buf.ReadString('\n'); err == nil {
			// output
			if err = w.logger.Output(3, w.pfx+line); err != nil {
				return
			}
		} else {
			if force {
				// if forced, output to logger
				if err = w.logger.Output(3, w.pfx+line); err != nil {
					return
				}
			} else {
				// write back
				if _, err = w.buf.WriteString(line); err != nil {
					return
				}
			}
			break
		}
	}

	return
}

func (w *loggerWriter) Close() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if err = w.finish(true); err != nil {
		return
	}
	return
}

func (w *loggerWriter) Write(p []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if n, err = w.buf.Write(p); err != nil {
		return
	}

	if err = w.finish(false); err != nil {
		return
	}
	return
}

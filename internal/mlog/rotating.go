package mlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type rotatingFile struct {
	opts RotatingFileOptions

	fd   *os.File
	size int64
	lock sync.Locker
}

// RotatingFileOptions options for creating a RotatingFile
type RotatingFileOptions struct {
	// Dir directory
	Dir string
	// Filename filename prefix
	Filename string
	// MaxFileSize max size of a single file, default to 128mb
	MaxFileSize int64
	// MaxFileCount max count of rotated files
	MaxFileCount int64
}

// NewRotatingFile create a new io.WriteCloser as a rotating log file
func NewRotatingFile(opts RotatingFileOptions) (w io.WriteCloser, err error) {
	if opts.MaxFileSize == 0 {
		opts.MaxFileSize = 128 * 1000 * 1000
	}
	rf := &rotatingFile{opts: opts, lock: &sync.Mutex{}}
	if err = rf.open(); err != nil {
		return
	}
	w = rf
	return
}

func (rf *rotatingFile) currentPath() string {
	return filepath.Join(rf.opts.Dir, rf.opts.Filename+".log")
}

func (rf *rotatingFile) rotatedPath(id int64) string {
	return filepath.Join(rf.opts.Dir, fmt.Sprintf("%s.%d.log", rf.opts.Filename, id))
}

func (rf *rotatingFile) nextRotatedID() (id int64, err error) {
	var entries []os.DirEntry
	if entries, err = os.ReadDir(rf.opts.Dir); err != nil {
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, rf.opts.Filename+".") &&
			strings.HasSuffix(name, ".log") {
			eIDStr := strings.TrimSuffix(strings.TrimPrefix(name, rf.opts.Filename+"."), ".log")
			eID, _ := strconv.ParseInt(eIDStr, 10, 64)
			if eID > id {
				id = eID
			}
		}
	}

	id += 1

	// if id exceeded MaxFileCount, back to 1
	if rf.opts.MaxFileCount > 0 && id > rf.opts.MaxFileCount {
		id = 1
	}
	return
}

func (rf *rotatingFile) open() (err error) {
	var fd *os.File
	if fd, err = os.OpenFile(
		rf.currentPath(),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0644,
	); err != nil {
		return
	}

	var info os.FileInfo
	if info, err = fd.Stat(); err != nil {
		_ = fd.Close()
		return
	}

	existed := rf.fd

	rf.fd = fd
	rf.size = info.Size()

	if existed != nil {
		_ = existed.Close()
	}

	return
}

func (rf *rotatingFile) reallocate() (err error) {
	rf.lock.Lock()
	defer rf.lock.Unlock()

	// recheck, in case of race condition
	if atomic.LoadInt64(&rf.size) <= rf.opts.MaxFileSize {
		return
	}

	// find next rotated id
	var id int64
	if id, err = rf.nextRotatedID(); err != nil {
		return
	}

	// try remove existed, in case id looped due to maxCount
	_ = os.Remove(rf.rotatedPath(id))

	// remove current file to rotated path
	if err = os.Rename(rf.currentPath(), rf.rotatedPath(id)); err != nil {
		return
	}

	// open current file, this will close existing file
	if err = rf.open(); err != nil {
		return
	}

	return nil
}

func (rf *rotatingFile) Write(p []byte) (n int, err error) {
	if n, err = rf.fd.Write(p); err != nil {
		return
	}

	// reallocate if exceeded
	if atomic.AddInt64(&rf.size, int64(n)) > rf.opts.MaxFileSize {
		if err = rf.reallocate(); err != nil {
			return
		}
	}

	return
}

func (rf *rotatingFile) Close() (err error) {
	rf.lock.Lock()
	defer rf.lock.Unlock()

	if rf.fd != nil {
		err = rf.fd.Close()
		rf.fd = nil
	}
	return
}

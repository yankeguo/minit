package msetups

import (
	"bytes"
	"github.com/guoyk93/minit/pkg/mlog"
	"os"
)

const (
	BannerFile = "/etc/banner.minit.txt"
)

func init() {
	Register(10, setupBanner)
}

func setupBanner(logger mlog.ProcLogger) (err error) {
	var buf []byte
	if buf, err = os.ReadFile(BannerFile); err != nil {
		err = nil
		return
	}

	lines := bytes.Split(buf, []byte{'\n'})
	for _, line := range lines {
		logger.Print(string(line))
	}

	return
}

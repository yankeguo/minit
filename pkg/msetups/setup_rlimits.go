//go:build linux

package msetups

import (
	"fmt"
	"github.com/guoyk93/minit/pkg/mlog"
	"golang.org/x/sys/unix"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const Unlimited = "unlimited"

var (
	knownRLimitNames = map[string]int{
		"AS":         unix.RLIMIT_AS,
		"CORE":       unix.RLIMIT_CORE,
		"CPU":        unix.RLIMIT_CPU,
		"DATA":       unix.RLIMIT_DATA,
		"FSIZE":      unix.RLIMIT_FSIZE,
		"LOCKS":      unix.RLIMIT_LOCKS,
		"MEMLOCK":    unix.RLIMIT_MEMLOCK,
		"MSGQUEUE":   unix.RLIMIT_MSGQUEUE,
		"NICE":       unix.RLIMIT_NICE,
		"NOFILE":     unix.RLIMIT_NOFILE,
		"NPROC":      unix.RLIMIT_NPROC,
		"RTPRIO":     unix.RLIMIT_RTPRIO,
		"SIGPENDING": unix.RLIMIT_SIGPENDING,
		"STACK":      unix.RLIMIT_STACK,
	}
)

func decodeRLimitValue(v *uint64, s string) (err error) {
	s = strings.TrimSpace(s)
	if s == "-" || s == "" {
		return
	}
	if strings.ToLower(s) == Unlimited {
		*v = unix.RLIM_INFINITY
	} else {
		if *v, err = strconv.ParseUint(s, 10, 64); err != nil {
			return
		}
	}
	return
}

func formatRLimitValue(v uint64) string {
	if v == unix.RLIM_INFINITY {
		return Unlimited
	} else {
		return strconv.FormatUint(v, 10)
	}
}

func init() {
	Register(30, setupRLimits)
}

func setupRLimits(logger mlog.ProcLogger) (err error) {
	for name, res := range knownRLimitNames {
		key := "MINIT_RLIMIT_" + name
		val := strings.TrimSpace(os.Getenv(key))
		if val == "-" || val == "-:-" || val == "" {
			continue
		}
		var limit syscall.Rlimit
		if err = syscall.Getrlimit(res, &limit); err != nil {
			err = fmt.Errorf("failed getting rlimit_%s: %s", name, err.Error())
			return
		}
		logger.Printf("current rlimit_%s=%s:%s", name, formatRLimitValue(limit.Cur), formatRLimitValue(limit.Max))
		if strings.Contains(val, ":") {
			splits := strings.Split(val, ":")
			if len(splits) != 2 {
				err = fmt.Errorf("invalid environment variable %s=%s", key, val)
				return
			}
			if err = decodeRLimitValue(&limit.Cur, splits[0]); err != nil {
				err = fmt.Errorf("invalid environment variable %s=%s: %s", key, val, err.Error())
				return
			}
			if err = decodeRLimitValue(&limit.Max, splits[1]); err != nil {
				err = fmt.Errorf("invalid environment variable %s=%s: %s", key, val, err.Error())
				return
			}
		} else {
			if err = decodeRLimitValue(&limit.Cur, val); err != nil {
				return
			}
			limit.Max = limit.Cur
		}
		logger.Printf("setting rlimit_%s=%s:%s", name, formatRLimitValue(limit.Cur), formatRLimitValue(limit.Max))
		if err = syscall.Setrlimit(res, &limit); err != nil {
			err = fmt.Errorf("failed setting rlimit_%s=%s: %s", name, val, err.Error())
			return
		}
	}

	return
}

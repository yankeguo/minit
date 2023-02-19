package ggos

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	hostnameSequenceIDRegexp = regexp.MustCompile(`[0-9]+$`)
)

// HostnameSequenceID extract a sequence id from hostname
func HostnameSequenceID() (id int, err error) {
	var hostname string
	if hostname = strings.TrimSpace(os.Getenv("HOSTNAME")); hostname == "" {
		if hostname, err = os.Hostname(); err != nil {
			return
		}
	}

	if match := hostnameSequenceIDRegexp.FindStringSubmatch(hostname); len(match) == 0 {
		err = errors.New("no sequence id in hostname")
	} else {
		id, err = strconv.Atoi(match[0])
	}

	return
}

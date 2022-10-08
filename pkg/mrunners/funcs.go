package mrunners

import (
	"errors"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"
)

// Funcs provided funcs for gracetmpl
var Funcs = map[string]interface{}{
	"netResolveIPAddr":    net.ResolveIPAddr,
	"netResolveIP":        netResolveIP,
	"osHostname":          os.Hostname,
	"osUserCacheDir":      os.UserCacheDir,
	"osUserConfigDir":     os.UserConfigDir,
	"osUserHomeDir":       os.UserHomeDir,
	"osGetegid":           os.Getegid,
	"osGetenv":            os.Getenv,
	"osGeteuid":           os.Geteuid,
	"osGetgid":            os.Getgid,
	"osGetgroups":         os.Getgroups,
	"osGetpagesize":       os.Getpagesize,
	"osGetpid":            os.Getpid,
	"osGetppid":           os.Getppid,
	"osGetuid":            os.Getuid,
	"osGetwd":             os.Getwd,
	"osTempDir":           os.TempDir,
	"osUserLookupGroup":   user.LookupGroup,
	"osUserLookupGroupId": user.LookupGroupId,
	"osUserCurrent":       user.Current,
	"osUserLookup":        user.Lookup,
	"osUserLookupId":      user.LookupId,
	"stringsContains":     strings.Contains,
	"stringsFields":       strings.Fields,
	"stringsIndex":        strings.Index,
	"stringsLastIndex":    strings.LastIndex,
	"stringsHasPrefix":    strings.HasPrefix,
	"stringsHasSuffix":    strings.HasSuffix,
	"stringsRepeat":       strings.Repeat,
	"stringsReplaceAll":   strings.ReplaceAll,
	"stringsSplit":        strings.Split,
	"stringsSplitN":       strings.SplitN,
	"stringsToLower":      strings.ToLower,
	"stringsToUpper":      strings.ToUpper,
	"stringsTrimPrefix":   strings.TrimPrefix,
	"stringsTrimSpace":    strings.TrimSpace,
	"stringsTrimSuffix":   strings.TrimSuffix,
	"strconvQuote":        strconv.Quote,
	"strconvUnquote":      strconv.Unquote,
	"strconvParseBool":    strconv.ParseBool,
	"strconvParseInt":     strconv.ParseInt,
	"strconvParseUint":    strconv.ParseUint,
	"strconvParseFloat":   strconv.ParseFloat,
	"strconvFormatBool":   strconv.FormatBool,
	"strconvFormatInt":    strconv.FormatInt,
	"strconvFormatUint":   strconv.FormatUint,
	"strconvFormatFloat":  strconv.FormatFloat,
	"strconvAoti":         strconv.Atoi,
	"strconvItoa":         strconv.Itoa,

	"add":        add,
	"neg":        neg,
	"intAdd":     add,
	"intNeg":     neg,
	"int64Add":   add,
	"int64Neg":   neg,
	"float32Add": add,
	"float32Neg": neg,
	"float64Add": add,
	"float64Neg": neg,

	"osHostnameSequenceID": osHostnameSequenceID,
	"k8sStatefulSetID":     osHostnameSequenceID,
}

func netResolveIP(s string) (ip string, err error) {
	var addr *net.IPAddr
	if addr, err = net.ResolveIPAddr("ip", s); err != nil {
		return
	}
	ip = addr.IP.String()
	return
}

func add(a, b interface{}) interface{} {
	switch a.(type) {
	case bool:
		return a.(bool) || b.(bool)
	case int:
		return a.(int) + b.(int)
	case int64:
		return a.(int64) + b.(int64)
	case int32:
		return a.(int32) + b.(int32)
	case float32:
		return a.(float32) + b.(float32)
	case float64:
		return a.(float64) + b.(float64)
	case string:
		return a.(string) + b.(string)
	}
	return nil
}

func neg(a interface{}) interface{} {
	switch a.(type) {
	case bool:
		return !a.(bool)
	case int:
		return -a.(int)
	case int64:
		return -a.(int64)
	case int32:
		return -a.(int32)
	case float32:
		return -a.(float32)
	case float64:
		return -a.(float64)
	}
	return nil
}

func osHostnameSequenceID() (id int, err error) {
	var hostname string
	if hostname = os.Getenv("HOSTNAME"); hostname == "" {
		if hostname, err = os.Hostname(); err != nil {
			return
		}
	}
	splits := strings.Split(hostname, "-")
	if len(splits) < 2 {
		err = errors.New("invalid stateful-set hostname")
		return
	}
	id, err = strconv.Atoi(splits[len(splits)-1])
	return
}

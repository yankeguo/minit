package mtmpl

import (
	"errors"
	"net"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"strings"
)

// Funcs provided funcs for render
var Funcs = map[string]any{
	"netResolveIPAddr":     net.ResolveIPAddr,
	"netResolveIP":         netResolveIP,
	"osHostname":           os.Hostname,
	"osHostnameSequenceID": osHostnameSequenceID,
	"osUserCacheDir":       os.UserCacheDir,
	"osUserConfigDir":      os.UserConfigDir,
	"osUserHomeDir":        os.UserHomeDir,
	"osGetegid":            os.Getegid,
	"osGetenv":             os.Getenv,
	"osGeteuid":            os.Geteuid,
	"osGetgid":             os.Getgid,
	"osGetgroups":          os.Getgroups,
	"osGetpagesize":        os.Getpagesize,
	"osGetpid":             os.Getpid,
	"osGetppid":            os.Getppid,
	"osGetuid":             os.Getuid,
	"osGetwd":              os.Getwd,
	"osTempDir":            os.TempDir,
	"osUserLookupGroup":    user.LookupGroup,
	"osUserLookupGroupId":  user.LookupGroupId,
	"osUserCurrent":        user.Current,
	"osUserLookup":         user.Lookup,
	"osUserLookupId":       user.LookupId,
	"stringsContains":      strings.Contains,
	"stringsFields":        strings.Fields,
	"stringsIndex":         strings.Index,
	"stringsLastIndex":     strings.LastIndex,
	"stringsHasPrefix":     strings.HasPrefix,
	"stringsHasSuffix":     strings.HasSuffix,
	"stringsRepeat":        strings.Repeat,
	"stringsReplaceAll":    strings.ReplaceAll,
	"stringsSplit":         strings.Split,
	"stringsSplitN":        strings.SplitN,
	"stringsToLower":       strings.ToLower,
	"stringsToUpper":       strings.ToUpper,
	"stringsTrimPrefix":    strings.TrimPrefix,
	"stringsTrimSpace":     strings.TrimSpace,
	"stringsTrimSuffix":    strings.TrimSuffix,
	"strconvQuote":         strconv.Quote,
	"strconvUnquote":       strconv.Unquote,
	"strconvParseBool":     strconv.ParseBool,
	"strconvParseInt":      strconv.ParseInt,
	"strconvParseUint":     strconv.ParseUint,
	"strconvParseFloat":    strconv.ParseFloat,
	"strconvFormatBool":    strconv.FormatBool,
	"strconvFormatInt":     strconv.FormatInt,
	"strconvFormatUint":    strconv.FormatUint,
	"strconvFormatFloat":   strconv.FormatFloat,
	"strconvAoti":          strconv.Atoi,
	"strconvItoa":          strconv.Itoa,

	"add": add,
	"neg": neg,

	// deprecated
	"intAdd":           add,
	"intNeg":           neg,
	"int64Add":         add,
	"int64Neg":         neg,
	"float32Add":       add,
	"float32Neg":       neg,
	"float64Add":       add,
	"float64Neg":       neg,
	"k8sStatefulSetID": osHostnameSequenceID,
}

func netResolveIP(s string) (ip string, err error) {
	var addr *net.IPAddr
	if addr, err = net.ResolveIPAddr("ip", s); err != nil {
		return
	}
	ip = addr.IP.String()
	return
}

func add(a, b any) (any, error) {
	switch a.(type) {
	case bool:
		return a.(bool) || b.(bool), nil
		// ___BEG_GEN:ADD___
	case uint8:
		return a.(uint8) + b.(uint8), nil
	case uint16:
		return a.(uint16) + b.(uint16), nil
	case uint32:
		return a.(uint32) + b.(uint32), nil
	case uint64:
		return a.(uint64) + b.(uint64), nil
	case int8:
		return a.(int8) + b.(int8), nil
	case int16:
		return a.(int16) + b.(int16), nil
	case int32:
		return a.(int32) + b.(int32), nil
	case int64:
		return a.(int64) + b.(int64), nil
	case float32:
		return a.(float32) + b.(float32), nil
	case float64:
		return a.(float64) + b.(float64), nil
	case complex64:
		return a.(complex64) + b.(complex64), nil
	case complex128:
		return a.(complex128) + b.(complex128), nil
	case int:
		return a.(int) + b.(int), nil
	case uint:
		return a.(uint) + b.(uint), nil
	case string:
		return a.(string) + b.(string), nil
	case uintptr:
		return a.(uintptr) + b.(uintptr), nil
		// ___END_GEN:ADD___
	}
	return nil, errors.New("add: type not supported: " + reflect.TypeOf(a).String())
}

func neg(a any) (any, error) {
	switch a := a.(type) {
	case bool:
		return !a, nil
		// ___BEG_GEN:NEG___
	case uint8:
		return -a, nil
	case uint16:
		return -a, nil
	case uint32:
		return -a, nil
	case uint64:
		return -a, nil
	case int8:
		return -a, nil
	case int16:
		return -a, nil
	case int32:
		return -a, nil
	case int64:
		return -a, nil
	case float32:
		return -a, nil
	case float64:
		return -a, nil
	case complex64:
		return -a, nil
	case complex128:
		return -a, nil
	case int:
		return -a, nil
	case uint:
		return -a, nil
		// ___END_GEN:NEG___
	}
	return nil, errors.New("neg: type not supported: " + reflect.TypeOf(a).String())
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
		err = errors.New("missing sequence id in hostname")
		return
	}
	id, err = strconv.Atoi(splits[len(splits)-1])
	return
}

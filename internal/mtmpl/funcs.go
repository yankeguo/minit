package mtmpl

import (
	"errors"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

//go:generate python3 funcs.gen.py

// Funcs provided funcs for render
var Funcs = map[string]any{
	"filepathJoin":         filepath.Join,
	"netResolveIPAddr":     net.ResolveIPAddr,
	"netResolveIP":         funcNetResolveIP,
	"osHostname":           os.Hostname,
	"osHostnameSequenceID": funcOsHostnameSequenceID,
	"osReadDir":            os.ReadDir,
	"osReadFile":           os.ReadFile,
	"osReadFileString":     funcOsReadFileString,
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
	"strconvAtoi":          strconv.Atoi,
	"strconvItoa":          strconv.Itoa,

	"add":   funcAdd,
	"neg":   funcNeg,
	"dict":  funcDict,
	"slice": funcSlice,

	"int64":   funcInt64,
	"uint64":  funcUint64,
	"float64": funcFloat64,

	// deprecated
	"intAdd":           funcAdd,
	"intNeg":           funcNeg,
	"int64Add":         funcAdd,
	"int64Neg":         funcNeg,
	"float32Add":       funcAdd,
	"float32Neg":       funcNeg,
	"float64Add":       funcAdd,
	"float64Neg":       funcNeg,
	"k8sStatefulSetID": funcOsHostnameSequenceID,
}

func funcNetResolveIP(s string) (ip string, err error) {
	var addr *net.IPAddr
	if addr, err = net.ResolveIPAddr("ip", s); err != nil {
		return
	}
	ip = addr.IP.String()
	return
}

func funcAdd(a, b any) (any, error) {
	switch a := a.(type) {
	case bool:
		return a || b.(bool), nil
		// ___BEG_GEN:ADD___
	case uint8:
		return a + b.(uint8), nil
	case uint16:
		return a + b.(uint16), nil
	case uint32:
		return a + b.(uint32), nil
	case uint64:
		return a + b.(uint64), nil
	case int8:
		return a + b.(int8), nil
	case int16:
		return a + b.(int16), nil
	case int32:
		return a + b.(int32), nil
	case int64:
		return a + b.(int64), nil
	case float32:
		return a + b.(float32), nil
	case float64:
		return a + b.(float64), nil
	case int:
		return a + b.(int), nil
	case uint:
		return a + b.(uint), nil
	case complex64:
		return a + b.(complex64), nil
	case complex128:
		return a + b.(complex128), nil
	case string:
		return a + b.(string), nil
	case uintptr:
		return a + b.(uintptr), nil
		// ___END_GEN:ADD___
	}
	return nil, errors.New("add: type not supported: " + reflect.TypeOf(a).String())
}

func funcNeg(a any) (any, error) {
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
	case int:
		return -a, nil
	case uint:
		return -a, nil
	case complex64:
		return -a, nil
	case complex128:
		return -a, nil
		// ___END_GEN:NEG___
	}
	return nil, errors.New("neg: type not supported: " + reflect.TypeOf(a).String())
}

func funcDict(items ...any) (map[string]any, error) {
	if len(items)%2 != 0 {
		return nil, errors.New("dict: odd number of items")
	}
	m := map[string]any{}
	for i := 0; i < len(items); i += 2 {
		k, ok := items[i].(string)
		if !ok {
			return nil, errors.New("dict: key is not a string")
		}
		m[k] = items[i+1]
	}
	return m, nil
}

func funcSlice(args ...any) []any {
	return args
}

func funcOsHostnameSequenceID() (id int, err error) {
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

func funcOsReadFileString(path string) (string, error) {
	buf, err := os.ReadFile(path)
	return string(buf), err
}

func funcInt64(v any) (int64, error) {
	switch v := v.(type) {
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	case string:
		return strconv.ParseInt(v, 10, 64)
	case complex64:
		return int64(real(v)), nil
	case complex128:
		return int64(real(v)), nil
		// __BEG_GEN:INT64__
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case uint:
		return int64(v), nil
		// __END_GEN:INT64__
	}
	return 0, errors.New("int64: type not supported: " + reflect.TypeOf(v).String())
}

func funcUint64(v any) (uint64, error) {
	switch v := v.(type) {
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	case string:
		return strconv.ParseUint(v, 10, 64)
	case complex64:
		return uint64(real(v)), nil
	case complex128:
		return uint64(real(v)), nil
		// __BEG_GEN:UINT64__
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return uint64(v), nil
	case int8:
		return uint64(v), nil
	case int16:
		return uint64(v), nil
	case int32:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case float32:
		return uint64(v), nil
	case float64:
		return uint64(v), nil
	case int:
		return uint64(v), nil
	case uint:
		return uint64(v), nil
		// __END_GEN:UINT64__
	}
	return 0, errors.New("uint64: type not supported: " + reflect.TypeOf(v).String())
}

func funcFloat64(v any) (float64, error) {
	switch v := v.(type) {
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	case string:
		return strconv.ParseFloat(v, 64)
	case complex64:
		return float64(real(v)), nil
	case complex128:
		return real(v), nil
		// __BEG_GEN:FLOAT64__
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return float64(v), nil
	case int:
		return float64(v), nil
	case uint:
		return float64(v), nil
		// __END_GEN:FLOAT64__
	}
	return 0, errors.New("float64: type not supported: " + reflect.TypeOf(v).String())
}

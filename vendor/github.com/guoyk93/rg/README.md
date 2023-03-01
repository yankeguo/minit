# rg

`rg (Royal Guard)` is a generics based throw-catch approach in Go

## Usage

Any function with the latest return value of type `error` can be wrapped by `rg.Must` (or `rg.Must2`, `rg.Must3` ...)

## Example

```go
package demo

import (
	"encoding/json"
	"github.com/guoyk93/rg"
	"gopkg.in/yaml.v3"
	"os"
)

// jsonFileToYAMLUgly this is a demo function WITHOUT rg
func jsonFileToYAMLUgly(filename string) (err error) {
	var buf []byte
	if buf, err = os.ReadFile(filename); err != nil {
		return
	}
	var m map[string]interface{}
	if err = json.Unmarshal(buf, &m); err != nil {
		return
	}
	if buf, err = yaml.Marshal(m); err != nil {
		return
	}
	buf = rg.Must(yaml.Marshal(m))
	if err = os.WriteFile(filename+".yaml", buf, 0640); err != nil {
		return
	}
	return
}

// jsonFileToYAML this is a demo function WITH rg
func jsonFileToYAML(filename string) (err error) {
	defer rg.Guard(&err)
	buf := rg.Must(os.ReadFile(filename))
	var m map[string]interface{}
	rg.Must0(json.Unmarshal(buf, &m))
	buf = rg.Must(yaml.Marshal(m))
	rg.Must0(os.WriteFile(filename+".yaml", buf, 0640))
	return
}
```

## Donation

See https://guoyk.xyz/donation

## Credits

Guo Y.K., MIT License
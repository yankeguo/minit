//go:build !linux

package mexec

import "os/exec"

func setSysProcAttr(cmd *exec.Cmd) {

}

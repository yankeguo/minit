package mexec

import (
	"errors"
	"github.com/guoyk93/grace/gracelog"
	"github.com/guoyk93/minit/pkg/shellquote"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

type ExecuteOptions struct {
	Dir     string
	Shell   string
	Env     map[string]string
	Command []string
	Charset string

	Logger          gracelog.ProcLogger
	IgnoreExecError bool
}

type Manager interface {
	Signal(sig os.Signal)
	Execute(opts ExecuteOptions) (err error)
}

type manager struct {
	childPIDs    map[int]struct{}
	childPIDLock sync.Locker
	charsets     map[string]encoding.Encoding
}

func NewManager() Manager {
	return &manager{
		childPIDs:    map[int]struct{}{},
		childPIDLock: &sync.Mutex{},
		charsets: map[string]encoding.Encoding{
			"gb18030": simplifiedchinese.GB18030,
			"gbk":     simplifiedchinese.GBK,
		},
	}
}

func (m *manager) addChildPID(fn func() (pid int, err error)) error {
	m.childPIDLock.Lock()
	defer m.childPIDLock.Unlock()
	pid, err := fn()
	if err == nil {
		m.childPIDs[pid] = struct{}{}
	}
	return err
}

func (m *manager) delChildPID(pid int) {
	m.childPIDLock.Lock()
	defer m.childPIDLock.Unlock()
	delete(m.childPIDs, pid)
}

func (m *manager) Signal(sig os.Signal) {
	m.childPIDLock.Lock()
	defer m.childPIDLock.Unlock()
	for pid := range m.childPIDs {
		if process, _ := os.FindProcess(pid); process != nil {
			_ = process.Signal(sig)
		}
	}
}

func (m *manager) Execute(opts ExecuteOptions) (err error) {
	var argv []string

	// check opts.Dir
	if opts.Dir != "" {
		var info os.FileInfo
		if info, err = os.Stat(opts.Dir); err != nil {
			err = errors.New("failed to stat opts.Dir: " + err.Error())
			return
		}
		if !info.IsDir() {
			err = errors.New("opts.Dir is not a directory: " + opts.Dir)
			return
		}
	}

	// build env
	env := make(map[string]string)

	for _, item := range os.Environ() {
		splits := strings.SplitN(item, "=", 2)
		var k, v string
		if len(splits) > 0 {
			k = splits[0]
			if len(splits) > 1 {
				v = splits[1]
			}
			env[k] = v
		}
	}
	MergeEnv(env, opts.Env)

	// build argv
	if opts.Shell != "" {
		if argv, err = shellquote.Split(opts.Shell); err != nil {
			err = errors.New("opts.Shell is invalid: " + err.Error())
			return
		}
	} else {
		for _, arg := range opts.Command {
			argv = append(argv, os.Expand(arg, func(s string) string {
				return env[s]
			}))
		}
	}

	// build exec.Cmd
	var outPipe, errPipe io.Reader
	cmd := exec.Command(argv[0], argv[1:]...)
	if opts.Shell != "" {
		cmd.Stdin = strings.NewReader(strings.Join(opts.Command, "\n"))
	}
	cmd.Dir = opts.Dir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// build out / err pipe
	if outPipe, err = cmd.StdoutPipe(); err != nil {
		return
	}
	if errPipe, err = cmd.StderrPipe(); err != nil {
		return
	}

	// charset
	if opts.Charset != "" {
		enc := m.charsets[strings.ToLower(opts.Charset)]
		if enc == nil {
			opts.Logger.Error("unknown charset:", opts.Charset)
		} else {
			outPipe = enc.NewDecoder().Reader(outPipe)
			errPipe = enc.NewDecoder().Reader(errPipe)
		}
	}

	// start process in the same lock with signal children
	if err = m.addChildPID(func() (pid int, err error) {
		if err = cmd.Start(); err != nil {
			return
		}
		pid = cmd.Process.Pid
		return
	}); err != nil {
		return
	}

	opts.Logger.Print("process started")

	// streaming
	go opts.Logger.Out().ReadFrom(outPipe)
	go opts.Logger.Err().ReadFrom(errPipe)

	// wait for process
	if err = cmd.Wait(); err != nil {
		opts.Logger.Error("process exited:", err.Error())

		if opts.IgnoreExecError {
			err = nil
		}
	} else {
		opts.Logger.Print("process exited")
	}

	m.delChildPID(cmd.Process.Pid)

	return
}

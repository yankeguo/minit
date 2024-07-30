package mexec

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/yankeguo/minit/internal/menv"
	"github.com/yankeguo/minit/internal/mlog"
	"github.com/yankeguo/minit/pkg/shellquote"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type ExecuteOptions struct {
	Name string

	Dir          string
	Shell        string
	Env          map[string]string
	Command      []string
	Charset      string
	SuccessCodes []int

	Logger mlog.ProcLogger
}

type Manager interface {
	Signal(sig os.Signal)
	Execute(opts ExecuteOptions) (err error)
}

type manager struct {
	managedPIDs    map[int]struct{}
	managedPIDLock sync.Locker
	charsets       map[string]encoding.Encoding
}

func NewManager() Manager {
	return &manager{
		managedPIDs:    map[int]struct{}{},
		managedPIDLock: &sync.Mutex{},
		charsets: map[string]encoding.Encoding{
			"gb18030": simplifiedchinese.GB18030,
			"gbk":     simplifiedchinese.GBK,
		},
	}
}

func (m *manager) StartCommand(cmd *exec.Cmd) (done func(), err error) {
	m.managedPIDLock.Lock()
	defer m.managedPIDLock.Unlock()

	if err = cmd.Start(); err != nil {
		return
	}

	pid := cmd.Process.Pid
	m.managedPIDs[pid] = struct{}{}
	done = func() {
		m.managedPIDLock.Lock()
		defer m.managedPIDLock.Unlock()
		delete(m.managedPIDs, pid)
	}
	return
}

func (m *manager) Signal(sig os.Signal) {
	m.managedPIDLock.Lock()
	defer m.managedPIDLock.Unlock()

	for pid := range m.managedPIDs {
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
	var env map[string]string
	if env, err = menv.Construct(menv.Environ(), opts.Env); err != nil {
		err = errors.New("failed constructing environment variables: " + err.Error())
		return
	}

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
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
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
	var done func()
	if done, err = m.StartCommand(cmd); err != nil {
		return
	}
	defer done()

	opts.Logger.Printf("minit: %s: process started", opts.Name)

	// streaming
	go opts.Logger.Out().ReadFrom(outPipe)
	go opts.Logger.Err().ReadFrom(errPipe)

	// wait for process
	err = cmd.Wait()

	var code int

	if err != nil {
		opts.Logger.Errorf("minit: %s: process exited with error: %s", opts.Name, err.Error())
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			return
		}
	}

	if checkSuccessCode(opts.SuccessCodes, code) {
		err = nil
		opts.Logger.Printf("minit: %s: exit code %d is in success_codes", opts.Name, code)
		return
	}

	err = fmt.Errorf("exit code: %d is not in success_codes", code)

	opts.Logger.Errorf("minit: %s: process exited with error: %s", opts.Name, err.Error())

	return
}

func checkSuccessCode(successCodes []int, code int) bool {
	if len(successCodes) == 0 {
		return code == 0
	}
	for _, c := range successCodes {
		if c == code {
			return true
		}
	}
	return false
}

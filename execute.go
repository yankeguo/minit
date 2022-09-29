package main

import (
	"fmt"
	"github.com/guoyk93/grace/gracelog"
	"github.com/guoyk93/minit/pkg/shellquote"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	knownCharsets = map[string]encoding.Encoding{
		"gb18030": simplifiedchinese.GB18030,
		"gbk":     simplifiedchinese.GBK,
	}
)

var (
	childPids                 = map[int]struct{}{}
	childPidsLock sync.Locker = &sync.Mutex{}
)

type ExecuteOptions struct {
	Dir     string   `yaml:"dir"`     // 所有涉及命令执行的单元，指定命令执行时的当前目录
	Shell   string   `yaml:"shell"`   // 使用 shell 来执行命令，比如 'bash'
	Command []string `yaml:"command"` // 所有涉及命令执行的单元，指定命令执行的内容
	Charset string   `yaml:"charset"` // output charset
}

func addChildPID(fn func() (pid int, err error)) error {
	childPidsLock.Lock()
	defer childPidsLock.Unlock()
	pid, err := fn()
	if err == nil {
		childPids[pid] = struct{}{}
	}
	return err
}

func removeChildPID(pid int) {
	childPidsLock.Lock()
	defer childPidsLock.Unlock()
	delete(childPids, pid)
}

func notifyChildPIDs(sig os.Signal) {
	childPidsLock.Lock()
	defer childPidsLock.Unlock()
	for pid := range childPids {
		if process, _ := os.FindProcess(pid); process != nil {
			_ = process.Signal(sig)
		}
	}
}

func execute(opts ExecuteOptions, logger *gracelog.ProcLogger) (err error) {
	argv := make([]string, 0)

	// 检查 opts.Dir
	if opts.Dir != "" {
		var info os.FileInfo
		if info, err = os.Stat(opts.Dir); err != nil {
			err = fmt.Errorf("无法访问指定的 dir, 请检查: %s", err.Error())
			return
		}
		if !info.IsDir() {
			err = fmt.Errorf("指定的 dir 不是目录: %s", opts.Dir)
			return
		}
	}

	// 构建 argv
	if opts.Shell != "" {
		if argv, err = shellquote.Split(opts.Shell); err != nil {
			err = fmt.Errorf("无法处理 shell 参数，请检查: %s", err.Error())
			return
		}
	} else {
		for _, arg := range opts.Command {
			argv = append(argv, os.ExpandEnv(arg))
		}
	}

	// 构建 cmd
	var outPipe, errPipe io.Reader
	cmd := exec.Command(argv[0], argv[1:]...)
	if opts.Shell != "" {
		cmd.Stdin = strings.NewReader(strings.Join(opts.Command, "\n"))
	}
	cmd.Dir = opts.Dir
	// 阻止信号传递
	setupCmdSysProcAttr(cmd)

	if outPipe, err = cmd.StdoutPipe(); err != nil {
		return
	}
	if errPipe, err = cmd.StderrPipe(); err != nil {
		return
	}

	// charset
	if opts.Charset != "" {
		enc := knownCharsets[strings.ToLower(opts.Charset)]
		if enc == nil {
			logger.Error("未知字符集: " + opts.Charset)
		} else {
			outPipe = enc.NewDecoder().Reader(outPipe)
			errPipe = enc.NewDecoder().Reader(errPipe)
		}
	}

	// 执行
	if err = addChildPID(func() (pid int, err error) {
		// 在同一个锁内部启动进程，并记录 PID
		if err = cmd.Start(); err != nil {
			return
		}
		pid = cmd.Process.Pid
		return
	}); err != nil {
		return
	}

	// 串流
	go logger.Out().ReadFrom(outPipe)
	go logger.Err().ReadFrom(errPipe)

	// 等待退出
	if err = cmd.Wait(); err != nil {
		logger.Errorf("进程退出: %s", err.Error())
		err = nil
	} else {
		logger.Print("进程退出")
	}

	// 移除 Pid
	removeChildPID(cmd.Process.Pid)

	return
}

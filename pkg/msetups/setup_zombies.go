//go:build linux

package msetups

import (
	"bytes"
	"fmt"
	"github.com/guoyk93/minit/pkg/mlog"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func init() {
	Register(60, setupZombies)
}

func setupZombies(log mlog.ProcLogger) (err error) {
	// 如果自己不是 PID 1，则不负责清理僵尸进程
	if os.Getpid() != 1 {
		log.Print("minit is not running as PID 1, skipping cleaning up zombies")
		return
	}

	go runZombieCleaner(log)

	return
}

func runZombieCleaner(log mlog.ProcLogger) {
	// SIGCHLD 触发
	chSig := make(chan os.Signal, 10)
	signal.Notify(chSig, syscall.SIGCHLD)

	// 周期触发
	tk := time.NewTicker(time.Second * 30)

	var chT <-chan time.Time

	for {
		select {
		case <-chSig:
			if chT == nil {
				chT = time.After(time.Second * 3)
			}
		case <-tk.C:
			if chT == nil {
				chT = time.After(time.Second * 5)
			}
		case <-chT:
			chT = nil
			cleanZombieProcesses(log)
		}
	}
}

func cleanZombieProcesses(log mlog.ProcLogger) {
	var (
		err  error
		pids []int
	)
	if pids, err = findZombieProcesses(); err != nil {
		log.Print("failed checking zombies:", err.Error())
		return
	}

	for _, pid := range pids {
		go waitZombieProcess(log, pid)
	}
}

func findZombieProcesses() (pids []int, err error) {
	var f *os.File
	if f, err = os.Open("/proc"); err != nil {
		return
	}
	defer f.Close()
	var dirnames []string
	if dirnames, err = f.Readdirnames(-1); err != nil {
		return
	}
	for _, dirname := range dirnames {
		if dirname[0] < '0' || dirname[0] > '9' {
			continue
		}
		var pid int
		if pid, err = strconv.Atoi(dirname); err != nil {
			return
		}
		var zombie bool
		if zombie, err = checkProcessIsZombie(pid); err != nil {
			err = nil
			continue
		}
		if zombie {
			pids = append(pids, pid)
		}
	}
	return
}

func checkProcessIsZombie(pid int) (zombie bool, err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid)); err != nil {
		return
	}
	zombie = checkProcStatIsZombie(buf)
	return
}

func checkProcStatIsZombie(buf []byte) bool {
	if len(buf) == 0 {
		return false
	}
	idx := bytes.LastIndexByte(buf, ')')
	if idx < 0 {
		return false
	}
	buf = buf[idx+1:]
	buf = bytes.TrimSpace(buf)
	if len(buf) == 0 {
		return false
	}
	return buf[0] == 'Z'
}

func waitZombieProcess(log mlog.ProcLogger, pid int) {
	var err error
	var ws syscall.WaitStatus
	for {
		_, err = syscall.Wait4(pid, &ws, 0, nil)
		for syscall.EINTR == err {
			_, err = syscall.Wait4(pid, &ws, 0, nil)
		}
		if syscall.ECHILD == err {
			break
		}
	}
	log.Printf("zombie cleaned %d", pid)
}

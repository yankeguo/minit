package main

import (
	"context"
	"fmt"
	"github.com/guoyk93/grace/gracelog"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"
)

var (
	optUnitDir   string
	optLogDir    string
	optQuickExit bool
)

var (
	UnitNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9]$`)
)

var (
	log *gracelog.ProcLogger
)

var (
	GitHash = ""
)

func exit(err *error) {
	if *err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[%s] 错误退出: %s\n", "minit", (*err).Error())
		os.Exit(1)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "[%s] 正常退出\n", "minit")
	}
}

func main() {
	var err error
	defer exit(&err)

	StringEnv(&optUnitDir, "MINIT_UNIT_DIR", "/etc/minit.d")
	StringEnv(&optLogDir, "MINIT_LOG_DIR", "/var/log/minit")
	BoolEnv(&optQuickExit, "MINIT_QUICK_EXIT")

	// 确保配置单元目录
	if err = os.MkdirAll(optUnitDir, 0755); err != nil {
		return
	}

	// 确保日志目录
	if err = os.MkdirAll(optLogDir, 0755); err != nil {
		return
	}

	if log, err = gracelog.NewProcLogger(gracelog.ProcLoggerOptions{
		ConsolePrefix: "[minit] ",
		RotatingFileOptions: gracelog.RotatingFileOptions{
			Dir:      optLogDir,
			Filename: "minit",
		},
	}); err != nil {
		return
	}

	if GitHash == "" {
		GitHash = "UNKNOWN"
	}

	log.Print("minit (#" + GitHash + ")")

	// 自述文件
	setupBanner()

	// 内核参数
	if err = setupSysctl(); err != nil {
		return
	}

	// 资源限制
	if err = setupRLimits(); err != nil {
		return
	}

	// 透明大页
	if err = setupTHP(); err != nil {
		return
	}

	// WebDAV
	if err = setupWebDAV(); err != nil {
		return
	}

	// 载入单元
	var units []Unit
	if units, err = LoadDir(optUnitDir); err != nil {
		return
	}

	// 载入环境变量
	var (
		extraUnit Unit
		extraOK   bool
	)
	if extraUnit, extraOK, err = LoadEnvMain(); err != nil {
		return
	}
	if extraOK {
		units = append(units, extraUnit)
	}

	// 载入命令参数
	if extraUnit, extraOK, err = LoadArgsMain(); err != nil {
		return
	}
	if extraOK {
		units = append(units, extraUnit)
	}

	// 检查单元命名
	unitNames := map[string]bool{"minit": true}
	for _, unit := range units {
		if unit.Name == "" {
			err = fmt.Errorf("缺少单元名称，检查 name 字段")
			return
		}
		if !UnitNamePattern.MatchString(unit.Name) {
			err = fmt.Errorf("单元名称 %s 不符合规则，检查 name 字段", unit.Name)
			return
		}
		if unitNames[unit.Name] {
			err = fmt.Errorf("单元名称 %s 重复出现，检查 name 字段", unit.Name)
			return
		}
		unitNames[unit.Name] = true
		log.Printf("载入单元 %s/%s", unit.Kind, unit.Name)
	}

	// 控制器组, L1 是 render (渲染配置文件), L2 是 once (一次性命令), L3 是 daemon 和 cron
	runners := map[RunnerLevel][]Runner{}

	// 创建控制器
	for _, unit := range units {
		fac := RunnerFactories[unit.Kind]
		if fac == nil {
			err = fmt.Errorf("单元 %s 类型 %s 未知，检查 kind 字段", unit.Name, unit.Kind)
			return
		}

		var logger *gracelog.ProcLogger
		if logger, err = gracelog.NewProcLogger(gracelog.ProcLoggerOptions{
			RotatingFileOptions: gracelog.RotatingFileOptions{
				Dir:      optLogDir,
				Filename: unit.CanonicalName(),
			},
			ConsolePrefix: "[" + unit.Name + "] ",
		}); err != nil {
			err = fmt.Errorf("无法为 %s 创建日志: %s", unit.Name, err.Error())
			return
		}

		var runner Runner
		if runner, err = fac.Create(unit, logger); err != nil {
			err = fmt.Errorf("无法为 %s 创建控制器: %s", unit.Name, err.Error())
			return
		}

		runners[fac.Level] = append(runners[fac.Level], runner)
	}

	// 运行 L1 控制器
	for _, runner := range runners[RunnerL1] {
		runner.Run(context.Background())
	}
	// 运行 L2 控制器
	for _, runner := range runners[RunnerL2] {
		runner.Run(context.Background())
	}

	if len(runners[RunnerL3]) == 0 && optQuickExit {
		log.Printf("没有 L3 任务")
		return
	}

	// 运行 L3 控制器
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	for _, runner := range runners[RunnerL3] {
		wg.Add(1)
		go func(runner Runner) {
			runner.Run(ctx)
			wg.Done()
		}(runner)
	}

	log.Printf("启动完毕")

	// 启动僵尸进程清理工具
	setupZombies(log)

	// 等待信号并退出
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	sig := <-chSig
	log.Printf("接收到信号: %s", sig.String())

	// 关闭主环境
	cancel()

	// 延迟 3 秒播发信号
	time.Sleep(time.Second * 3)
	notifyChildPIDs(sig)

	// 等待控制器退出
	wg.Wait()
}

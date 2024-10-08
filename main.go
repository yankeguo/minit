package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/yankeguo/minit/internal/menv"
	"github.com/yankeguo/minit/internal/mexec"
	"github.com/yankeguo/minit/internal/mlog"
	"github.com/yankeguo/minit/internal/mrunners"
	"github.com/yankeguo/minit/internal/msetups"
	"github.com/yankeguo/minit/internal/munit"
	"github.com/yankeguo/rg"
)

var (
	AppVersion = "unknown"
)

func exit(err *error) {
	if *err == nil {
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, "%s: exited with error: %s\n", "minit", (*err).Error())
	os.Exit(1)
}

func envStr(key string, out *string) {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		*out = val
	}
}

func envBool(key string, out *bool) {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		*out, _ = strconv.ParseBool(val)
	}
}

func main() {
	var err error
	defer exit(&err)
	defer rg.Guard(&err)

	var (
		optPprofPort = ""
		optUnitDir   = "/etc/minit.d"
		optLogDir    = ""
		optQuickExit bool
	)

	// pprof
	if envStr("MINIT_PPROF_PORT", &optPprofPort); optPprofPort != "" {
		go func() {
			_ = http.ListenAndServe(":"+optPprofPort, nil)
		}()
	}

	envStr("MINIT_UNIT_DIR", &optUnitDir)
	envStr("MINIT_LOG_DIR", &optLogDir)
	envBool("MINIT_QUICK_EXIT", &optQuickExit)

	log := rg.Must(mlog.CreateSimpleLogger(optLogDir, "minit", "minit: "))

	exem := mexec.NewManager()

	log.Print("starting (" + AppVersion + ")")

	// run through setups
	rg.Must0(msetups.Setup(log))

	// load units
	units, skips := rg.Must2(
		munit.Load(
			munit.LoadOptions{
				Args: os.Args[1:],
				Env:  menv.Environ(),
				Dirs: munit.ParseUnitDirPattern(optUnitDir),
			},
		),
	)

	for _, skip := range skips {
		log.Print("unit skipped: " + skip.Name)
	}

	// load runners
	var (
		runnersS []mrunners.Runner
		runnersL []mrunners.Runner
	)

	{
		var runners []mrunners.Runner

		// convert units to runners
		for _, unit := range units {
			runners = append(
				runners,
				rg.Must(mrunners.Create(mrunners.RunnerOptions{
					Unit:   unit,
					Exec:   exem,
					Logger: rg.Must(mlog.CreateSimpleLogger(optLogDir, unit.Name, "")),
				})),
			)
		}

		// split short runners and long runners
		for _, runner := range runners {
			if runner.Long {
				runnersL = append(runnersL, runner)
			} else {
				runnersS = append(runnersS, runner)
			}
		}
	}

	// execute short runners
	for _, runner := range runnersS {
		if err = runner.Action.Do(context.Background()); err != nil {
			return
		}
	}

	// quick exit
	if len(runnersL) == 0 && optQuickExit {
		log.Printf("no long runners and MINIT_QUICK_EXIT is set")
		return
	}

	// run long runners
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	chErr := make(chan error, 1)

	for _, runner := range runnersL {
		wg.Add(1)
		go func(runner mrunners.Runner) {
			if err := runner.Action.Do(ctx); err != nil {
				select {
				case chErr <- err:
				default:
				}
			}
			wg.Done()
		}(runner)
	}

	// wait for signals
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)

	var sig os.Signal

	select {
	case sig = <-chSig:
		log.Printf("signal caught: %s", sig.String())
	case err = <-chErr:
		log.Printf("critical error caught: %s", err.Error())
	}

	if sig == nil {
		sig = syscall.SIGTERM
	}

	// shutdown context
	cancel()

	// delay 3 seconds
	time.Sleep(time.Second * 3)

	// broadcast signals
	exem.Signal(sig)

	// wait for long runners
	wg.Wait()
}

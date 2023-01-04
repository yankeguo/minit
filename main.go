package main

import (
	"context"
	"fmt"
	"github.com/guoyk93/gg"
	"github.com/guoyk93/gg/ggos"
	"github.com/guoyk93/minit/pkg/mexec"
	"github.com/guoyk93/minit/pkg/mlog"
	"github.com/guoyk93/minit/pkg/mrunners"
	"github.com/guoyk93/minit/pkg/msetups"
	"github.com/guoyk93/minit/pkg/munit"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

var (
	GitHash = "UNKNOWN"
)

func exit(err *error) {
	if *err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: exited with error: %s\n", "minit", (*err).Error())
		os.Exit(1)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "%s: exited\n", "minit")
	}
}

func main() {
	var err error
	defer exit(&err)
	defer gg.Guard(&err)

	var (
		optUnitDir   = "/etc/minit.d"
		optLogDir    = "/var/log/minit"
		optQuickExit bool
	)

	ggos.MustEnv("MINIT_UNIT_DIR", &optUnitDir)
	ggos.MustEnv("MINIT_LOG_DIR", &optLogDir)
	ggos.BoolEnv("MINIT_QUICK_EXIT", &optQuickExit)

	gg.Must0(os.MkdirAll(optUnitDir, 0755))
	gg.Must0(os.MkdirAll(optLogDir, 0755))

	log := gg.Must(mlog.NewProcLogger(mlog.ProcLoggerOptions{
		ConsolePrefix: "minit: ",
		RotatingFileOptions: mlog.RotatingFileOptions{
			Dir:      optLogDir,
			Filename: "minit",
		},
	}))

	exem := mexec.NewManager()

	log.Print("starting (#" + GitHash + ")")

	// run through setups
	gg.Must0(msetups.Setup(log))

	// load units
	loader := munit.NewLoader()
	units, skips := gg.Must2(
		loader.Load(
			munit.LoadOptions{
				Args: os.Args[1:],
				Env:  true,
				Dir:  optUnitDir,
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
				gg.Must(mrunners.Create(mrunners.RunnerOptions{
					Unit: unit,
					Exec: exem,
					Logger: gg.Must(mlog.NewProcLogger(mlog.ProcLoggerOptions{
						RotatingFileOptions: mlog.RotatingFileOptions{
							Dir:      optLogDir,
							Filename: unit.Name,
						},
					})),
				})),
			)
		}

		// sort runners
		sort.Slice(runners, func(i, j int) bool {
			return runners[i].Order < runners[j].Order
		})

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
		runner.Func.Do(context.Background())
	}

	// quick exit
	if len(runnersL) == 0 && optQuickExit {
		log.Printf("no long runners and MINIT_QUICK_EXIT is set")
		return
	}

	// run long runners
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	for _, runner := range runnersL {
		wg.Add(1)
		go func(runner mrunners.Runner) {
			runner.Func.Do(ctx)
			wg.Done()
		}(runner)
	}

	log.Printf("started")

	// wait for signals
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	sig := <-chSig
	log.Printf("signal caught: %s", sig.String())

	// shutdown context
	cancel()

	// delay 3 seconds
	time.Sleep(time.Second * 3)

	// broadcast signals
	exem.Signal(sig)

	// wait for long runners
	wg.Wait()
}

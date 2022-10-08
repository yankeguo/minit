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
		_, _ = fmt.Fprintf(os.Stderr, "[%s] exited with error: %s\n", "minit", (*err).Error())
		os.Exit(1)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "[%s] exited\n", "minit")
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
		ConsolePrefix: "[minit] ",
		RotatingFileOptions: mlog.RotatingFileOptions{
			Dir:      optLogDir,
			Filename: "minit",
		},
	}))

	exem := mexec.NewManager()

	log.Print("minit (#" + GitHash + ")")

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
	var runners []mrunners.Runner

	for _, unit := range units {
		runners = append(
			runners,
			gg.Must(mrunners.Create(mrunners.RunnerOptions{
				Unit: unit,
				Exec: exem,
				Logger: gg.Must(mlog.NewProcLogger(mlog.ProcLoggerOptions{
					ConsolePrefix: "[" + unit.Kind + "/" + unit.Name + "] ",
					RotatingFileOptions: mlog.RotatingFileOptions{
						Dir:      optLogDir,
						Filename: unit.Name,
					},
				})),
			})),
		)
	}

	sort.Slice(runners, func(i, j int) bool {
		return runners[i].Order > runners[j].Order
	})

	// run and remove short runners
	var n int
	for _, runner := range runners {
		if runner.Long {
			runners[n] = runner
			n++
		} else {
			runner.Func.Do(context.Background())
		}
	}
	runners = runners[:n]

	// quick exit
	if len(runners) == 0 && optQuickExit {
		log.Printf("no long runners and MINIT_QUICK_EXIT is set")
		return
	}

	// run long runners
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	for _, runner := range runners {
		wg.Add(1)
		go func(runner mrunners.Runner) {
			runner.Func.Do(ctx)
			wg.Done()
		}(runner)
	}

	log.Printf("booted")

	// wait for signals
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	sig := <-chSig
	log.Printf("signal caught: %s", sig.String())

	// shutdown context
	cancel()

	// dely 3 seconds
	time.Sleep(time.Second * 3)

	// broadcast signals
	exem.Signal(sig)

	// wait for long runners
	wg.Wait()
}

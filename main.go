package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const (
	PerPlotThreads = 2
)

var (
	PerPlotMem    = ByteSzFromMB(3200)
	TmpPlotSpace  = ByteSzFromGiB(356)
	FarmPlotSpace = ByteSzFromGiB(101.4 + .2)
)

var (
	ErrMaxProcessesReached = fmt.Errorf("max processes reached")
)

func main() {
	loadEnv()
	mem := getMemStats()
	logF("Starting chiarunner...\n"+
		"System CPU threads: %d\n"+
		"System Free mem: %s\n"+
		"System Total mem: %s\n",
		runtime.NumCPU(),
		mem.Free.String(),
		mem.Total.String())

	r := newRunner(ByteSzFromMB(float64(env.MaxMemoryMB)))
	logF("Max parallel plots: %d\n", r.MaxParallelPlots())

	for _, d := range env.FarmDirs {
		fd := NewFarmDir(d)
		r.FarmPool.AddDirs(fd)
		logF("added farm directory %s", d)
	}

	for _, d := range env.PlotDirs {
		pd := newPlotDir(d)
		r.PlotPool.AddDirs(pd)
		logF("added plot directory %s\n", d)
	}

	// log the current status
	logLn(r.StatusString())

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		sig := <-sigs
		logErrLn("signal", sig, "called", ". Terminating...")
		cancel()
	}()
	r.runner(ctx, time.Minute)
	runtime.SetFinalizer(r, func(r *Runner) {
		cancel()
	})
}

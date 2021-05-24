package main

import (
	"context"
	"flag"
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
	// parse flags
	memMB := flag.Int("m", 0, "max memory in MB")
	pltDirs := flag.String("p", "", "comma delimited list of plotting dirs")
	frmDirs := flag.String("f", "", "comma delimited list of farming dirs")
	logDir := flag.String("o", "", "log file")
	flag.Parse()

	initLogger(logDir)

	mem := getMemStats()
	logF("Starting chiarunner...\n" +
		"System CPU threads: %d\n" +
		"System Free mem: %s\n" +
		"System Total mem: %s\n",
		runtime.NumCPU(),
		mem.Free.String(),
		mem.Total.String())

	if memMB == nil || *memMB <= 0 {
		logFatalLn("-m must be set")
	}
	if pltDirs == nil || len(*pltDirs) < 1 {
		logFatalLn("-p must be set")
	}
	if frmDirs == nil || len(*frmDirs) < 1 {
		logFatalLn("-f must be set")
	}
	//fmt.Printf("Farm Space: %s\n", FarmPlotSpace)
	//fmt.Printf("plot Space: %s\n", TmpPlotSpace)

	r := newRunner(ByteSzFromMB(float64(*memMB)))
	logF("Max parallel plots: %d\n", r.MaxParallelPlots())

	for _, d := range farmDirsFromString(*frmDirs) {
		r.FarmPool.AddDirs(d)
		d.FreeSpace()
		d.FarmingSpaceAvail()
		d.CanAddPlot()
		logF("adding farm directory %s...\n" +
			"%s free space: %s\n" +
			"%s farm plots available: %d\n",
			d.dirStr,
			d.dirStr, d.FreeSpace(),
			d.dirStr, int(d.FreeSpace() / FarmPlotSpace))
	}

	for _, d := range plotDirsFromString(*pltDirs) {
		r.PlotPool.AddDirs(d)
		d.FreeSpace()
		d.PlottingSpaceAvail()
		d.CanPlot()
		logF("adding plot directory %s...\n" +
			"%s free space: %s\n" +
			"%s temp plots available: %d\n",
			d.dirStr,
			d.dirStr, d.FreeSpace(),
			d.dirStr, int(d.FreeSpace() / TmpPlotSpace))
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		sig := <- sigs
		logErrLn("signal", sig, "called", ". Terminating...")
		cancel()
	}()
	r.runner(ctx, time.Minute)
	runtime.SetFinalizer(r, func(r *Runner){
		cancel()
	})
}

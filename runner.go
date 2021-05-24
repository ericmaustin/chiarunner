package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

//newRunner creates a new Runner with the given maximum memory bytes
func newRunner(maxMemBytes ByteSz) *Runner {
	cpuMax := runtime.NumCPU() / PerPlotThreads
	memMax := maxMemBytes / PerPlotMem

	math.Floor(math.Min(float64(cpuMax), float64(memMax)))

	return &Runner{
		PlotPool: &PlotPool{
			mu: &sync.RWMutex{},
		},
		FarmPool: &FarmPool{
			mu: &sync.RWMutex{},
		},
		activeProcesses: map[int]*os.Process{},
		maxPlots:        int(math.Floor(math.Min(float64(cpuMax), float64(memMax)))),
		mu:              &sync.RWMutex{},
	}
}

//Runner maintains a PlotPool and FarmPool as well as a map of active processes
// Runner is used to spawn new plot processes
type Runner struct {
	PlotPool        *PlotPool
	FarmPool        *FarmPool
	activeProcesses map[int]*os.Process
	maxPlots        int
	mu              *sync.RWMutex
}

//MaxParallelPlots returns the maximum number of new plots
func (r *Runner) MaxParallelPlots() int {
	return r.maxPlots - len(r.activeProcesses)
}

// plot attempts to create a new plot by running the chia plots create command using the next available
// plotting dir and farming dir
// if no space is available or not enough memory or cpu resources are available, then this returns
// an ErrMaxProcessesReached error
// commands are executed and then waited on in a separate go routine
func (r *Runner) plot() error {
	if r.MaxParallelPlots() < 1 {
		return ErrMaxProcessesReached
	}

	logLn("starting plot process...")

	r.mu.Lock()
	defer r.mu.Unlock()

	plotDir, err := r.PlotPool.NextUp()
	if err != nil {
		return err
	}
	logLn("plot dir", plotDir.dirStr, "has been selected with", plotDir.FreeSpace() ,"free space")

	farmDir, err := r.FarmPool.NextUp()
	if err != nil {
		return err
	}
	logLn("farm dir", plotDir.dirStr, "has been selected with", farmDir.FreeSpace() ,"free space")

	args := []string{
		"plots",
		"create",
		"-k", "32",
		"-r", "2",
		"-b", fmt.Sprintf("%d", int(PerPlotMem.MB())),
		"-t", plotDir.dirStr,
		"-d", farmDir.dirStr,
	}

	cmd := exec.Command("chia", args...)
	logLn("running cmd: chia", strings.Join(args, " "))

	err = cmd.Start()
	if err != nil {
		logErrLn("cmd failed!")
		return err
	}
	pid := cmd.Process.Pid
	r.activeProcesses[pid] = cmd.Process
	plotDir.AddPID(pid)
	farmDir.AddPID(pid)
	logF("[%d] now plotting. plot dir:%s farm dir:%s\n", pid, plotDir.dirStr, farmDir.dirStr)
	go r.waitForCmd(cmd, plotDir, farmDir)
	return nil
}

//waitForCmd waits for an exec.Cmd to complete, then removes the PID from the plot and farm dirs and removes
// the process from the active process slice
func (r *Runner) waitForCmd(cmd *exec.Cmd, plotDir *PlotDir, farmDir *FarmDir) {
	pid := cmd.Process.Pid
	err := cmd.Wait()
	if err != nil {
		logF("process %d finished with error: %v\n", pid, err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	// cleanup after our process
	plotDir.RmPID(pid)
	farmDir.RmPID(pid)
	delete(r.activeProcesses, pid)
	logF("process %d finished\n", pid)
}

//killAll kills all the active processes
func (r *Runner) killAll() {
	for pid, proc := range r.activeProcesses {
		if err := proc.Kill(); err != nil {
			logErrLn("failed to kill plot process", pid)
		} else {
			logLn("killed plot process", pid)
		}
	}
}

//runner is the actual worker
func (r *Runner) runner(ctx context.Context, waitDur time.Duration) {
	ticker := time.NewTicker(waitDur)

	// first plot cmd before the for loop
	if err := r.plot(); err != nil && err != ErrMaxProcessesReached {
		logFatalLn("plot error:", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			logLn("context done, runner exiting...")
			r.killAll()
			return
		case  <- ticker.C:
			// got tick, try to plot
			if err := r.plot(); err != nil && err != ErrMaxProcessesReached {
				logFatalLn("plot error:", err)
				return
			}
		}
	}
}
package main

import (
	"fmt"
	"strings"
	"sync"
)

//newDir creates a new dir with the given dir string
func newDir(dirStr string) dir {
	return dir{
		dirStr:     dirStr,
		mu:         &sync.RWMutex{},
		activePIDs: map[int]int{},
	}
}

//dir represents a dir
type dir struct {
	dirStr     string
	activePIDs map[int]int
	mu         *sync.RWMutex
}

//newPlotDir crates anew PlotDir with the given dir string
func newPlotDir(dirStr string) *PlotDir {
	return &PlotDir{
		dir: newDir(dirStr),
	}
}

func (d *dir) DiskStat() *DiskStat {
	return diskStat(d.dirStr)
}

//PlotDir represents a dir used for plotting
type PlotDir struct {
	dir
}

//AddPID adds the given PID int to the active pid map
func (p *PlotDir) AddPID(pid int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.activePIDs[pid] = pid
}

func (p *PlotDir) RmPID(pid int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.activePIDs, pid)
}

func (p *PlotDir) TempSpace() ByteSz {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return ByteSz(int64(len(p.activePIDs)) * int64(PerPlotMem))
}

func (p *PlotDir) AvailableSpace() ByteSz {
	return diskStat(p.dirStr).Available
}

func (p *PlotDir) PlottingSpaceAvail() ByteSz {
	return p.AvailableSpace().Sub(p.TempSpace())
}

func (p *PlotDir) CanPlot() bool {
	return p.PlottingSpaceAvail() > TmpPlotSpace
}

func NewFarmDir(dir string) *FarmDir {
	return &FarmDir{
		dir: newDir(dir),
	}
}

//FarmDir represents a dir used for farming
type FarmDir struct {
	dir
}

//AddPID adds the given PID int to the active pid map
func (f *FarmDir) AddPID(pid int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.activePIDs[pid] = pid
}

//AddPID adds the given PID int to the active pid map
func (f *FarmDir) RmPID(pid int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.activePIDs, pid)
}

func (f *FarmDir) TempSpace() ByteSz {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return ByteSz(int64(len(f.activePIDs)) * int64(FarmPlotSpace))
}

func (f *FarmDir) AvailableSpace() ByteSz {
	return diskStat(f.dirStr).Available
}

func (f *FarmDir) FarmingSpaceAvail() ByteSz {
	return f.AvailableSpace().Sub(f.TempSpace())
}

func (f *FarmDir) CanAddPlot() bool {
	return f.FarmingSpaceAvail() > FarmPlotSpace
}

type PlotPool struct {
	PlotDirs []*PlotDir
	ptr      int
	mu       *sync.RWMutex
}

func (p *PlotPool) AddDirs(plotDirs ...*PlotDir) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, pl := range plotDirs {
		p.PlotDirs = append(p.PlotDirs, pl)
	}
}

func (p *PlotPool) DirCnt() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.PlotDirs)
}

func (p *PlotPool) next() *PlotDir {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.PlotDirs) < 1 {
		panic(fmt.Errorf("no plot dirs"))
	}
	newPtr := p.ptr + 1
	if newPtr > len(p.PlotDirs)-1 {
		newPtr = 0
	}
	p.ptr = newPtr
	return p.PlotDirs[newPtr]
}

func (p *PlotPool) NextUp() (*PlotDir, error) {
	for i := 0; i < p.DirCnt(); i++ {
		pl := p.next()
		if pl.CanPlot() {
			return pl, nil
		}
	}
	return nil, ErrMaxProcessesReached
}

type FarmPool struct {
	FarmDirs []*FarmDir
	ptr      int
	mu       *sync.RWMutex
}

func (f *FarmPool) AddDirs(farmDirs ...*FarmDir) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, fd := range farmDirs {
		f.FarmDirs = append(f.FarmDirs, fd)
	}
}

func (f *FarmPool) DirCnt() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.FarmDirs)
}

func (f *FarmPool) next() *FarmDir {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.FarmDirs) < 1 {
		panic(fmt.Errorf("no plot dirs"))
	}
	newPtr := f.ptr + 1
	if newPtr > len(f.FarmDirs)-1 {
		newPtr = 0
	}
	f.ptr = newPtr
	return f.FarmDirs[newPtr]
}

func (f *FarmPool) NextUp() (*FarmDir, error) {
	for i := 0; i < f.DirCnt(); i++ {
		fd := f.next()
		if fd.CanAddPlot() {
			return fd, nil
		}
	}
	return nil, ErrMaxProcessesReached
}

func plotDirsFromString(dirStr string) []*PlotDir {
	dirs := strings.Split(dirStr, ",")
	out := make([]*PlotDir, len(dirs))
	for i, d := range dirs {
		out[i] = newPlotDir(strings.TrimSpace(d))
	}
	return out
}

func farmDirsFromString(dirStr string) []*FarmDir {
	dirs := strings.Split(dirStr, ",")
	out := make([]*FarmDir, len(dirs))
	for i, d := range dirs {
		out[i] = NewFarmDir(strings.TrimSpace(d))
	}
	return out
}

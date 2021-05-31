package main

import (
	"fmt"
	"github.com/mackerelio/go-osstat/memory"
	"golang.org/x/sys/unix"
	"math"
)

//ByteSz is used for byte sizes
type ByteSz int64

//B returns the the ByteSz in bytes
func (b ByteSz) B() int64 {
	return int64(b)
}

//KB returns the the ByteSz in SCI kilobytes
func (b ByteSz) KB() float64 {
	return float64(b) / 1000
}

//KiB returns the the ByteSz in IEC kilobytes
func (b ByteSz) KiB() float64 {
	return float64(b) / 1024
}

//MB returns the the ByteSz in SCI megabytes
func (b ByteSz) MB() float64 {
	return b.KB() / 1000
}

//MiB returns the the ByteSz in IEC megabytes
func (b ByteSz) MiB() float64 {
	return b.KiB() / 1024
}

//GB returns the the ByteSz in SCI gigabytes
func (b ByteSz) GB() float64 {
	return b.MB() / 1000
}

//GiB returns the the ByteSz in IEC gigabytes
func (b ByteSz) GiB() float64 {
	return b.MiB() / 1024
}

//TB returns the the ByteSz in SCI terabytes
func (b ByteSz) TB() float64 {
	return b.GB() / 1000
}

//TiB returns the the ByteSz in IEC terabytes
func (b ByteSz) TiB() float64 {
	return b.GiB() / 1024
}

//Sub subtracts another ByteSz
func (b ByteSz) Sub(other ByteSz) ByteSz {
	return ByteSz(int64(b) - int64(other))
}

//Add adds another ByteSz
func (b ByteSz) Add(other ByteSz) ByteSz {
	return ByteSz(int64(b) + int64(other))
}

//String prints the ByteSz as a smart-formatted string according to its size
func (b ByteSz) String() string {
	i := float64(b)
	if i > 1e12 {
		return fmt.Sprintf("%.6f TB", b.TB())
	}
	if i > 1e9 {
		return fmt.Sprintf("%.6f GB", b.GB())
	}
	if i > 1e6 {
		return fmt.Sprintf("%.6f MB", b.MB())
	}
	if i > 1e3 {
		return fmt.Sprintf("%.6f KB", b.KB())
	}

	return fmt.Sprintf("%d B", b)
}

//ByteSzFromGB creates a new ByteSz from a float64 of SCI gigabytes
func ByteSzFromGB(gb float64) ByteSz {
	return ByteSz(gb * 1e9)
}

//ByteSzFromMB creates a new ByteSz from a float64 of SCI megabytes
func ByteSzFromMB(mb float64) ByteSz {
	return ByteSz(mb * 1e6)
}

//ByteSzFromGiB creates a new ByteSz from a float64 of IEC gigabytes
func ByteSzFromGiB(gb float64) ByteSz {
	return ByteSz(int64(gb * math.Pow(1024, 3)))
}

// MemStats represents memory statistics for darwin
type MemStats struct {
	Total, Used, Cached, Free, Active, Inactive, SwapTotal, SwapUsed, SwapFree ByteSz
}

//getMemStats gets a MemStats ptr with sizes in ByteSz
func getMemStats() *MemStats {
	mem, err := memory.Get()
	if err != nil {
		panic(err)
	}
	return &MemStats{
		Total:     ByteSz(mem.Total),
		Used:      ByteSz(mem.Used),
		Cached:    ByteSz(mem.Cached),
		Free:      ByteSz(mem.Free),
		Active:    ByteSz(mem.Active),
		Inactive:  ByteSz(mem.Inactive),
		SwapTotal: ByteSz(mem.SwapTotal),
		SwapUsed:  ByteSz(mem.SwapUsed),
		SwapFree:  ByteSz(mem.SwapFree),
	}
}

// DiskStat contains the disk stat details
type DiskStat struct {
	Used      ByteSz
	Available ByteSz
	Total     ByteSz
}

//diskStat gets the free disk space as a ByteSz for the given dir
func diskStat(dir string) *DiskStat {
	var stat unix.Statfs_t

	if err := unix.Statfs(dir, &stat); err != nil {
		logFatalF("could not get disk status of %s: %v", dir, err)
	}

	ds := &DiskStat{
		Available: ByteSz(stat.Bavail * uint64(stat.Bsize)),
		Total:     ByteSz(stat.Blocks * uint64(stat.Bsize)),
	}

	ds.Used = ds.Total.Sub(ds.Available)

	return ds
}

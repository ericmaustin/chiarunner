package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"math"
	"os"
	"runtime"
	"strings"
)

type envVars struct {
	ChiaDir          string
	MaxMemoryMB      int
	PlotDirs         []string
	FarmDirs         []string
	LogFile          string
	SMTPHost         string
	SMTPPort         int
	SMTPUser         string
	SMTPPassword     string
	EmailFrom        string
	EmailTo          []string
	PerPlotMemMB     int
	PerPlotThreads   int
	MaxParallelPlots int
}

func (e *envVars) PerPlotMem() ByteSz {
	return ByteSzFromMB(float64(e.PerPlotMemMB))
}

func (e *envVars) MaxMemory() ByteSz {
	return ByteSzFromMB(float64(e.MaxMemoryMB))
}

var env *envVars

var (
	flagConfigFile,
	flagLogFile,
	flagPlottingDirs,
	flagFarmingDirs,
	flagSMTPHost,
	flagSMTPUser,
	flagSMTPPass,
	flagEmailTo,
	flagEmailFrom,
	flagChiaDir string

	flagMaxMem,
	flagPerPlotMem,
	flagPerPlotThreads,
	flagSMTPPort int
)

func loadEnv() {

	flag.Parse()

	env = new(envVars)

	//if config file is passed, attempt to parse the toml file
	if len(flagConfigFile) > 0 {
		b, err := os.ReadFile(flagConfigFile)
		if err != nil {
			logFatalLn(err)
		}
		if err = toml.Unmarshal(b, env); err != nil {
			logFatalLn(err)
		}
	}

	if flagMaxMem > 0 {
		env.MaxMemoryMB = flagMaxMem
	}

	if flagPerPlotMem > 0 {
		env.PerPlotMemMB = flagPerPlotMem
	} else if env.PerPlotMemMB <= 0 {
		// default per plot mem MB
		env.PerPlotMemMB = 3200
	}

	if flagPerPlotThreads > 0 {
		env.PerPlotThreads = flagPerPlotThreads
	} else if env.PerPlotThreads <= 0 {
		// default per plot threads
		env.PerPlotThreads = 2
	}

	if len(flagLogFile) > 0 {
		env.LogFile = flagLogFile
	}

	if len(flagPlottingDirs) > 0 {
		dirs := strings.Split(flagPlottingDirs, ",")
		env.PlotDirs = make([]string, len(dirs))
		for i, d := range dirs {
			env.PlotDirs[i] = strings.TrimSpace(d)
		}
	}

	if len(flagFarmingDirs) > 0 {
		dirs := strings.Split(flagFarmingDirs, ",")
		env.FarmDirs = make([]string, len(dirs))
		for i, d := range dirs {
			env.FarmDirs[i] = strings.TrimSpace(d)
		}
	}

	if len(flagSMTPHost) > 0 {
		env.SMTPHost = flagSMTPHost
	}

	if flagSMTPPort > 0 {
		env.SMTPPort = flagSMTPPort
	}

	if len(flagSMTPUser) > 0 {
		env.SMTPUser = flagSMTPUser
	}

	if len(flagSMTPPass) > 0 {
		env.SMTPPassword = flagSMTPPass
	}

	if len(flagEmailFrom) > 0 {
		env.EmailFrom = flagEmailFrom
	}

	if len(flagEmailTo) > 0 {
		addresses := strings.Split(flagEmailTo, ",")
		env.EmailTo = make([]string, len(addresses))
		for i, to := range addresses {
			env.EmailTo[i] = to
		}
	}

	if env.MaxParallelPlots <= 0 {
		cpuMax := runtime.NumCPU() / env.PerPlotThreads
		memMax := env.MaxMemoryMB / env.PerPlotMemMB
		env.MaxParallelPlots = int(math.Floor(math.Min(float64(cpuMax), float64(memMax))))
	}
}

func init() {
	// chia blockchain directory
	flag.StringVar(&flagChiaDir, "chia-dir", "~/chia-blockchain", "chia blockchain directory")
	flag.StringVar(&flagConfigFile, "config", "", "config TOML file to use")
	// max memory flag
	flag.IntVar(&flagMaxMem, "max-mem", 0, "max memory in MB")
	flag.IntVar(&flagPerPlotMem, "plot-mem", 0, "max memory to use per plot")
	flag.IntVar(&flagPerPlotThreads, "plot-threads", 0, "cpu threads to use per plot")
	// log file flag
	flag.StringVar(&flagLogFile, "log", "", "log output file")
	// plotting dirs flag
	flag.StringVar(&flagPlottingDirs, "temp-dirs", "", "comma delimited list of temporary plotting dirs")
	// farming dirs flag
	flag.StringVar(&flagFarmingDirs, "farm-dirs", "", "comma delimited list of farming dirs")
	// smtp flags
	flag.StringVar(&flagSMTPHost, "smtp-host", "", "SMTP server host")
	flag.IntVar(&flagSMTPPort, "smtp-port", 0, "SMTP server port")
	flag.StringVar(&flagSMTPUser, "smtp-user", "", "SMTP username")
	flag.StringVar(&flagSMTPPass, "smtp-pass", "", "SMTP password")
	// email flags
	flag.StringVar(&flagEmailTo, "email-to", "", "comma delimited list email addresses to send alerts to")
	flag.StringVar(&flagEmailFrom, "email-from", "", "from email address")

}

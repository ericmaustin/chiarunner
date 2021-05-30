package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"os"
	"strings"
)

type envVars struct {
	MaxMemoryMB  ByteSz
	PlotDirs     []string
	FarmDirs     []string
	LogFile      string
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	EmailFrom    string
	EmailTo      []string
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
	flagEmailFrom string

	flagMaxMem,
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
		env.MaxMemoryMB = ByteSzFromMB(float64(flagMaxMem))
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
}

func init() {
	flag.StringVar(&flagConfigFile, "c", "", "config TOML file to use")
	// max memory flag
	flag.IntVar(&flagMaxMem, "m", 0, "max memory in MB")
	// log file flag
	flag.StringVar(&flagLogFile, "o", "", "log output file")
	// plotting dirs flag
	flag.StringVar(&flagPlottingDirs, "p", "", "comma delimited list of plotting dirs")
	// farming dirs flag
	flag.StringVar(&flagFarmingDirs, "f", "", "comma delimited list of farming dirs")
	// smtp flags
	flag.StringVar(&flagSMTPHost, "smtp-host", "", "SMTP server host")
	flag.IntVar(&flagSMTPPort, "smtp-port", 0, "SMTP server port")
	flag.StringVar(&flagSMTPUser, "smtp-user", "", "SMTP username")
	flag.StringVar(&flagSMTPPass, "smtp-pass", "", "SMTP password")
	// email flags
	flag.StringVar(&flagEmailTo, "email-to", "", "comma delimited list email addresses to send alerts to")
	flag.StringVar(&flagEmailFrom, "email-from", "", "from email address")

}

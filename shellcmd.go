package main

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
)

type ShellCmdBuilder struct {
	shell string
	shellArgs []string
	cmdSep string
	Commands []*exec.Cmd
	cmd      *exec.Cmd
}

func (s *ShellCmdBuilder) AddCmd(cmd *exec.Cmd) *ShellCmdBuilder {
	s.Commands = append(s.Commands, cmd)
	s.buildCmd()
	return s
}

func (s*ShellCmdBuilder) buildCmd() {
	cmdStrs := make([]string, len(s.Commands))
	for i, cmd := range s.Commands {
		cmdStrs[i] = cmd.String()
	}
	s.cmd = exec.Command(s.shell, append(s.shellArgs, strings.Join(cmdStrs, s.cmdSep))...)
}

func (s *ShellCmdBuilder) SetCmdSep(sep string) *exec.Cmd {
	s.cmdSep = sep
	return s.cmd
}

func (s *ShellCmdBuilder) Cmd() *exec.Cmd {
	return s.cmd
}

func NewShellCmdBuilder(shell string, arg ...string) *ShellCmdBuilder {
	return &ShellCmdBuilder{
		shell:    shell,
		shellArgs: arg,
		cmdSep: "; ",
	}
}

func newChiaBaseCmd() *ShellCmdBuilder {
	cmd := NewShellCmdBuilder("/bin/bash", "-c")
	cmd.AddCmd(exec.Command("source", path.Join(env.ChiaDir, "activate")))
	return cmd
}

func PlotCmd(tmpDir, farmDir string) *exec.Cmd {
	shellCmd := newChiaBaseCmd()
	shellCmd.AddCmd(exec.Command("chia",
		"plots",
		"create",
		"-k", "32",
		"-r", fmt.Sprintf("%d", env.PerPlotThreads),
		"-b", fmt.Sprintf("%d", env.PerPlotMemMB),
		"-t", tmpDir,
		"-d", farmDir))
	return shellCmd.Cmd()
}

func WalletShowCmd() *exec.Cmd {
	shellCmd := newChiaBaseCmd()
	shellCmd.AddCmd(exec.Command("chia","wallet", "show"))
	return shellCmd.Cmd()
}

func FarmSummaryCmd() *exec.Cmd {
	shellCmd := newChiaBaseCmd()
	shellCmd.AddCmd(exec.Command("chia","farm", "summary"))
	return shellCmd.Cmd()
}
package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestShellCmdBuilder(t *testing.T){
	shCmd := new(ShellCmdBuilder)
	shCmd.AddCmd(exec.Command("echo", "hello"))
	shCmd.AddCmd(exec.Command("echo", "world"))
	fmt.Println(shCmd.Cmd().String())
	out, err := shCmd.Cmd().Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("cmd output:\n", string(out))
}

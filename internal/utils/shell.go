package utils

import (
	"os"
	"os/exec"

	"github.com/mitchellh/go-ps"
)

type Shell struct {
	name string
}

func CurrentShell() (s Shell) {
	proc, _ := ps.FindProcess(os.Getppid())
	parentName := proc.Executable()
	if parentName == "go" {
		// Executed with `go run .`
		proc, _ = ps.FindProcess(proc.PPid())
		parentName = proc.Executable()
	}
	switch parentName {
	case "bash", "zsh", "fish", "pwsh":
		s = Shell{name: parentName}
	default:
		s = Shell{name: "sh"}
	}
	return
}

func (s Shell) Run(command string, args []string) error {
	var escapedArgs string
	for _, arg := range args {
		escapedArgs += "\"" + arg + "\" "
	}

	cmd := exec.Command(s.name, "-c", command+" "+escapedArgs)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func RunInCurrentShell(command string, args []string) error {
	return CurrentShell().Run(command, args)
}

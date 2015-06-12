package gitdeploy

import (
	"os/exec"
	"syscall"
)

func runCmd(name string, args ...string) ([]byte, int, error) {
	return runCmdInDir("", name, args...)
}

func runCmdInDir(dir, name string, args ...string) ([]byte, int, error) {
	c := exec.Command(name, args...)

	if dir != "" {
		c.Dir = dir
	}

	output, err := c.CombinedOutput()

	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return nil, status.ExitStatus(), err
		}

		return output, -1, err
	}

	return output, 0, err
}

func overlayFiles(source, target *Environment) {
	// same as above, what happens when one fails?
	// clear the target out
	target.runCmdInDir("git", "rm", "-r", ".")

	source.runCmdInDir("rsync", "-av", ".", "../../"+target.getDir(), "--exclude", ".git")
}

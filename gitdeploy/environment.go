package gitdeploy

import (
	"path"
	"strings"
)

type Environment struct {
	Name          string
	RepoURL       string
	Branch        string
	PreventDeploy bool

	DefaultSourceEnvironmentName string

	project *Project
}

func (e *Environment) getWholeDir() string {
	return path.Join(e.project.baseDirectory, e.getDir())
}

func (e *Environment) getDir() string {
	return path.Join(e.project.Name, e.Name)
}

func (e *Environment) runCmdInDir(name string, args ...string) ([]byte, int, error) {
	return runCmdInDir(e.getWholeDir(), name, args...)
}

func (e *Environment) update() {
	e.runCmdInDir("git", "pull")
}

func (e *Environment) CurrentTag() string {
	e.update()
	e.runCmdInDir("git", "checkout", e.Branch)
	b, _, _ := e.runCmdInDir("git", "describe", e.Branch, "--tags")
	return strings.TrimSpace(string(b))
}

func (e *Environment) checkout(tag string) error {
	_, _, err := e.runCmdInDir("git", "checkout", tag)
	return err
}

func (e *Environment) commitAndTag(tag string) {
	// i should return an error, but what do i do when one of these commands fails?? a checkout? what if that fails?
	e.runCmdInDir("git", "add", "-A")
	e.runCmdInDir("git", "commit", "-m", tag)
	e.runCmdInDir("git", "tag", "-f", tag)
}

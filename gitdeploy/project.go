package gitdeploy

import (
	"errors"
	"sync"
)

type Project struct {
	Name                         string
	DefaultSourceEnvironmentName string
	Environments                 map[string]*Environment // maybe it makes sense for this to be a map[string]*Environment?

	baseDirectory string

	deployMu sync.Mutex
}

func (p *Project) Deploy(source, target *Environment, buildNumber string) error {
	if target == nil {
		return errors.New("Target needed")
	}

	// first, try to get default from env
	if source == nil {
		source, _ = p.GetEnvionmentByName(target.DefaultSourceEnvironmentName)
	}

	// then, try to get it from project
	if source == nil {
		source, _ = p.GetEnvionmentByName(p.DefaultSourceEnvironmentName)
	}

	if source == nil {
		return errors.New("Could not find a source environment")
	}

	p.deployMu.Lock()
	defer p.deployMu.Unlock()

	// update source
	source.update()

	if buildNumber == "" {
		// set from source
		buildNumber = source.CurrentTag()
	}

	// checkout build in source
	source.checkout(buildNumber)

	// overlay files (remove everything and copy)
	overlayFiles(source, target)

	// commit target and tag
	target.commitAndTag(buildNumber)

	// push, what if this rejects??? retry?
	target.runCmdInDir("git", "push", "-f")
	target.runCmdInDir("git", "push", "--tags")

	return nil
}

func (s *Project) GetEnvionmentByName(name string) (*Environment, error) {
	for _, e := range s.Environments {
		if e.Name == name {
			return e, nil
		}
	}

	return nil, errors.New("Could not find environment with name " + name)
}

func (p *Project) runCmdInDir(name string, args ...string) ([]byte, int, error) {
	return runCmdInDir(p.baseDirectory, name, args...)
}

func (p *Project) AddEnvironments(es ...*Environment) {
	for _, e := range es {
		// should we return an error here?
		if _, pres := p.Environments[e.Name]; pres {
			continue
		}

		e.project = p
		p.Environments[e.Name] = e
	}
}

func (p *Project) InitEnviroments() {
	for _, e := range p.Environments {
		p.runCmdInDir("mkdir", "-p", e.getDir())

		if _, code, _ := e.runCmdInDir("git", "status"); code == 0 {
			// is a repo, so just update it?
			e.update()
		} else {
			e.runCmdInDir("git", "clone", e.RepoURL, ".")
		}
	}
}

func NewProject(name, baseDirectory string) *Project {
	return &Project{
		Name:          name,
		baseDirectory: baseDirectory,
		Environments:  make(map[string]*Environment),
	}
}

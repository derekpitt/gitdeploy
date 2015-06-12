# Git Deploy

Simple project deployments with git and go.

Use tags to mark builds. Deploy builds from environment to environment with built-in git history.

## Install

`go get github.com/derekpitt/gitdeploy/gitdeploy`

## Configure

    // name of project and base directory to hold enviroment repos
    project := gitdeploy.NewProject("ui", "./projects") 

    auto := &gitdeploy.Environment{
      Name:          "auto",
      RepoURL:       "...",
      Branch:        "gh-pages",
      PreventDeploy: true,
    }

    qa := &gitdeploy.Environment{
      Name:    "qa",
      Branch:  "master",
      RepoURL: "...",
    }

    project.AddEnvironments(auto, qa)

    project.InitEnviroments()


## Deploy

    project.Deploy(auto, qa, "") // deploy latest from auto to qa

## All Together

    package main

    import "github.com/derekpitt/gitdeploy/gitdeploy"

    func main() {
      project := gitdeploy.NewProject("ui", "./projects")

      auto := &gitdeploy.Environment{
        Name:          "auto",
        RepoURL:       "...",
        Branch:        "gh-pages",
        PreventDeploy: true,
      }

      qa := &gitdeploy.Environment{
        Name:    "qa",
        Branch:  "master",
        RepoURL: "...",
      }

      project.AddEnvironments(auto, qa)

      project.InitEnviroments()

      project.Deploy(auto, qa, "")
    }

### TODO

- Some sort of roll-back
- More error checking
- A simple command line tool

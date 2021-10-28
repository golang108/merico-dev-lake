package main // must be main for plugin entry point

import (
	"context"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github-domain/tasks"
	"github.com/mitchellh/mapstructure"
)

type GithubDomainOptions struct {
	Tasks []string `json:"tasks,omitempty"`
}

// plugin interface
type GithubDomain string

func (plugin GithubDomain) Init() {
}

func (plugin GithubDomain) Description() string {
	return "Convert Github Entities to Domain Layer Entities"
}

func (plugin GithubDomain) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	// process options
	var op GithubDomainOptions
	var err error
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return err
	}

	tasksToRun := make(map[string]bool, len(op.Tasks))
	for _, task := range op.Tasks {
		tasksToRun[task] = true
	}
	if len(tasksToRun) == 0 {
		tasksToRun = map[string]bool{
			"convertIssues":  true,
			"convertPrs":     true,
			"convertCommits": true,
		}
	}

	// run tasks
	logger.Print("start GithubDomain plugin execution")

	progress <- 0.01
	if tasksToRun["convertIssues"] {
		err = tasks.ConvertIssues()
		if err != nil {
			return err
		}
	}
	progress <- 0.05
	if tasksToRun["convertPrs"] {
		err = tasks.ConvertPullRequests()
		if err != nil {
			return err
		}
	}
	progress <- 0.07
	if tasksToRun["convertCommits"] {
		err = tasks.ConvertCommits()
		if err != nil {
			return err
		}
	}

	progress <- 1
	logger.Print("end GithubDomain plugin execution")
	close(progress)
	return nil
}

func (plugin GithubDomain) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/github-domain"
}

func (plugin GithubDomain) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return make(map[string]map[string]core.ApiResourceHandler)
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry GithubDomain //nolint

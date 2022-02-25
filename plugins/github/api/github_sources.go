package api

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/mitchellh/mapstructure"
)

// This object conforms to what the frontend currently sends.
type GithubSource struct {
	Endpoint string `mapstructure:"GITHUB_ENDPOINT"`
	Auth     string `mapstructure:"GITHUB_AUTH"`
	Proxy    string `mapstructure:"GITHUB_PROXY"`
}

type GithubConfig struct {
	Endpoint string `mapstructure:"GITHUB_ENDPOINT"`
	Auth     string `mapstructure:"GITHUB_AUTH"`
	Proxy    string `mapstructure:"GITHUB_PROXY"`

	PrType            string `mapstructure:"GITHUB_PR_TYPE"`
	PrComponent       string `mapstructure:"GITHUB_PR_COMPONENT"`
	IssueSeverity     string `mapstructure:"GITHUB_ISSUE_SEVERITY"`
	IssuePriority     string `mapstructure:"GITHUB_ISSUE_PRIORITY"`
	IssueRequirement  string `mapstructure:"GITHUB_ISSUE_REQUIREMENT"`
	IssueCompoent     string `mapstructure:"GITHUB_ISSUE_COMPONENT"`
	IssueTypeBug      string `mapstructure:"GITHUB_ISSUE_TYPE_BUG"`
	IssueTypeIncident string `mapstructure:"GITHUB_ISSUE_TYPE_INCIDENT"`
}

// This object conforms to what the frontend currently expects.
type GithubResponse struct {
	Name  string
	ID    int
	Proxy string `json:"proxy"`

	GithubConfig
}

/*
PUT /plugins/github/sources/:sourceId
*/
func PutSource(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	githubSource := GithubSource{}
	err := mapstructure.Decode(input.Body, &githubSource)
	if err != nil {
		return nil, err
	}
	v := config.GetConfig()
	if githubSource.Endpoint != "" {
		v.Set("GITHUB_ENDPOINT", githubSource.Endpoint)
	}
	if githubSource.Auth != "" {
		v.Set("GITHUB_AUTH", githubSource.Auth)
	}
	v.Set("GITHUB_PROXY", githubSource.Proxy)
	err = v.WriteConfig()
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: "Success"}, nil
}

/*
GET /plugins/github/sources
*/
func ListSources(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-source is developed.
	githubSources, err := GetSourceFromEnv()
	response := []GithubResponse{*githubSources}
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: response}, nil
}

/*
GET /plugins/github/sources/:sourceId
*/
func GetSource(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-source is developed.)
	githubSources, err := GetSourceFromEnv()
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: githubSources}, nil
}

func GetSourceFromEnv() (*GithubResponse, error) {
	v := config.GetConfig()
	var configJson GithubConfig
	err := v.Unmarshal(&configJson)
	if err != nil {
		return nil, err
	}

	return &GithubResponse{
		Name:  "Github",
		ID:    1,
		Proxy: configJson.Proxy,

		GithubConfig: GithubConfig{
			Endpoint:          configJson.Endpoint,
			Auth:              configJson.Auth,
			Proxy:             configJson.Proxy,
			PrType:            configJson.PrType,
			PrComponent:       configJson.PrComponent,
			IssueSeverity:     configJson.IssueSeverity,
			IssuePriority:     configJson.IssuePriority,
			IssueRequirement:  configJson.IssueRequirement,
			IssueCompoent:     configJson.IssueCompoent,
			IssueTypeBug:      configJson.IssueTypeBug,
			IssueTypeIncident: configJson.IssueTypeIncident,
		},
	}, nil
}

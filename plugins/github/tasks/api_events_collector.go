package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/plugins/helper"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
)

const RAW_EVENTS_TABLE = "github_api_events"

// this struct should be moved to `gitub_api_common.go`

var CollectApiEventsMeta = core.SubTaskMeta{
	Name:             "collectApiEvents",
	EntryPoint:       CollectApiEvents,
	EnabledByDefault: true,
	Description:      "Collect Events data from Github api",
}

func CollectApiEvents(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)

	since := data.Since
	incremental := false
	// user didn't specify a time range to sync, try load from database
	// actually, for github pull, since doesn't make any sense, github pull api doesn't support it
	if since == nil {
		var latestUpdatedIssueEvent models.GithubIssueEvent
		err := db.Model(&latestUpdatedIssueEvent).
			Joins("left join _tool_github_issues on _tool_github_issues.github_id = _tool_github_issue_events.issue_id").
			Where("_tool_github_issues.repo_id = ?", data.Repo.GithubId).
			Order("github_created_at DESC").Limit(1).Find(&latestUpdatedIssueEvent).Error
		if err != nil {
			return fmt.Errorf("failed to get latest github issue record: %w", err)
		}

		if latestUpdatedIssueEvent.GithubId > 0 {
			since = &latestUpdatedIssueEvent.GithubCreatedAt
			incremental = true
		}

	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_EVENTS_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/issues/events",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("state", "all")
			if since != nil {
				query.Set("since", since.String())
			}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var items []json.RawMessage
			err := helper.UnmarshalResponse(res, &items)
			if err != nil {
				return nil, err
			}
			return items, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

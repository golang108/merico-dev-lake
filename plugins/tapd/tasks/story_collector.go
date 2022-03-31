package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"net/http"
	"net/url"
)

const RAW_STORY_TABLE = "tapd_api_stories"

var _ core.SubTaskEntryPoint = CollectStorys

func CollectStorys(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect stories")
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_STORY_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "stories",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Stories []json.RawMessage `json:"data"`
			}
			err := helper.UnmarshalResponse(res, &data)
			return data.Stories, err
		},
	})
	if err != nil {
		logger.Error("collect story error:", err)
		return err
	}
	return collector.Execute()
}

var CollectStoryMeta = core.SubTaskMeta{
	Name:        "collectStorys",
	EntryPoint:  CollectStorys,
	Required:    true,
	Description: "collect Tapd stories",
}

/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

const (
	rawEnterpriseAiCreditUsageTable = "copilot_ai_credit_usage_enterprise"
	rawOrgAiCreditUsageTable        = "copilot_ai_credit_usage_org"
)

// CollectEnterpriseAiCreditUsage collects enterprise-level daily AI credit billing usage.
func CollectEnterpriseAiCreditUsage(taskCtx plugin.SubTaskContext) errors.Error {
	return collectAiCreditUsage(taskCtx, models.AiCreditUsageLevelEnterprise)
}

// CollectOrgAiCreditUsage collects organization-level daily AI credit billing usage.
func CollectOrgAiCreditUsage(taskCtx plugin.SubTaskContext) errors.Error {
	return collectAiCreditUsage(taskCtx, models.AiCreditUsageLevelOrganization)
}

// collectAiCreditUsage iterates day-by-day over the GitHub Billing Usage
// "AI credit usage" endpoint, passing year/month/day query parameters, and stores
// each daily report as a raw record. Shared by the enterprise and organization subtasks.
func collectAiCreditUsage(taskCtx plugin.SubTaskContext, level string) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*GhCopilotTaskData)
	if !ok {
		return errors.Default.New("task data is not GhCopilotTaskData")
	}
	connection := data.Connection
	connection.Normalize()
	logger := taskCtx.GetLogger()

	var urlTemplate, rawTable string
	switch level {
	case models.AiCreditUsageLevelEnterprise:
		if !connection.HasEnterprise() {
			logger.Info("No enterprise configured, skipping enterprise AI credit usage collection")
			return nil
		}
		urlTemplate = fmt.Sprintf("enterprises/%s/settings/billing/ai_credit/usage", connection.Enterprise)
		rawTable = rawEnterpriseAiCreditUsageTable
	case models.AiCreditUsageLevelOrganization:
		if connection.Organization == "" {
			logger.Info("No organization configured, skipping org AI credit usage collection")
			return nil
		}
		// Note: the billing usage endpoint uses "organizations/{org}", not "orgs/{org}".
		urlTemplate = fmt.Sprintf("organizations/%s/settings/billing/ai_credit/usage", connection.Organization)
		rawTable = rawOrgAiCreditUsageTable
	default:
		return errors.Default.New(fmt.Sprintf("unknown AI credit usage level: %s", level))
	}

	apiClient, err := CreateApiClient(taskCtx.TaskContext(), connection)
	if err != nil {
		return err
	}

	rawArgs := helper.RawDataSubTaskArgs{
		Ctx:   taskCtx,
		Table: rawTable,
		Options: copilotRawParams{
			ConnectionId: data.Options.ConnectionId,
			ScopeId:      data.Options.ScopeId,
			Organization: connection.Organization,
			Endpoint:     connection.Endpoint,
		},
	}

	collector, err := helper.NewStatefulApiCollector(rawArgs)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	start, until := computeReportDateRange(now, collector.GetSince())
	start = clampDailyMetricsStartForBackfill(start, until)

	dayIter := newDayIterator(start, until)

	err = collector.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   apiClient,
		Input:       dayIter,
		UrlTemplate: urlTemplate,
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			input := reqData.Input.(*dayInput)
			y, m, d, parseErr := splitDay(input.Day)
			if parseErr != nil {
				return nil, parseErr
			}
			q := url.Values{}
			q.Set("year", strconv.Itoa(y))
			q.Set("month", strconv.Itoa(m))
			q.Set("day", strconv.Itoa(d))
			return q, nil
		},
		Incremental:   true,
		Concurrency:   1,
		AfterResponse: ignoreNoContent,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			return parseAiCreditUsageResponse(res, logger)
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()
}

// splitDay parses a YYYY-MM-DD string into year, month and day integers.
func splitDay(day string) (int, int, int, errors.Error) {
	t, parseErr := time.Parse("2006-01-02", day)
	if parseErr != nil {
		return 0, 0, 0, errors.Default.Wrap(parseErr, fmt.Sprintf("invalid day %q", day))
	}
	return t.Year(), int(t.Month()), t.Day(), nil
}

// parseAiCreditUsageResponse reads the single JSON report object returned for one day.
// Empty bodies (no usage) are skipped; the whole object is stored as one raw record
// and exploded into per-usageItem rows by the extractor.
func parseAiCreditUsageResponse(res *http.Response, logger log.Logger) ([]json.RawMessage, errors.Error) {
	body, readErr := io.ReadAll(res.Body)
	res.Body.Close()
	if readErr != nil {
		return nil, errors.Default.Wrap(readErr, "failed to read AI credit usage response")
	}
	if isEmptyReport(body) {
		return nil, nil
	}
	trimmed := bytes.TrimSpace(body)
	// Defensive: only forward objects that actually carry usage items.
	if !strings.Contains(string(trimmed), "usageItems") {
		if logger != nil {
			logger.Info("AI credit usage response had no usageItems, skipping")
		}
		return nil, nil
	}
	return []json.RawMessage{json.RawMessage(trimmed)}, nil
}

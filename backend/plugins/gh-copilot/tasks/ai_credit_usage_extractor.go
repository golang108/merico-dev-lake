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
	"encoding/json"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

// aiCreditUsageReport mirrors the GitHub Billing Usage "AI credit usage" response.
type aiCreditUsageReport struct {
	TimePeriod struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
	} `json:"timePeriod"`
	Enterprise   string `json:"enterprise"`
	Organization string `json:"organization"`
	User         string `json:"user"`
	CostCenter   *struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"costCenter"`
	UsageItems []aiCreditUsageItem `json:"usageItems"`
}

type aiCreditUsageItem struct {
	Product          string  `json:"product"`
	Sku              string  `json:"sku"`
	Model            string  `json:"model"`
	UnitType         string  `json:"unitType"`
	PricePerUnit     float64 `json:"pricePerUnit"`
	GrossQuantity    float64 `json:"grossQuantity"`
	GrossAmount      float64 `json:"grossAmount"`
	DiscountQuantity float64 `json:"discountQuantity"`
	DiscountAmount   float64 `json:"discountAmount"`
	NetQuantity      float64 `json:"netQuantity"`
	NetAmount        float64 `json:"netAmount"`
}

// ExtractEnterpriseAiCreditUsage parses raw enterprise AI credit usage reports.
func ExtractEnterpriseAiCreditUsage(taskCtx plugin.SubTaskContext) errors.Error {
	return extractAiCreditUsage(taskCtx, models.AiCreditUsageLevelEnterprise, rawEnterpriseAiCreditUsageTable)
}

// ExtractOrgAiCreditUsage parses raw organization AI credit usage reports.
func ExtractOrgAiCreditUsage(taskCtx plugin.SubTaskContext) errors.Error {
	return extractAiCreditUsage(taskCtx, models.AiCreditUsageLevelOrganization, rawOrgAiCreditUsageTable)
}

// extractAiCreditUsage explodes each daily report's usageItems into one
// GhCopilotAiCreditUsage row per product/sku/model/unitType. Shared by the
// enterprise and organization extractor subtasks.
func extractAiCreditUsage(taskCtx plugin.SubTaskContext, level, rawTable string) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*GhCopilotTaskData)
	if !ok {
		return errors.Default.New("task data is not GhCopilotTaskData")
	}
	connection := data.Connection
	connection.Normalize()

	params := copilotRawParams{
		ConnectionId: data.Options.ConnectionId,
		ScopeId:      data.Options.ScopeId,
		Organization: connection.Organization,
		Endpoint:     connection.Endpoint,
	}

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Table:   rawTable,
			Options: params,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			rows, err := mapAiCreditUsageReport(level, data.Options.ConnectionId, data.Options.ScopeId, row.Data)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, len(rows))
			for i, r := range rows {
				results[i] = r
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}

// mapAiCreditUsageReport unmarshals a single AI credit usage report and explodes its
// usageItems into one GhCopilotAiCreditUsage row per product/sku/model/unitType.
func mapAiCreditUsageReport(level string, connectionId uint64, scopeId string, raw []byte) ([]*models.GhCopilotAiCreditUsage, errors.Error) {
	var report aiCreditUsageReport
	if err := errors.Convert(json.Unmarshal(raw, &report)); err != nil {
		return nil, err
	}

	tp := report.TimePeriod
	var date time.Time
	if tp.Year > 0 && tp.Month > 0 && tp.Day > 0 {
		date = time.Date(tp.Year, time.Month(tp.Month), tp.Day, 0, 0, 0, 0, time.UTC)
	}

	costCenterId, costCenterName := "", ""
	if report.CostCenter != nil {
		costCenterId = report.CostCenter.Id
		costCenterName = report.CostCenter.Name
	}

	results := make([]*models.GhCopilotAiCreditUsage, 0, len(report.UsageItems))
	for _, item := range report.UsageItems {
		results = append(results, &models.GhCopilotAiCreditUsage{
			ConnectionId:     connectionId,
			ScopeId:          scopeId,
			Level:            level,
			Year:             tp.Year,
			Month:            tp.Month,
			Day:              tp.Day,
			Product:          item.Product,
			Sku:              item.Sku,
			Model:            item.Model,
			UnitType:         item.UnitType,
			CostCenterId:     costCenterId,
			Date:             date,
			Enterprise:       report.Enterprise,
			Organization:     report.Organization,
			UserLogin:        report.User,
			CostCenterName:   costCenterName,
			PricePerUnit:     item.PricePerUnit,
			GrossQuantity:    item.GrossQuantity,
			GrossAmount:      item.GrossAmount,
			DiscountQuantity: item.DiscountQuantity,
			DiscountAmount:   item.DiscountAmount,
			NetQuantity:      item.NetQuantity,
			NetAmount:        item.NetAmount,
		})
	}
	return results, nil
}

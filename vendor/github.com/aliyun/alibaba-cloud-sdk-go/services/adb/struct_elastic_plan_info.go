package adb

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// ElasticPlanInfo is a nested struct in adb response
type ElasticPlanInfo struct {
	PlanName         string `json:"PlanName" xml:"PlanName"`
	ResourcePoolName string `json:"ResourcePoolName" xml:"ResourcePoolName"`
	ElasticNodeNum   int    `json:"ElasticNodeNum" xml:"ElasticNodeNum"`
	StartTime        string `json:"StartTime" xml:"StartTime"`
	EndTime          string `json:"EndTime" xml:"EndTime"`
	WeeklyRepeat     string `json:"WeeklyRepeat" xml:"WeeklyRepeat"`
	StartDay         string `json:"StartDay" xml:"StartDay"`
	EndDay           string `json:"EndDay" xml:"EndDay"`
	Enable           bool   `json:"Enable" xml:"Enable"`
}

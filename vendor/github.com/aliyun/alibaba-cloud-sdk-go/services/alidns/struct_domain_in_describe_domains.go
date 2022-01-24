package alidns

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

// DomainInDescribeDomains is a nested struct in alidns response
type DomainInDescribeDomains struct {
	DomainId        string                      `json:"DomainId" xml:"DomainId"`
	DomainName      string                      `json:"DomainName" xml:"DomainName"`
	PunyCode        string                      `json:"PunyCode" xml:"PunyCode"`
	AliDomain       bool                        `json:"AliDomain" xml:"AliDomain"`
	RecordCount     int64                       `json:"RecordCount" xml:"RecordCount"`
	RegistrantEmail string                      `json:"RegistrantEmail" xml:"RegistrantEmail"`
	Remark          string                      `json:"Remark" xml:"Remark"`
	GroupId         string                      `json:"GroupId" xml:"GroupId"`
	GroupName       string                      `json:"GroupName" xml:"GroupName"`
	InstanceId      string                      `json:"InstanceId" xml:"InstanceId"`
	VersionCode     string                      `json:"VersionCode" xml:"VersionCode"`
	VersionName     string                      `json:"VersionName" xml:"VersionName"`
	InstanceEndTime string                      `json:"InstanceEndTime" xml:"InstanceEndTime"`
	InstanceExpired bool                        `json:"InstanceExpired" xml:"InstanceExpired"`
	Starmark        bool                        `json:"Starmark" xml:"Starmark"`
	CreateTime      string                      `json:"CreateTime" xml:"CreateTime"`
	CreateTimestamp int64                       `json:"CreateTimestamp" xml:"CreateTimestamp"`
	ResourceGroupId string                      `json:"ResourceGroupId" xml:"ResourceGroupId"`
	DnsServers      DnsServersInDescribeDomains `json:"DnsServers" xml:"DnsServers"`
	Tags            TagsInDescribeDomains       `json:"Tags" xml:"Tags"`
}

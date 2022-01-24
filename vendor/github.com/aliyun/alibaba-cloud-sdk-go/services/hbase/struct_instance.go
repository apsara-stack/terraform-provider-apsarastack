package hbase

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

// Instance is a nested struct in hbase response
type Instance struct {
	Engine               string                  `json:"Engine" xml:"Engine"`
	MasterDiskSize       int                     `json:"MasterDiskSize" xml:"MasterDiskSize"`
	CoreDiskCount        string                  `json:"CoreDiskCount" xml:"CoreDiskCount"`
	ExpireTimeUTC        string                  `json:"ExpireTimeUTC" xml:"ExpireTimeUTC"`
	CoreDiskSize         int                     `json:"CoreDiskSize" xml:"CoreDiskSize"`
	CoreNodeCount        int                     `json:"CoreNodeCount" xml:"CoreNodeCount"`
	ModuleStackVersion   string                  `json:"ModuleStackVersion" xml:"ModuleStackVersion"`
	MajorVersion         string                  `json:"MajorVersion" xml:"MajorVersion"`
	DeleteTime           string                  `json:"DeleteTime" xml:"DeleteTime"`
	RegionId             string                  `json:"RegionId" xml:"RegionId"`
	CreatedTime          string                  `json:"CreatedTime" xml:"CreatedTime"`
	ResourceGroupId      string                  `json:"ResourceGroupId" xml:"ResourceGroupId"`
	IsDefault            bool                    `json:"IsDefault" xml:"IsDefault"`
	Duration             int                     `json:"Duration" xml:"Duration"`
	InstanceId           string                  `json:"InstanceId" xml:"InstanceId"`
	CreatedTimeUTC       string                  `json:"CreatedTimeUTC" xml:"CreatedTimeUTC"`
	AutoRenewal          bool                    `json:"AutoRenewal" xml:"AutoRenewal"`
	VswitchId            string                  `json:"VswitchId" xml:"VswitchId"`
	ExpireTime           string                  `json:"ExpireTime" xml:"ExpireTime"`
	ClusterName          string                  `json:"ClusterName" xml:"ClusterName"`
	VpcId                string                  `json:"VpcId" xml:"VpcId"`
	NetworkType          string                  `json:"NetworkType" xml:"NetworkType"`
	IsDeletionProtection bool                    `json:"IsDeletionProtection" xml:"IsDeletionProtection"`
	CoreDiskType         string                  `json:"CoreDiskType" xml:"CoreDiskType"`
	MasterNodeCount      int                     `json:"MasterNodeCount" xml:"MasterNodeCount"`
	MasterInstanceType   string                  `json:"MasterInstanceType" xml:"MasterInstanceType"`
	IsHa                 bool                    `json:"IsHa" xml:"IsHa"`
	ColdStorageStatus    string                  `json:"ColdStorageStatus" xml:"ColdStorageStatus"`
	ClusterId            string                  `json:"ClusterId" xml:"ClusterId"`
	ClusterType          string                  `json:"ClusterType" xml:"ClusterType"`
	ParentId             string                  `json:"ParentId" xml:"ParentId"`
	PayType              string                  `json:"PayType" xml:"PayType"`
	InstanceName         string                  `json:"InstanceName" xml:"InstanceName"`
	ModuleId             int                     `json:"ModuleId" xml:"ModuleId"`
	ZoneId               string                  `json:"ZoneId" xml:"ZoneId"`
	BackupStatus         string                  `json:"BackupStatus" xml:"BackupStatus"`
	CoreInstanceType     string                  `json:"CoreInstanceType" xml:"CoreInstanceType"`
	Status               string                  `json:"Status" xml:"Status"`
	MasterDiskType       string                  `json:"MasterDiskType" xml:"MasterDiskType"`
	Tags                 TagsInDescribeInstances `json:"Tags" xml:"Tags"`
}

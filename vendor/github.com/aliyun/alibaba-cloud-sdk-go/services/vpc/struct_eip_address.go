package vpc

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

// EipAddress is a nested struct in vpc response
type EipAddress struct {
	RegionId                      string                               `json:"RegionId" xml:"RegionId"`
	IpAddress                     string                               `json:"IpAddress" xml:"IpAddress"`
	PrivateIpAddress              string                               `json:"PrivateIpAddress" xml:"PrivateIpAddress"`
	AllocationId                  string                               `json:"AllocationId" xml:"AllocationId"`
	Status                        string                               `json:"Status" xml:"Status"`
	InstanceId                    string                               `json:"InstanceId" xml:"InstanceId"`
	Bandwidth                     string                               `json:"Bandwidth" xml:"Bandwidth"`
	EipBandwidth                  string                               `json:"EipBandwidth" xml:"EipBandwidth"`
	InternetChargeType            string                               `json:"InternetChargeType" xml:"InternetChargeType"`
	AllocationTime                string                               `json:"AllocationTime" xml:"AllocationTime"`
	InstanceType                  string                               `json:"InstanceType" xml:"InstanceType"`
	InstanceRegionId              string                               `json:"InstanceRegionId" xml:"InstanceRegionId"`
	ChargeType                    string                               `json:"ChargeType" xml:"ChargeType"`
	ExpiredTime                   string                               `json:"ExpiredTime" xml:"ExpiredTime"`
	HDMonitorStatus               string                               `json:"HDMonitorStatus" xml:"HDMonitorStatus"`
	Name                          string                               `json:"Name" xml:"Name"`
	ISP                           string                               `json:"ISP" xml:"ISP"`
	Descritpion                   string                               `json:"Descritpion" xml:"Descritpion"`
	BandwidthPackageId            string                               `json:"BandwidthPackageId" xml:"BandwidthPackageId"`
	BandwidthPackageType          string                               `json:"BandwidthPackageType" xml:"BandwidthPackageType"`
	BandwidthPackageBandwidth     string                               `json:"BandwidthPackageBandwidth" xml:"BandwidthPackageBandwidth"`
	ResourceGroupId               string                               `json:"ResourceGroupId" xml:"ResourceGroupId"`
	HasReservationData            string                               `json:"HasReservationData" xml:"HasReservationData"`
	ReservationBandwidth          string                               `json:"ReservationBandwidth" xml:"ReservationBandwidth"`
	ReservationInternetChargeType string                               `json:"ReservationInternetChargeType" xml:"ReservationInternetChargeType"`
	ReservationActiveTime         string                               `json:"ReservationActiveTime" xml:"ReservationActiveTime"`
	ReservationOrderType          string                               `json:"ReservationOrderType" xml:"ReservationOrderType"`
	Mode                          string                               `json:"Mode" xml:"Mode"`
	DeletionProtection            bool                                 `json:"DeletionProtection" xml:"DeletionProtection"`
	SecondLimited                 bool                                 `json:"SecondLimited" xml:"SecondLimited"`
	SegmentInstanceId             string                               `json:"SegmentInstanceId" xml:"SegmentInstanceId"`
	Netmode                       string                               `json:"Netmode" xml:"Netmode"`
	AvailableRegions              AvailableRegions                     `json:"AvailableRegions" xml:"AvailableRegions"`
	OperationLocks                OperationLocksInDescribeEipAddresses `json:"OperationLocks" xml:"OperationLocks"`
	Tags                          TagsInDescribeEipAddresses           `json:"Tags" xml:"Tags"`
}

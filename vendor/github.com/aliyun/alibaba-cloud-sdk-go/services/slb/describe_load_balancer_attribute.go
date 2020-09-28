package slb

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

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeLoadBalancerAttribute invokes the slb.DescribeLoadBalancerAttribute API synchronously
// api document: https://help.aliyun.com/api/slb/describeloadbalancerattribute.html
func (client *Client) DescribeLoadBalancerAttribute(request *DescribeLoadBalancerAttributeRequest) (response *DescribeLoadBalancerAttributeResponse, err error) {
	response = CreateDescribeLoadBalancerAttributeResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeLoadBalancerAttributeWithChan invokes the slb.DescribeLoadBalancerAttribute API asynchronously
// api document: https://help.aliyun.com/api/slb/describeloadbalancerattribute.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLoadBalancerAttributeWithChan(request *DescribeLoadBalancerAttributeRequest) (<-chan *DescribeLoadBalancerAttributeResponse, <-chan error) {
	responseChan := make(chan *DescribeLoadBalancerAttributeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeLoadBalancerAttribute(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DescribeLoadBalancerAttributeWithCallback invokes the slb.DescribeLoadBalancerAttribute API asynchronously
// api document: https://help.aliyun.com/api/slb/describeloadbalancerattribute.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLoadBalancerAttributeWithCallback(request *DescribeLoadBalancerAttributeRequest, callback func(response *DescribeLoadBalancerAttributeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeLoadBalancerAttributeResponse
		var err error
		defer close(result)
		response, err = client.DescribeLoadBalancerAttribute(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DescribeLoadBalancerAttributeRequest is the request struct for api DescribeLoadBalancerAttribute
type DescribeLoadBalancerAttributeRequest struct {
	*requests.RpcRequest
	AccessKeyId          string           `position:"Query" name:"access_key_id"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	IncludeReservedData  requests.Boolean `position:"Query" name:"IncludeReservedData"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	Tags                 string           `position:"Query" name:"Tags"`
	LoadBalancerId       string           `position:"Query" name:"LoadBalancerId"`
}

// DescribeLoadBalancerAttributeResponse is the response struct for api DescribeLoadBalancerAttribute
type DescribeLoadBalancerAttributeResponse struct {
	*responses.BaseResponse
	RequestId                      string                                        `json:"RequestId" xml:"RequestId"`
	LoadBalancerId                 string                                        `json:"LoadBalancerId" xml:"LoadBalancerId"`
	ResourceGroupId                string                                        `json:"ResourceGroupId" xml:"ResourceGroupId"`
	LoadBalancerName               string                                        `json:"LoadBalancerName" xml:"LoadBalancerName"`
	LoadBalancerStatus             string                                        `json:"LoadBalancerStatus" xml:"LoadBalancerStatus"`
	RegionId                       string                                        `json:"RegionId" xml:"RegionId"`
	RegionIdAlias                  string                                        `json:"RegionIdAlias" xml:"RegionIdAlias"`
	Address                        string                                        `json:"Address" xml:"Address"`
	AddressType                    string                                        `json:"AddressType" xml:"AddressType"`
	VpcId                          string                                        `json:"VpcId" xml:"VpcId"`
	VSwitchId                      string                                        `json:"VSwitchId" xml:"VSwitchId"`
	NetworkType                    string                                        `json:"NetworkType" xml:"NetworkType"`
	InternetChargeType             string                                        `json:"InternetChargeType" xml:"InternetChargeType"`
	AutoReleaseTime                int64                                         `json:"AutoReleaseTime" xml:"AutoReleaseTime"`
	Bandwidth                      int                                           `json:"Bandwidth" xml:"Bandwidth"`
	LoadBalancerSpec               string                                        `json:"LoadBalancerSpec" xml:"LoadBalancerSpec"`
	CreateTime                     string                                        `json:"CreateTime" xml:"CreateTime"`
	CreateTimeStamp                int64                                         `json:"CreateTimeStamp" xml:"CreateTimeStamp"`
	EndTime                        string                                        `json:"EndTime" xml:"EndTime"`
	EndTimeStamp                   int64                                         `json:"EndTimeStamp" xml:"EndTimeStamp"`
	PayType                        string                                        `json:"PayType" xml:"PayType"`
	MasterZoneId                   string                                        `json:"MasterZoneId" xml:"MasterZoneId"`
	SlaveZoneId                    string                                        `json:"SlaveZoneId" xml:"SlaveZoneId"`
	AddressIPVersion               string                                        `json:"AddressIPVersion" xml:"AddressIPVersion"`
	CloudType                      string                                        `json:"CloudType" xml:"CloudType"`
	RenewalDuration                int                                           `json:"RenewalDuration" xml:"RenewalDuration"`
	RenewalStatus                  string                                        `json:"RenewalStatus" xml:"RenewalStatus"`
	RenewalCycUnit                 string                                        `json:"RenewalCycUnit" xml:"RenewalCycUnit"`
	HasReservedInfo                string                                        `json:"HasReservedInfo" xml:"HasReservedInfo"`
	ReservedInfoOrderType          string                                        `json:"ReservedInfoOrderType" xml:"ReservedInfoOrderType"`
	ReservedInfoInternetChargeType string                                        `json:"ReservedInfoInternetChargeType" xml:"ReservedInfoInternetChargeType"`
	ReservedInfoBandwidth          string                                        `json:"ReservedInfoBandwidth" xml:"ReservedInfoBandwidth"`
	ReservedInfoActiveTime         string                                        `json:"ReservedInfoActiveTime" xml:"ReservedInfoActiveTime"`
	DeleteProtection               string                                        `json:"DeleteProtection" xml:"DeleteProtection"`
	AssociatedCenId                string                                        `json:"AssociatedCenId" xml:"AssociatedCenId"`
	AssociatedCenStatus            string                                        `json:"AssociatedCenStatus" xml:"AssociatedCenStatus"`
	CloudInstanceType              string                                        `json:"CloudInstanceType" xml:"CloudInstanceType"`
	CloudInstanceId                string                                        `json:"CloudInstanceId" xml:"CloudInstanceId"`
	TunnelType                     string                                        `json:"TunnelType" xml:"TunnelType"`
	CloudInstanceUid               int64                                         `json:"CloudInstanceUid" xml:"CloudInstanceUid"`
	SupportPrivateLink             bool                                          `json:"SupportPrivateLink" xml:"SupportPrivateLink"`
	BusinessStatus                 string                                        `json:"BusinessStatus" xml:"BusinessStatus"`
	ModificationProtectionStatus   string                                        `json:"ModificationProtectionStatus" xml:"ModificationProtectionStatus"`
	ModificationProtectionReason   string                                        `json:"ModificationProtectionReason" xml:"ModificationProtectionReason"`
	ListenerPorts                  ListenerPorts                                 `json:"ListenerPorts" xml:"ListenerPorts"`
	Labels                         Labels                                        `json:"Labels" xml:"Labels"`
	ListenerPortsAndProtocal       ListenerPortsAndProtocal                      `json:"ListenerPortsAndProtocal" xml:"ListenerPortsAndProtocal"`
	ListenerPortsAndProtocol       ListenerPortsAndProtocol                      `json:"ListenerPortsAndProtocol" xml:"ListenerPortsAndProtocol"`
	BackendServers                 BackendServersInDescribeLoadBalancerAttribute `json:"BackendServers" xml:"BackendServers"`
}

// CreateDescribeLoadBalancerAttributeRequest creates a request to invoke DescribeLoadBalancerAttribute API
func CreateDescribeLoadBalancerAttributeRequest() (request *DescribeLoadBalancerAttributeRequest) {
	request = &DescribeLoadBalancerAttributeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Slb", "2014-05-15", "DescribeLoadBalancerAttribute", "slb", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeLoadBalancerAttributeResponse creates a response to parse from DescribeLoadBalancerAttribute response
func CreateDescribeLoadBalancerAttributeResponse() (response *DescribeLoadBalancerAttributeResponse) {
	response = &DescribeLoadBalancerAttributeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

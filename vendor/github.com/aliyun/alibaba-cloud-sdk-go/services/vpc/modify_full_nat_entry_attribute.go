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

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ModifyFullNatEntryAttribute invokes the vpc.ModifyFullNatEntryAttribute API synchronously
func (client *Client) ModifyFullNatEntryAttribute(request *ModifyFullNatEntryAttributeRequest) (response *ModifyFullNatEntryAttributeResponse, err error) {
	response = CreateModifyFullNatEntryAttributeResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyFullNatEntryAttributeWithChan invokes the vpc.ModifyFullNatEntryAttribute API asynchronously
func (client *Client) ModifyFullNatEntryAttributeWithChan(request *ModifyFullNatEntryAttributeRequest) (<-chan *ModifyFullNatEntryAttributeResponse, <-chan error) {
	responseChan := make(chan *ModifyFullNatEntryAttributeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyFullNatEntryAttribute(request)
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

// ModifyFullNatEntryAttributeWithCallback invokes the vpc.ModifyFullNatEntryAttribute API asynchronously
func (client *Client) ModifyFullNatEntryAttributeWithCallback(request *ModifyFullNatEntryAttributeRequest, callback func(response *ModifyFullNatEntryAttributeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyFullNatEntryAttributeResponse
		var err error
		defer close(result)
		response, err = client.ModifyFullNatEntryAttribute(request)
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

// ModifyFullNatEntryAttributeRequest is the request struct for api ModifyFullNatEntryAttribute
type ModifyFullNatEntryAttributeRequest struct {
	*requests.RpcRequest
	FullNatEntryDescription string           `position:"Query" name:"FullNatEntryDescription"`
	ResourceOwnerId         requests.Integer `position:"Query" name:"ResourceOwnerId"`
	AccessIp                string           `position:"Query" name:"AccessIp"`
	ClientToken             string           `position:"Query" name:"ClientToken"`
	FullNatEntryId          string           `position:"Query" name:"FullNatEntryId"`
	NatIpPort               string           `position:"Query" name:"NatIpPort"`
	FullNatTableId          string           `position:"Query" name:"FullNatTableId"`
	AccessPort              string           `position:"Query" name:"AccessPort"`
	DryRun                  requests.Boolean `position:"Query" name:"DryRun"`
	ResourceOwnerAccount    string           `position:"Query" name:"ResourceOwnerAccount"`
	IpProtocol              string           `position:"Query" name:"IpProtocol"`
	OwnerAccount            string           `position:"Query" name:"OwnerAccount"`
	OwnerId                 requests.Integer `position:"Query" name:"OwnerId"`
	FullNatEntryName        string           `position:"Query" name:"FullNatEntryName"`
	NatIp                   string           `position:"Query" name:"NatIp"`
	NetworkInterfaceId      string           `position:"Query" name:"NetworkInterfaceId"`
}

// ModifyFullNatEntryAttributeResponse is the response struct for api ModifyFullNatEntryAttribute
type ModifyFullNatEntryAttributeResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyFullNatEntryAttributeRequest creates a request to invoke ModifyFullNatEntryAttribute API
func CreateModifyFullNatEntryAttributeRequest() (request *ModifyFullNatEntryAttributeRequest) {
	request = &ModifyFullNatEntryAttributeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "ModifyFullNatEntryAttribute", "vpc", "openAPI")
	request.Method = requests.POST
	return
}

// CreateModifyFullNatEntryAttributeResponse creates a response to parse from ModifyFullNatEntryAttribute response
func CreateModifyFullNatEntryAttributeResponse() (response *ModifyFullNatEntryAttributeResponse) {
	response = &ModifyFullNatEntryAttributeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

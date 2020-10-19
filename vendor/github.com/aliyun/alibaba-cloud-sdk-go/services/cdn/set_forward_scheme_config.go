package cdn

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

// SetForwardSchemeConfig invokes the cdn.SetForwardSchemeConfig API synchronously
func (client *Client) SetForwardSchemeConfig(request *SetForwardSchemeConfigRequest) (response *SetForwardSchemeConfigResponse, err error) {
	response = CreateSetForwardSchemeConfigResponse()
	err = client.DoAction(request, response)
	return
}

// SetForwardSchemeConfigWithChan invokes the cdn.SetForwardSchemeConfig API asynchronously
func (client *Client) SetForwardSchemeConfigWithChan(request *SetForwardSchemeConfigRequest) (<-chan *SetForwardSchemeConfigResponse, <-chan error) {
	responseChan := make(chan *SetForwardSchemeConfigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SetForwardSchemeConfig(request)
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

// SetForwardSchemeConfigWithCallback invokes the cdn.SetForwardSchemeConfig API asynchronously
func (client *Client) SetForwardSchemeConfigWithCallback(request *SetForwardSchemeConfigRequest, callback func(response *SetForwardSchemeConfigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SetForwardSchemeConfigResponse
		var err error
		defer close(result)
		response, err = client.SetForwardSchemeConfig(request)
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

// SetForwardSchemeConfigRequest is the request struct for api SetForwardSchemeConfig
type SetForwardSchemeConfigRequest struct {
	*requests.RpcRequest
	SchemeOriginPort string           `position:"Query" name:"SchemeOriginPort"`
	Enable           string           `position:"Query" name:"Enable"`
	SchemeOrigin     string           `position:"Query" name:"SchemeOrigin"`
	DomainName       string           `position:"Query" name:"DomainName"`
	OwnerId          requests.Integer `position:"Query" name:"OwnerId"`
	ConfigId         requests.Integer `position:"Query" name:"ConfigId"`
}

// SetForwardSchemeConfigResponse is the response struct for api SetForwardSchemeConfig
type SetForwardSchemeConfigResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSetForwardSchemeConfigRequest creates a request to invoke SetForwardSchemeConfig API
func CreateSetForwardSchemeConfigRequest() (request *SetForwardSchemeConfigRequest) {
	request = &SetForwardSchemeConfigRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "SetForwardSchemeConfig", "", "")
	request.Method = requests.POST
	return
}

// CreateSetForwardSchemeConfigResponse creates a response to parse from SetForwardSchemeConfig response
func CreateSetForwardSchemeConfigResponse() (response *SetForwardSchemeConfigResponse) {
	response = &SetForwardSchemeConfigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

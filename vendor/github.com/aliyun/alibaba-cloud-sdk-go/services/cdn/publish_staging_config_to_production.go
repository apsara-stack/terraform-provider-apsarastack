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

// PublishStagingConfigToProduction invokes the cdn.PublishStagingConfigToProduction API synchronously
func (client *Client) PublishStagingConfigToProduction(request *PublishStagingConfigToProductionRequest) (response *PublishStagingConfigToProductionResponse, err error) {
	response = CreatePublishStagingConfigToProductionResponse()
	err = client.DoAction(request, response)
	return
}

// PublishStagingConfigToProductionWithChan invokes the cdn.PublishStagingConfigToProduction API asynchronously
func (client *Client) PublishStagingConfigToProductionWithChan(request *PublishStagingConfigToProductionRequest) (<-chan *PublishStagingConfigToProductionResponse, <-chan error) {
	responseChan := make(chan *PublishStagingConfigToProductionResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.PublishStagingConfigToProduction(request)
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

// PublishStagingConfigToProductionWithCallback invokes the cdn.PublishStagingConfigToProduction API asynchronously
func (client *Client) PublishStagingConfigToProductionWithCallback(request *PublishStagingConfigToProductionRequest, callback func(response *PublishStagingConfigToProductionResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *PublishStagingConfigToProductionResponse
		var err error
		defer close(result)
		response, err = client.PublishStagingConfigToProduction(request)
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

// PublishStagingConfigToProductionRequest is the request struct for api PublishStagingConfigToProduction
type PublishStagingConfigToProductionRequest struct {
	*requests.RpcRequest
	FunctionName string           `position:"Query" name:"FunctionName"`
	DomainName   string           `position:"Query" name:"DomainName"`
	OwnerId      requests.Integer `position:"Query" name:"OwnerId"`
}

// PublishStagingConfigToProductionResponse is the response struct for api PublishStagingConfigToProduction
type PublishStagingConfigToProductionResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreatePublishStagingConfigToProductionRequest creates a request to invoke PublishStagingConfigToProduction API
func CreatePublishStagingConfigToProductionRequest() (request *PublishStagingConfigToProductionRequest) {
	request = &PublishStagingConfigToProductionRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "PublishStagingConfigToProduction", "", "")
	request.Method = requests.POST
	return
}

// CreatePublishStagingConfigToProductionResponse creates a response to parse from PublishStagingConfigToProduction response
func CreatePublishStagingConfigToProductionResponse() (response *PublishStagingConfigToProductionResponse) {
	response = &PublishStagingConfigToProductionResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

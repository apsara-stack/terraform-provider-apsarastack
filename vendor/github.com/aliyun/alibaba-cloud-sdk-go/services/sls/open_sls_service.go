package sls

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

// OpenSlsService invokes the sls.OpenSlsService API synchronously
func (client *Client) OpenSlsService(request *OpenSlsServiceRequest) (response *OpenSlsServiceResponse, err error) {
	response = CreateOpenSlsServiceResponse()
	err = client.DoAction(request, response)
	return
}

// OpenSlsServiceWithChan invokes the sls.OpenSlsService API asynchronously
func (client *Client) OpenSlsServiceWithChan(request *OpenSlsServiceRequest) (<-chan *OpenSlsServiceResponse, <-chan error) {
	responseChan := make(chan *OpenSlsServiceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.OpenSlsService(request)
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

// OpenSlsServiceWithCallback invokes the sls.OpenSlsService API asynchronously
func (client *Client) OpenSlsServiceWithCallback(request *OpenSlsServiceRequest, callback func(response *OpenSlsServiceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *OpenSlsServiceResponse
		var err error
		defer close(result)
		response, err = client.OpenSlsService(request)
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

// OpenSlsServiceRequest is the request struct for api OpenSlsService
type OpenSlsServiceRequest struct {
	*requests.RpcRequest
}

// OpenSlsServiceResponse is the response struct for api OpenSlsService
type OpenSlsServiceResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Success   bool   `json:"Success" xml:"Success"`
	Message   string `json:"Message" xml:"Message"`
	Code      string `json:"Code" xml:"Code"`
}

// CreateOpenSlsServiceRequest creates a request to invoke OpenSlsService API
func CreateOpenSlsServiceRequest() (request *OpenSlsServiceRequest) {
	request = &OpenSlsServiceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Sls", "2019-10-23", "OpenSlsService", "", "")
	request.Method = requests.POST
	return
}

// CreateOpenSlsServiceResponse creates a response to parse from OpenSlsService response
func CreateOpenSlsServiceResponse() (response *OpenSlsServiceResponse) {
	response = &OpenSlsServiceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

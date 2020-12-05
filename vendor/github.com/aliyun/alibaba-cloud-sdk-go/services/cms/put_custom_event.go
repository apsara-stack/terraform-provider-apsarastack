package cms

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

// PutCustomEvent invokes the cms.PutCustomEvent API synchronously
func (client *Client) PutCustomEvent(request *PutCustomEventRequest) (response *PutCustomEventResponse, err error) {
	response = CreatePutCustomEventResponse()
	err = client.DoAction(request, response)
	return
}

// PutCustomEventWithChan invokes the cms.PutCustomEvent API asynchronously
func (client *Client) PutCustomEventWithChan(request *PutCustomEventRequest) (<-chan *PutCustomEventResponse, <-chan error) {
	responseChan := make(chan *PutCustomEventResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.PutCustomEvent(request)
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

// PutCustomEventWithCallback invokes the cms.PutCustomEvent API asynchronously
func (client *Client) PutCustomEventWithCallback(request *PutCustomEventRequest, callback func(response *PutCustomEventResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *PutCustomEventResponse
		var err error
		defer close(result)
		response, err = client.PutCustomEvent(request)
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

// PutCustomEventRequest is the request struct for api PutCustomEvent
type PutCustomEventRequest struct {
	*requests.RpcRequest
	EventInfo *[]PutCustomEventEventInfo `position:"Query" name:"EventInfo"  type:"Repeated"`
}

// PutCustomEventEventInfo is a repeated param struct in PutCustomEventRequest
type PutCustomEventEventInfo struct {
	GroupId   string `name:"GroupId"`
	Time      string `name:"Time"`
	EventName string `name:"EventName"`
	Content   string `name:"Content"`
}

// PutCustomEventResponse is the response struct for api PutCustomEvent
type PutCustomEventResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
}

// CreatePutCustomEventRequest creates a request to invoke PutCustomEvent API
func CreatePutCustomEventRequest() (request *PutCustomEventRequest) {
	request = &PutCustomEventRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cms", "2019-01-01", "PutCustomEvent", "cms", "openAPI")
	request.Method = requests.POST
	return
}

// CreatePutCustomEventResponse creates a response to parse from PutCustomEvent response
func CreatePutCustomEventResponse() (response *PutCustomEventResponse) {
	response = &PutCustomEventResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

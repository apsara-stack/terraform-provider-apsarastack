package edas

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

// InsertDegradeControl invokes the edas.InsertDegradeControl API synchronously
func (client *Client) InsertDegradeControl(request *InsertDegradeControlRequest) (response *InsertDegradeControlResponse, err error) {
	response = CreateInsertDegradeControlResponse()
	err = client.DoAction(request, response)
	return
}

// InsertDegradeControlWithChan invokes the edas.InsertDegradeControl API asynchronously
func (client *Client) InsertDegradeControlWithChan(request *InsertDegradeControlRequest) (<-chan *InsertDegradeControlResponse, <-chan error) {
	responseChan := make(chan *InsertDegradeControlResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.InsertDegradeControl(request)
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

// InsertDegradeControlWithCallback invokes the edas.InsertDegradeControl API asynchronously
func (client *Client) InsertDegradeControlWithCallback(request *InsertDegradeControlRequest, callback func(response *InsertDegradeControlResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *InsertDegradeControlResponse
		var err error
		defer close(result)
		response, err = client.InsertDegradeControl(request)
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

// InsertDegradeControlRequest is the request struct for api InsertDegradeControl
type InsertDegradeControlRequest struct {
	*requests.RoaRequest
	Duration    requests.Integer `position:"Query" name:"Duration"`
	RuleType    string           `position:"Query" name:"RuleType"`
	AppId       string           `position:"Query" name:"AppId"`
	UrlVar      string           `position:"Query" name:"UrlVar"`
	RtThreshold requests.Integer `position:"Query" name:"RtThreshold"`
	ServiceName string           `position:"Query" name:"ServiceName"`
	MethodName  string           `position:"Query" name:"MethodName"`
}

// InsertDegradeControlResponse is the response struct for api InsertDegradeControl
type InsertDegradeControlResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateInsertDegradeControlRequest creates a request to invoke InsertDegradeControl API
func CreateInsertDegradeControlRequest() (request *InsertDegradeControlRequest) {
	request = &InsertDegradeControlRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "InsertDegradeControl", "/pop/v5/degradeControl", "Edas", "openAPI")
	request.Method = requests.POST
	return
}

// CreateInsertDegradeControlResponse creates a response to parse from InsertDegradeControl response
func CreateInsertDegradeControlResponse() (response *InsertDegradeControlResponse) {
	response = &InsertDegradeControlResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

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

// GetSpringCloudTestMethod invokes the edas.GetSpringCloudTestMethod API synchronously
func (client *Client) GetSpringCloudTestMethod(request *GetSpringCloudTestMethodRequest) (response *GetSpringCloudTestMethodResponse, err error) {
	response = CreateGetSpringCloudTestMethodResponse()
	err = client.DoAction(request, response)
	return
}

// GetSpringCloudTestMethodWithChan invokes the edas.GetSpringCloudTestMethod API asynchronously
func (client *Client) GetSpringCloudTestMethodWithChan(request *GetSpringCloudTestMethodRequest) (<-chan *GetSpringCloudTestMethodResponse, <-chan error) {
	responseChan := make(chan *GetSpringCloudTestMethodResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetSpringCloudTestMethod(request)
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

// GetSpringCloudTestMethodWithCallback invokes the edas.GetSpringCloudTestMethod API asynchronously
func (client *Client) GetSpringCloudTestMethodWithCallback(request *GetSpringCloudTestMethodRequest, callback func(response *GetSpringCloudTestMethodResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetSpringCloudTestMethodResponse
		var err error
		defer close(result)
		response, err = client.GetSpringCloudTestMethod(request)
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

// GetSpringCloudTestMethodRequest is the request struct for api GetSpringCloudTestMethod
type GetSpringCloudTestMethodRequest struct {
	*requests.RoaRequest
	AppId            string `position:"Query" name:"appId"`
	Namespace        string `position:"Query" name:"namespace"`
	HttpMethod       string `position:"Query" name:"httpMethod"`
	MethodSignature  string `position:"Query" name:"methodSignature"`
	ServiceName      string `position:"Query" name:"serviceName"`
	Region           string `position:"Query" name:"region"`
	RequiredPath     string `position:"Query" name:"requiredPath"`
	MethodController string `position:"Query" name:"methodController"`
}

// GetSpringCloudTestMethodResponse is the response struct for api GetSpringCloudTestMethod
type GetSpringCloudTestMethodResponse struct {
	*responses.BaseResponse
	Code    int    `json:"Code" xml:"Code"`
	Message string `json:"Message" xml:"Message"`
	Success bool   `json:"Success" xml:"Success"`
	Data    Data   `json:"Data" xml:"Data"`
}

// CreateGetSpringCloudTestMethodRequest creates a request to invoke GetSpringCloudTestMethod API
func CreateGetSpringCloudTestMethodRequest() (request *GetSpringCloudTestMethodRequest) {
	request = &GetSpringCloudTestMethodRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "GetSpringCloudTestMethod", "/pop/sp/api/mse/test/springcloud/method", "Edas", "openAPI")
	request.Method = requests.POST
	return
}

// CreateGetSpringCloudTestMethodResponse creates a response to parse from GetSpringCloudTestMethod response
func CreateGetSpringCloudTestMethodResponse() (response *GetSpringCloudTestMethodResponse) {
	response = &GetSpringCloudTestMethodResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

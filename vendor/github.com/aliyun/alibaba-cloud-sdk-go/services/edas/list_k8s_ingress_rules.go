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

// ListK8sIngressRules invokes the edas.ListK8sIngressRules API synchronously
func (client *Client) ListK8sIngressRules(request *ListK8sIngressRulesRequest) (response *ListK8sIngressRulesResponse, err error) {
	response = CreateListK8sIngressRulesResponse()
	err = client.DoAction(request, response)
	return
}

// ListK8sIngressRulesWithChan invokes the edas.ListK8sIngressRules API asynchronously
func (client *Client) ListK8sIngressRulesWithChan(request *ListK8sIngressRulesRequest) (<-chan *ListK8sIngressRulesResponse, <-chan error) {
	responseChan := make(chan *ListK8sIngressRulesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListK8sIngressRules(request)
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

// ListK8sIngressRulesWithCallback invokes the edas.ListK8sIngressRules API asynchronously
func (client *Client) ListK8sIngressRulesWithCallback(request *ListK8sIngressRulesRequest, callback func(response *ListK8sIngressRulesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListK8sIngressRulesResponse
		var err error
		defer close(result)
		response, err = client.ListK8sIngressRules(request)
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

// ListK8sIngressRulesRequest is the request struct for api ListK8sIngressRules
type ListK8sIngressRulesRequest struct {
	*requests.RoaRequest
	Condition string `position:"Query" name:"Condition"`
	Namespace string `position:"Query" name:"Namespace"`
	ClusterId string `position:"Query" name:"ClusterId"`
}

// ListK8sIngressRulesResponse is the response struct for api ListK8sIngressRules
type ListK8sIngressRulesResponse struct {
	*responses.BaseResponse
	Code      int        `json:"Code" xml:"Code"`
	Message   string     `json:"Message" xml:"Message"`
	RequestId string     `json:"RequestId" xml:"RequestId"`
	Data      []DataItem `json:"Data" xml:"Data"`
}

// CreateListK8sIngressRulesRequest creates a request to invoke ListK8sIngressRules API
func CreateListK8sIngressRulesRequest() (request *ListK8sIngressRulesRequest) {
	request = &ListK8sIngressRulesRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "ListK8sIngressRules", "/pop/v5/k8s/acs/k8s_ingress", "Edas", "openAPI")
	request.Method = requests.GET
	return
}

// CreateListK8sIngressRulesResponse creates a response to parse from ListK8sIngressRules response
func CreateListK8sIngressRulesResponse() (response *ListK8sIngressRulesResponse) {
	response = &ListK8sIngressRulesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

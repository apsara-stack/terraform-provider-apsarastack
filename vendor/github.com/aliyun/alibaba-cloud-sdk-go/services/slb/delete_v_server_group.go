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

// DeleteVServerGroup invokes the slb.DeleteVServerGroup API synchronously
func (client *Client) DeleteVServerGroup(request *DeleteVServerGroupRequest) (response *DeleteVServerGroupResponse, err error) {
	response = CreateDeleteVServerGroupResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteVServerGroupWithChan invokes the slb.DeleteVServerGroup API asynchronously
func (client *Client) DeleteVServerGroupWithChan(request *DeleteVServerGroupRequest) (<-chan *DeleteVServerGroupResponse, <-chan error) {
	responseChan := make(chan *DeleteVServerGroupResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteVServerGroup(request)
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

// DeleteVServerGroupWithCallback invokes the slb.DeleteVServerGroup API asynchronously
func (client *Client) DeleteVServerGroupWithCallback(request *DeleteVServerGroupRequest, callback func(response *DeleteVServerGroupResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteVServerGroupResponse
		var err error
		defer close(result)
		response, err = client.DeleteVServerGroup(request)
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

// DeleteVServerGroupRequest is the request struct for api DeleteVServerGroup
type DeleteVServerGroupRequest struct {
	*requests.RpcRequest
	AccessKeyId          string           `position:"Query" name:"access_key_id"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	VServerGroupId       string           `position:"Query" name:"VServerGroupId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	Tags                 string           `position:"Query" name:"Tags"`
}

// DeleteVServerGroupResponse is the response struct for api DeleteVServerGroup
type DeleteVServerGroupResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteVServerGroupRequest creates a request to invoke DeleteVServerGroup API
func CreateDeleteVServerGroupRequest() (request *DeleteVServerGroupRequest) {
	request = &DeleteVServerGroupRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Slb", "2014-05-15", "DeleteVServerGroup", "slb", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDeleteVServerGroupResponse creates a response to parse from DeleteVServerGroup response
func CreateDeleteVServerGroupResponse() (response *DeleteVServerGroupResponse) {
	response = &DeleteVServerGroupResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

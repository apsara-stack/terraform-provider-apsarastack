package bssopenapi

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

// DescribeResourcePackageProduct invokes the bssopenapi.DescribeResourcePackageProduct API synchronously
func (client *Client) DescribeResourcePackageProduct(request *DescribeResourcePackageProductRequest) (response *DescribeResourcePackageProductResponse, err error) {
	response = CreateDescribeResourcePackageProductResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeResourcePackageProductWithChan invokes the bssopenapi.DescribeResourcePackageProduct API asynchronously
func (client *Client) DescribeResourcePackageProductWithChan(request *DescribeResourcePackageProductRequest) (<-chan *DescribeResourcePackageProductResponse, <-chan error) {
	responseChan := make(chan *DescribeResourcePackageProductResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeResourcePackageProduct(request)
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

// DescribeResourcePackageProductWithCallback invokes the bssopenapi.DescribeResourcePackageProduct API asynchronously
func (client *Client) DescribeResourcePackageProductWithCallback(request *DescribeResourcePackageProductRequest, callback func(response *DescribeResourcePackageProductResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeResourcePackageProductResponse
		var err error
		defer close(result)
		response, err = client.DescribeResourcePackageProduct(request)
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

// DescribeResourcePackageProductRequest is the request struct for api DescribeResourcePackageProduct
type DescribeResourcePackageProductRequest struct {
	*requests.RpcRequest
	ProductCode string `position:"Query" name:"ProductCode"`
}

// DescribeResourcePackageProductResponse is the response struct for api DescribeResourcePackageProduct
type DescribeResourcePackageProductResponse struct {
	*responses.BaseResponse
	Code      string                               `json:"Code" xml:"Code"`
	Message   string                               `json:"Message" xml:"Message"`
	RequestId string                               `json:"RequestId" xml:"RequestId"`
	Success   bool                                 `json:"Success" xml:"Success"`
	OrderId   int64                                `json:"OrderId" xml:"OrderId"`
	Data      DataInDescribeResourcePackageProduct `json:"Data" xml:"Data"`
}

// CreateDescribeResourcePackageProductRequest creates a request to invoke DescribeResourcePackageProduct API
func CreateDescribeResourcePackageProductRequest() (request *DescribeResourcePackageProductRequest) {
	request = &DescribeResourcePackageProductRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("BssOpenApi", "2017-12-14", "DescribeResourcePackageProduct", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeResourcePackageProductResponse creates a response to parse from DescribeResourcePackageProduct response
func CreateDescribeResourcePackageProductResponse() (response *DescribeResourcePackageProductResponse) {
	response = &DescribeResourcePackageProductResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

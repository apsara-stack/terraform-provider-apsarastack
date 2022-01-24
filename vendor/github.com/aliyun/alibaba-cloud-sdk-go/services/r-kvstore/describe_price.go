package r_kvstore

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

// DescribePrice invokes the r_kvstore.DescribePrice API synchronously
func (client *Client) DescribePrice(request *DescribePriceRequest) (response *DescribePriceResponse, err error) {
	response = CreateDescribePriceResponse()
	err = client.DoAction(request, response)
	return
}

// DescribePriceWithChan invokes the r_kvstore.DescribePrice API asynchronously
func (client *Client) DescribePriceWithChan(request *DescribePriceRequest) (<-chan *DescribePriceResponse, <-chan error) {
	responseChan := make(chan *DescribePriceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribePrice(request)
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

// DescribePriceWithCallback invokes the r_kvstore.DescribePrice API asynchronously
func (client *Client) DescribePriceWithCallback(request *DescribePriceRequest, callback func(response *DescribePriceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribePriceResponse
		var err error
		defer close(result)
		response, err = client.DescribePrice(request)
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

// DescribePriceRequest is the request struct for api DescribePrice
type DescribePriceRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	NodeType             string           `position:"Query" name:"NodeType"`
	Instances            string           `position:"Query" name:"Instances"`
	CouponNo             string           `position:"Query" name:"CouponNo"`
	InstanceClass        string           `position:"Query" name:"InstanceClass"`
	Capacity             requests.Integer `position:"Query" name:"Capacity"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	BusinessInfo         string           `position:"Query" name:"BusinessInfo"`
	Period               requests.Integer `position:"Query" name:"Period"`
	Product              string           `position:"Query" name:"Product"`
	Quantity             requests.Integer `position:"Query" name:"Quantity"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OrderParamOut        string           `position:"Query" name:"OrderParamOut"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	InstanceId           string           `position:"Query" name:"InstanceId"`
	ZoneId               string           `position:"Query" name:"ZoneId"`
	ChargeType           string           `position:"Query" name:"ChargeType"`
	Category             string           `position:"Query" name:"Category"`
	ForceUpgrade         requests.Boolean `position:"Query" name:"ForceUpgrade"`
	OrderType            string           `position:"Query" name:"OrderType"`
}

// DescribePriceResponse is the response struct for api DescribePrice
type DescribePriceResponse struct {
	*responses.BaseResponse
	RequestId   string    `json:"RequestId" xml:"RequestId"`
	OrderParams string    `json:"OrderParams" xml:"OrderParams"`
	Order       Order     `json:"Order" xml:"Order"`
	Rules       Rules     `json:"Rules" xml:"Rules"`
	SubOrders   SubOrders `json:"SubOrders" xml:"SubOrders"`
}

// CreateDescribePriceRequest creates a request to invoke DescribePrice API
func CreateDescribePriceRequest() (request *DescribePriceRequest) {
	request = &DescribePriceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("R-kvstore", "2015-01-01", "DescribePrice", "redisa", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribePriceResponse creates a response to parse from DescribePrice response
func CreateDescribePriceResponse() (response *DescribePriceResponse) {
	response = &DescribePriceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

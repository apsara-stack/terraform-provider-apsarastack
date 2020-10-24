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

// DescribeDomainCertificateInfo invokes the cdn.DescribeDomainCertificateInfo API synchronously
func (client *Client) DescribeDomainCertificateInfo(request *DescribeDomainCertificateInfoRequest) (response *DescribeDomainCertificateInfoResponse, err error) {
	response = CreateDescribeDomainCertificateInfoResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDomainCertificateInfoWithChan invokes the cdn.DescribeDomainCertificateInfo API asynchronously
func (client *Client) DescribeDomainCertificateInfoWithChan(request *DescribeDomainCertificateInfoRequest) (<-chan *DescribeDomainCertificateInfoResponse, <-chan error) {
	responseChan := make(chan *DescribeDomainCertificateInfoResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDomainCertificateInfo(request)
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

// DescribeDomainCertificateInfoWithCallback invokes the cdn.DescribeDomainCertificateInfo API asynchronously
func (client *Client) DescribeDomainCertificateInfoWithCallback(request *DescribeDomainCertificateInfoRequest, callback func(response *DescribeDomainCertificateInfoResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDomainCertificateInfoResponse
		var err error
		defer close(result)
		response, err = client.DescribeDomainCertificateInfo(request)
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

// DescribeDomainCertificateInfoRequest is the request struct for api DescribeDomainCertificateInfo
type DescribeDomainCertificateInfoRequest struct {
	*requests.RpcRequest
	DomainName string           `position:"Query" name:"DomainName"`
	OwnerId    requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeDomainCertificateInfoResponse is the response struct for api DescribeDomainCertificateInfo
type DescribeDomainCertificateInfoResponse struct {
	*responses.BaseResponse
	RequestId string                                   `json:"RequestId" xml:"RequestId"`
	CertInfos CertInfosInDescribeDomainCertificateInfo `json:"CertInfos" xml:"CertInfos"`
}

// CreateDescribeDomainCertificateInfoRequest creates a request to invoke DescribeDomainCertificateInfo API
func CreateDescribeDomainCertificateInfoRequest() (request *DescribeDomainCertificateInfoRequest) {
	request = &DescribeDomainCertificateInfoRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "DescribeDomainCertificateInfo", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeDomainCertificateInfoResponse creates a response to parse from DescribeDomainCertificateInfo response
func CreateDescribeDomainCertificateInfoResponse() (response *DescribeDomainCertificateInfoResponse) {
	response = &DescribeDomainCertificateInfoResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

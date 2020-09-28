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

// BatchSetCdnDomainServerCertificate invokes the cdn.BatchSetCdnDomainServerCertificate API synchronously
// api document: https://help.aliyun.com/api/cdn/batchsetcdndomainservercertificate.html
func (client *Client) BatchSetCdnDomainServerCertificate(request *BatchSetCdnDomainServerCertificateRequest) (response *BatchSetCdnDomainServerCertificateResponse, err error) {
	response = CreateBatchSetCdnDomainServerCertificateResponse()
	err = client.DoAction(request, response)
	return
}

// BatchSetCdnDomainServerCertificateWithChan invokes the cdn.BatchSetCdnDomainServerCertificate API asynchronously
// api document: https://help.aliyun.com/api/cdn/batchsetcdndomainservercertificate.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) BatchSetCdnDomainServerCertificateWithChan(request *BatchSetCdnDomainServerCertificateRequest) (<-chan *BatchSetCdnDomainServerCertificateResponse, <-chan error) {
	responseChan := make(chan *BatchSetCdnDomainServerCertificateResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.BatchSetCdnDomainServerCertificate(request)
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

// BatchSetCdnDomainServerCertificateWithCallback invokes the cdn.BatchSetCdnDomainServerCertificate API asynchronously
// api document: https://help.aliyun.com/api/cdn/batchsetcdndomainservercertificate.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) BatchSetCdnDomainServerCertificateWithCallback(request *BatchSetCdnDomainServerCertificateRequest, callback func(response *BatchSetCdnDomainServerCertificateResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *BatchSetCdnDomainServerCertificateResponse
		var err error
		defer close(result)
		response, err = client.BatchSetCdnDomainServerCertificate(request)
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

// BatchSetCdnDomainServerCertificateRequest is the request struct for api BatchSetCdnDomainServerCertificate
type BatchSetCdnDomainServerCertificateRequest struct {
	*requests.RpcRequest
	SSLProtocol   string           `position:"Query" name:"SSLProtocol"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	CertType      string           `position:"Query" name:"CertType"`
	SSLPri        string           `position:"Query" name:"SSLPri"`
	ForceSet      string           `position:"Query" name:"ForceSet"`
	CertName      string           `position:"Query" name:"CertName"`
	DomainName    string           `position:"Query" name:"DomainName"`
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	SSLPub        string           `position:"Query" name:"SSLPub"`
	Region        string           `position:"Query" name:"Region"`
}

// BatchSetCdnDomainServerCertificateResponse is the response struct for api BatchSetCdnDomainServerCertificate
type BatchSetCdnDomainServerCertificateResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateBatchSetCdnDomainServerCertificateRequest creates a request to invoke BatchSetCdnDomainServerCertificate API
func CreateBatchSetCdnDomainServerCertificateRequest() (request *BatchSetCdnDomainServerCertificateRequest) {
	request = &BatchSetCdnDomainServerCertificateRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "BatchSetCdnDomainServerCertificate", "", "")
	request.Method = requests.POST
	return
}

// CreateBatchSetCdnDomainServerCertificateResponse creates a response to parse from BatchSetCdnDomainServerCertificate response
func CreateBatchSetCdnDomainServerCertificateResponse() (response *BatchSetCdnDomainServerCertificateResponse) {
	response = &BatchSetCdnDomainServerCertificateResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

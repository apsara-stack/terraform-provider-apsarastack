package hbase

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

// DescribeBackupTables invokes the hbase.DescribeBackupTables API synchronously
func (client *Client) DescribeBackupTables(request *DescribeBackupTablesRequest) (response *DescribeBackupTablesResponse, err error) {
	response = CreateDescribeBackupTablesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeBackupTablesWithChan invokes the hbase.DescribeBackupTables API asynchronously
func (client *Client) DescribeBackupTablesWithChan(request *DescribeBackupTablesRequest) (<-chan *DescribeBackupTablesResponse, <-chan error) {
	responseChan := make(chan *DescribeBackupTablesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeBackupTables(request)
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

// DescribeBackupTablesWithCallback invokes the hbase.DescribeBackupTables API asynchronously
func (client *Client) DescribeBackupTablesWithCallback(request *DescribeBackupTablesRequest, callback func(response *DescribeBackupTablesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeBackupTablesResponse
		var err error
		defer close(result)
		response, err = client.DescribeBackupTables(request)
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

// DescribeBackupTablesRequest is the request struct for api DescribeBackupTables
type DescribeBackupTablesRequest struct {
	*requests.RpcRequest
	PageNumber     requests.Integer `position:"Query" name:"PageNumber"`
	PageSize       requests.Integer `position:"Query" name:"PageSize"`
	BackupRecordId string           `position:"Query" name:"BackupRecordId"`
	ClusterId      string           `position:"Query" name:"ClusterId"`
}

// DescribeBackupTablesResponse is the response struct for api DescribeBackupTables
type DescribeBackupTablesResponse struct {
	*responses.BaseResponse
	RequestId     string                       `json:"RequestId" xml:"RequestId"`
	Total         int64                        `json:"Total" xml:"Total"`
	PageSize      int                          `json:"PageSize" xml:"PageSize"`
	PageNumber    int                          `json:"PageNumber" xml:"PageNumber"`
	Tables        TablesInDescribeBackupTables `json:"Tables" xml:"Tables"`
	BackupRecords BackupRecords                `json:"BackupRecords" xml:"BackupRecords"`
}

// CreateDescribeBackupTablesRequest creates a request to invoke DescribeBackupTables API
func CreateDescribeBackupTablesRequest() (request *DescribeBackupTablesRequest) {
	request = &DescribeBackupTablesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("HBase", "2019-01-01", "DescribeBackupTables", "hbase", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeBackupTablesResponse creates a response to parse from DescribeBackupTables response
func CreateDescribeBackupTablesResponse() (response *DescribeBackupTablesResponse) {
	response = &DescribeBackupTablesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

package ascm

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

type CreateOrganizationRequest struct {
	*requests.RpcRequest
	ClientToken      string `position:"Query" name:"ClientToken"`
	Name             string `position:"Query" name:"Name"`
	PersonNum        string `position:"Query" name:"PersonNum"`
	ParentId         string `position:"Query" name:"ParentId"`
	ResourceGroupNum string `position:"Query" name:"ResourceGroupNum"`
	ResourceGroupId  string `position:"Query" name:"ResourceGroupId"`
}

type CreateOrganizationResponse struct {
	*responses.BaseResponse
	RequestId        string `json:"RequestId" xml:"RequestId"`
	Id               string `json:"Id" xml:"Id"`
	Name             string `json:"Name" xml:"Name"`
	ParentId         string `json:"ParentId" xml:"ParentId"`
	PersonNum        string `json:"PersonNum" xml:"PersonNum"`
	ResourceGroupNum string `json:"ResourceGroupNum" xml:"ResourceGroupNum"`
	ResourceGroupId  string `json:"ResourceGroupId" xml:"ResourceGroupId"`
}

func (client *Client) CreateOrganization(request *CreateOrganizationRequest) (response *CreateOrganizationResponse, err error) {
	response = CreateCreateOrganizationResponse()
	err = client.DoAction(request, response)
	return
}

func CreateCreateOrganizationRequest() (request *CreateOrganizationRequest) {
	request = &CreateOrganizationRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ascm", "2019-05-10", "CreateOrganization", "ascm", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCreateAscmResponse creates a response to parse from CreateOrganization response
func CreateCreateOrganizationResponse() (response *CreateOrganizationResponse) {
	response = &CreateOrganizationResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

package ascm

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) ModifyOrganization(request *ModifyOrganizationRequest) (response *ModifyOrganizationResponse, err error) {
	response = CreateModifyOrganizationResponse()
	err = client.DoAction(request, response)
	return
}

type ModifyOrganizationRequest struct {
	*requests.RpcRequest
	Name             string `position:"Query" name:"Name"`
	PersonNum        string `position:"Query" name:"PersonName"`
	ParentId         string `position:"Query" name:"ParentId"`
	ResourceGroupNum string `position:"Query" name:"ResourceGroupNum"`
	Id               string `position:"Query" name:"Id"`
}

type ModifyOrganizationResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

func CreateModifyOrganizationRequest() (request *ModifyOrganizationRequest) {
	request = &ModifyOrganizationRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ascm", "2019-05-10", "ModifyOrganization", "ascm", "openAPI")
	request.Method = requests.POST
	return
}

func CreateModifyOrganizationResponse() (response *ModifyOrganizationResponse) {
	response = &ModifyOrganizationResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

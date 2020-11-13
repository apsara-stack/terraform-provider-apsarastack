package ascm

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) RemoveOrganization(request *RemoveOrganizationRequest) (response *RemoveOrganizationResponse, err error) {
	response = CreateDeleteOrganizationResponse()
	err = client.DoAction(request, response)
	return
}

type RemoveOrganizationRequest struct {
	*requests.RpcRequest
	Id string `position:"Query" name:"Id"`
}

type RemoveOrganizationResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

func CreateDeleteOrganizationRequest() (request *RemoveOrganizationRequest) {
	request = &RemoveOrganizationRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ascm", "2019-05-10", "DeleteOrganization", "ascm", "openAPI")
	request.Method = requests.POST
	return
}

func CreateDeleteOrganizationResponse() (response *RemoveOrganizationResponse) {
	response = &RemoveOrganizationResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

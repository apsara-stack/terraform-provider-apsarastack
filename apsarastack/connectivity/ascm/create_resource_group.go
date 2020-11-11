package ascm

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) CreateResourceGroup(request *CreateResourceGroupRequest) (response *CreateResourceGroupResponse, err error) {
	response = CreateCreateResourceGroupResponse()
	err = client.DoAction(request, response)
	return
}

type CreateResourceGroupRequest struct {
	*requests.RpcRequest
	ClientToken       string           `position:"Query" name:"ClientToken"`
	ResourceGroupName string           `position:"Query" name:"ResourceGroupName"`
	OrganizationId    requests.Integer `position:"Query" name:"OrganizationId"`
}

type CreateResourceGroupResponse struct {
	*responses.BaseResponse
	RequestId       string `json:"RequestId" xml:"RequestId"`
	ResourceGroupId string `json:"DiskId" xml:"Id"`
}

func CreateCreateResourceGroupRequest() (request *CreateResourceGroupRequest) {
	request = &CreateResourceGroupRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ascm", "2019-05-10", "CreateResourceGroup", "ascm", "openAPI")
	request.Method = requests.POST
	return
}

func CreateCreateResourceGroupResponse() (response *CreateResourceGroupResponse) {
	response = &CreateResourceGroupResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

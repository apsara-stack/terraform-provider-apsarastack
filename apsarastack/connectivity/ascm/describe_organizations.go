package ascm

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) GetOrganizations(request *GetOrganizationsRequest) (response *GetOrganizationsResponse, err error) {
	response = CreateGetOrganizationsResponse()
	err = client.DoAction(request, response)
	return
}

type GetOrganizationsRequest struct {
	*requests.RpcRequest
	PageNumber      requests.Integer `position:"Query" name:"PageNumber"`
	Name            string           `position:"Query" name:"Name"`
	PersonNum       string           `position:"Query" name:"PersonNum"`
	ParentId        string           `position:"Query" name:"ParentId"`
	ResourceGroupId string           `position:"Query" name:"ResourceGroupId"`
	PageSize        requests.Integer `position:"Query" name:"PageSize"`
	Id              string           `position:"Query" name:"Id"`
}

type GetOrganizationsResponse struct {
	*responses.BaseResponse
	RequestId        string        `json:"RequestId" xml:"RequestId"`
	Id               string        `json:"Id" name:"Id"`
	Name             string        `json:"Name" name:"Name"`
	PersonNum        string        `json:"PersonNum" name:"PersonNum"`
	ParentId         string        `json:"ParentId" name:"ParentId"`
	ResourceGroupNum string        `json:"ResourceGroupNum" name:"ResourceGroupNum"`
	PageNumber       int           `json:"PageNumber" xml:"PageNumber"`
	PageSize         int           `json:"PageSize" xml:"PageSize"`
	Status           string        `json:"Status" xml:"Status"`
	Organizations    Organizations `json:"Organizations" xml:"Organizations"`
}

func CreateGetOrganizationsRequest() (request *GetOrganizationsRequest) {
	request = &GetOrganizationsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ascm", "2019-05-10", "GetOrganizations", "ascm", "openAPI")
	request.Method = requests.POST
	return
}
func CreateGetOrganizationsResponse() (response *GetOrganizationsResponse) {
	response = &GetOrganizationsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

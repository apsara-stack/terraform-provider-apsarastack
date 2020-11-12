package ascm

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) ListResourceGroup(request *DescribeResourceGroupsRequest) (response *DescribeResourceGroupsResponse, err error) {
	response = CreateDescribeResourceGroupsResponse()
	err = client.DoAction(request, response)
	return
}

type DescribeResourceGroupsRequest struct {
	*requests.RpcRequest
	PageNumber        requests.Integer `position:"Query" name:"PageNumber"`
	ResourceGroupName string           `position:"Query" name:"ResourceGroupName"`
	OrganizationId    string           `position:"Query" name:"OrganizationId"`
	ResourceGroupId   string           `position:"Query" name:"ResourceGroupId"`
	PageSize          requests.Integer `position:"Query" name:"PageSize"`
	ResourceGroupIds  string           `position:"Query" name:"ResourceGroupIds"`
}

type DescribeResourceGroupsResponse struct {
	*responses.BaseResponse
	RequestId      string                 `json:"RequestId" xml:"RequestId"`
	PageNumber     int                    `json:"PageNumber" xml:"PageNumber"`
	PageSize       int                    `json:"PageSize" xml:"PageSize"`
	Status         string                 `json:"Status" xml:"Status"`
	ResourceGroups DescribeResourceGroups `json:"Organizations" xml:"ResourceGroups"`
}

func CreateDescribeResourceGroupsRequest() (request *DescribeResourceGroupsRequest) {
	request = &DescribeResourceGroupsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ascm", "2019-05-10", "ListResourceGroup", "ascm", "openAPI")
	request.Method = requests.POST
	return
}
func CreateDescribeResourceGroupsResponse() (response *DescribeResourceGroupsResponse) {
	response = &DescribeResourceGroupsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

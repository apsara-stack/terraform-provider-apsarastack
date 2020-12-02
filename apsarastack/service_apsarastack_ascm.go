package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity/ascm"
)

type AscmService struct {
	client *connectivity.ApsaraStackClient
}

func (s *AscmService) DescribeAscmResourceGroup(id string) (response *ResourceG, err error) {
	var requestInfo *ascm.Client
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":          s.client.RegionId,
		"AccessKeySecret":   s.client.SecretKey,
		"Product":           "ascm",
		"Action":            "ListResourceGroup",
		"Version":           "2019-05-10",
		"resourceGroupName": id,
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	request.Scheme = "http"
	request.ApiName = "ListResourceGroup"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &ResourceG{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorResourceGroupNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListResourceGroup", ApsaraStackSdkGoERROR)

	}
	addDebug("ListResourceGroup", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}

type ResourceG struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
		GmtCreated        int64  `json:"gmtCreated"`
		ID                int    `json:"id"`
		OrganizationID    int    `json:"organizationID"`
		OrganizationName  string `json:"organizationName"`
		ResourceGroupName string `json:"resourceGroupName"`
		RsID              string `json:"rsId"`
		Creator           string `json:"creator,omitempty"`
		GmtModified       int64  `json:"gmtModified,omitempty"`
		ResourceGroupType int    `json:"resourceGroupType,omitempty"`
	} `json:"data"`
	Message  string `json:"message"`
	PageInfo struct {
		CurrentPage int   `json:"currentPage"`
		PageSize    int64 `json:"pageSize"`
		Total       int   `json:"total"`
		TotalPage   int   `json:"totalPage"`
	} `json:"pageInfo"`
	PureListData bool `json:"pureListData"`
	Redirect     bool `json:"redirect"`
	Success      bool `json:"success"`
}

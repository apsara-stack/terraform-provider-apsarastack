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
type LoginPolicy struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
		CuserID                string        `json:"cuserId"`
		Default                bool          `json:"default"`
		Enable                 bool          `json:"enable"`
		ID                     int           `json:"id"`
		IPRanges               []interface{} `json:"ipRanges"`
		LpID                   string        `json:"lpId"`
		MuserID                string        `json:"muserId"`
		Name                   string        `json:"name"`
		OrganizationVisibility string        `json:"organizationVisibility"`
		Description            string        `json:"description"`
		OwnerOrganizationID    int           `json:"ownerOrganizationId"`
		Rule                   string        `json:"rule"`
		TimeRanges             []interface{} `json:"timeRanges"`
	} `json:"data"`
	Message      string `json:"message"`
	PureListData bool   `json:"pureListData"`
	Redirect     bool   `json:"redirect"`
	Success      bool   `json:"success"`
}

func (s *AscmService) ListLoginPolicies(id string) (response *LoginPolicy, err error) {
	var requestInfo *ascm.Client
	request := requests.NewCommonRequest()
	request.Method = "POST"  // Set request method
	request.Product = "ascm" // Specify product
	//request.Domain = endpoint          // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2019-05-10" // Specify product version
	request.Scheme = "http"        // Set request scheme. Default: http
	request.ApiName = "ListLoginPolicies"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "ascm",
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"RegionId":        s.client.RegionId,
		"Action":          "ListLoginPolicies",
		"Version":         "2019-05-10",
		//"ParentId": "17",
		"Name": id,
		//"Id":"54438",
	}
	var resp = &LoginPolicy{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorResourceGroupNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListLoginPolicy", ApsaraStackSdkGoERROR)

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
	return response, nil
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

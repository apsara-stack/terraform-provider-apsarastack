package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity/ascm"
	"strings"
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

func (s *AscmService) DescribeAscmUser(id string) (response *User, err error) {
	var requestInfo *ascm.Client
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "ascm",
		"Action":          "ListUsers",
		"Version":         "2019-05-10",
		"loginName":       id,
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	request.Scheme = "http"
	request.ApiName = "ListUsers"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &User{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListUsers", ApsaraStackSdkGoERROR)

	}
	addDebug("ListUsers", response, requestInfo, request)

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

func (s *AscmService) DescribeAscmOrganization(id string) (response *ascm.Organization, err error) {
	var requestInfo *ascm.Client
	did := strings.Split(id, SLASH_SEPARATED)
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "GetOrganizationList",
		"Version":         "2019-05-10",
		"name":            did[0],
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	request.Scheme = "http"
	request.ApiName = "GetOrganizationList"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &ascm.Organization{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorOrganizationNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "GetOrganization", ApsaraStackSdkGoERROR)

	}
	addDebug("GetOrganization", response, requestInfo, request)

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

type User struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
		AccessKeys []struct {
			AccesskeyID string `json:"accesskeyId"`
			Ctime       int64  `json:"ctime"`
			CuserID     string `json:"cuserId"`
			ID          int    `json:"id"`
			Region      string `json:"region"`
			Status      string `json:"status"`
		} `json:"accessKeys"`
		CellphoneNum string `json:"cellphoneNum"`
		Default      bool   `json:"default"`
		DefaultRole  struct {
			Active                 bool   `json:"active"`
			ArID                   string `json:"arId"`
			Code                   string `json:"code"`
			CuserID                string `json:"cuserId"`
			Default                bool   `json:"default"`
			Description            string `json:"description"`
			Enable                 bool   `json:"enable"`
			ID                     int    `json:"id"`
			MuserID                string `json:"muserId"`
			OrganizationVisibility string `json:"organizationVisibility"`
			OwnerOrganizationID    int    `json:"ownerOrganizationId"`
			RAMRole                bool   `json:"rAMRole"`
			RoleLevel              int64  `json:"roleLevel"`
			RoleName               string `json:"roleName"`
			RoleRange              string `json:"roleRange"`
			RoleType               string `json:"roleType"`
		} `json:"defaultRole"`
		Deleted            bool   `json:"deleted"`
		DisplayName        string `json:"displayName"`
		Email              string `json:"email"`
		EnableDingTalk     bool   `json:"enableDingTalk"`
		EnableEmail        bool   `json:"enableEmail"`
		EnableShortMessage bool   `json:"enableShortMessage"`
		ID                 int    `json:"id"`
		LastLoginTime      int64  `json:"lastLoginTime"`
		LoginName          string `json:"loginName"`
		LoginPolicy        struct {
			CuserID  string `json:"cuserId"`
			Default  bool   `json:"default"`
			Enable   bool   `json:"enable"`
			ID       int    `json:"id"`
			IPRanges []struct {
				IPRange       string `json:"ipRange"`
				LoginPolicyID int    `json:"loginPolicyId"`
				Protocol      string `json:"protocol"`
			} `json:"ipRanges"`
			LpID                   string `json:"lpId"`
			MuserID                string `json:"muserId"`
			Name                   string `json:"name"`
			OrganizationVisibility string `json:"organizationVisibility"`
			OwnerOrganizationID    int    `json:"ownerOrganizationId"`
			Rule                   string `json:"rule"`
			TimeRanges             []struct {
				EndTime       string `json:"endTime"`
				LoginPolicyID int    `json:"loginPolicyId"`
				StartTime     string `json:"startTime"`
			} `json:"timeRanges"`
		} `json:"loginPolicy"`
		MobileNationCode string `json:"mobileNationCode"`
		Organization     struct {
			Alias             string        `json:"alias"`
			Ctime             int64         `json:"ctime"`
			CuserID           string        `json:"cuserId"`
			ID                int           `json:"id"`
			Internal          bool          `json:"internal"`
			Level             string        `json:"level"`
			Mtime             int64         `json:"mtime"`
			MultiCloudStatus  string        `json:"multiCloudStatus"`
			MuserID           string        `json:"muserId"`
			Name              string        `json:"name"`
			ParentID          int           `json:"parentId"`
			SupportRegionList []interface{} `json:"supportRegionList"`
			UUID              string        `json:"uuid"`
		} `json:"organization,omitempty"`
		ParentPk   string `json:"parentPk"`
		PrimaryKey string `json:"primaryKey"`
		Roles      []struct {
			Active                 bool   `json:"active"`
			ArID                   string `json:"arId"`
			Code                   string `json:"code"`
			Default                bool   `json:"default"`
			Description            string `json:"description"`
			Enable                 bool   `json:"enable"`
			ID                     int    `json:"id"`
			OrganizationVisibility string `json:"organizationVisibility"`
			OwnerOrganizationID    int    `json:"ownerOrganizationId"`
			RAMRole                bool   `json:"rAMRole"`
			RoleLevel              int64  `json:"roleLevel"`
			RoleName               string `json:"roleName"`
			RoleRange              string `json:"roleRange"`
			RoleType               string `json:"roleType"`
		} `json:"roles"`
		Status         string        `json:"status"`
		UserGroupRoles []interface{} `json:"userGroupRoles"`
		UserGroups     []interface{} `json:"userGroups"`
		UserRoles      []struct {
			Active                 bool   `json:"active"`
			ArID                   string `json:"arId"`
			Code                   string `json:"code"`
			CuserID                string `json:"cuserId"`
			Default                bool   `json:"default"`
			Description            string `json:"description"`
			Enable                 bool   `json:"enable"`
			ID                     int    `json:"id"`
			MuserID                string `json:"muserId"`
			OrganizationVisibility string `json:"organizationVisibility"`
			OwnerOrganizationID    int    `json:"ownerOrganizationId"`
			RAMRole                bool   `json:"rAMRole"`
			RoleLevel              int64  `json:"roleLevel"`
			RoleName               string `json:"roleName"`
			RoleRange              string `json:"roleRange"`
			RoleType               string `json:"roleType"`
		} `json:"userRoles"`
	} `json:"data"`
	Message  string `json:"message"`
	PageInfo struct {
		CurrentPage int `json:"currentPage"`
		PageSize    int `json:"pageSize"`
		Total       int `json:"total"`
		TotalPage   int `json:"totalPage"`
	} `json:"pageInfo"`
	PureListData bool `json:"pureListData"`
	Redirect     bool `json:"redirect"`
	Success      bool `json:"success"`
}

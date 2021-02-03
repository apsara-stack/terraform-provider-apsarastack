package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"strings"
)

type AscmService struct {
	client *connectivity.ApsaraStackClient
}

func (s *AscmService) DescribeAscmLogonPolicy(id string) (response *LoginPolicy, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
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
		"Name":            id,
	}
	var resp = &LoginPolicy{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorLoginPolicyNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListLoginPolicy", ApsaraStackSdkGoERROR)

	}
	addDebug("LoginPolicy", response, requestInfo, request)

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
func (s *AscmService) DescribeAscmResourceGroup(id string) (response *ResourceGroup, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)

	request := requests.NewCommonRequest()

	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":          s.client.RegionId,
		"AccessKeySecret":   s.client.SecretKey,
		"Product":           "ascm",
		"Action":            "ListResourceGroup",
		"Version":           "2019-05-10",
		"resourceGroupName": did[0],
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListResourceGroup"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &ResourceGroup{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorResourceGroupNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, did[0], "ListResourceGroup", ApsaraStackSdkGoERROR)

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
func (s *AscmService) DescribeAscmRamRole(id string) (response *AscmRoles, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)

	request := requests.NewCommonRequest()

	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "ListRoles",
		"Version":         "2019-05-10",
		"roleName":        did[0],
		"roleType":        "ROLETYPE_RAM",
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListRoles"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &AscmRoles{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorRamRoleNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListRoles", ApsaraStackSdkGoERROR)

	}
	addDebug("ListRoles", response, requestInfo, request)

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

func (s *AscmService) DescribeAscmRamServiceRole(id string) (response *RamRole, err error) {
	var requestInfo *ecs.Client

	request := requests.NewCommonRequest()

	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"id":              id,
		"Action":          "GetRAMServiceRole",
		"Version":         "2019-05-10",
		"roleType":        "ROLETYPE_RAM",
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListRAMServiceRoles"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &RamRole{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorRamServiceRoleNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListRAMServiceRoles", ApsaraStackSdkGoERROR)

	}
	addDebug("ListRAMServiceRoles", response, requestInfo, request)

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

func (s *AscmService) DescribeAscmResourceGroupUserAttachment(id string) (response *User, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()

	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "ascm",
		"Action":          "ListAscmUsersInsideResourceGroup",
		"Version":         "2019-05-10",
		"ResourceGroupId": id,
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListAscmUsersInsideResourceGroup"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &User{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorListAscmUsersInsideResourceGroupNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListAscmUsersInsideResourceGroup", ApsaraStackSdkGoERROR)

	}
	addDebug("ListAscmUsersInsideResourceGroup", response, requestInfo, request)

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
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
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
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
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
func (s *AscmService) DescribeAscmUserRoleBinding(id string) (response *User, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
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
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
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
func (s *AscmService) DescribeAscmDeletedUser(id string) (response *DeletedUser, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "ascm",
		"Action":          "ListDeletedUsers",
		"Version":         "2019-05-10",
		"loginName":       id,
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListDeletedUsers"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &DeletedUser{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListDeletedUsers", ApsaraStackSdkGoERROR)

	}
	addDebug("ListDeletedUsers", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}
	if resp.Data != nil {
		return resp, WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmOrganization(id string) (response *Organization, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
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
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "GetOrganizationList"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &Organization{}
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

	if resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}
func (s *AscmService) DescribeAscmRamPolicy(id string) (response *RamPolicies, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeyId":     s.client.AccessKey,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "ListRAMPolicies",
		"Version":         "2019-05-10",
		"policyName":      did[0],
	}
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListRAMPolicies"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &RamPolicies{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorRamPolicyNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListRAMPolicies", ApsaraStackSdkGoERROR)

	}
	addDebug("ListRAMPolicies", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}
func (s *AscmService) DescribeAscmQuota(id string) (response *AscmQuota, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)
	var targetType string
	if did[0] == "RDS" {
		targetType = "MySql"
	} else if did[0] == "R-KVSTORE" {
		targetType = "redis"
	} else if did[0] == "DDS" {
		targetType = "mongodb"
	} else {
		targetType = ""
	}
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"regionName ":     s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "GetQuota",
		"Version":         "2019-05-10",
		"productName":     did[0],
		"quotaType":       did[1],
		"quotaTypeId":     did[2],
		"targetType":      targetType,
	}
	request.Method = "GET"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.ApiName = "GetQuota"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &AscmQuota{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorQuotaNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, did[0], "GetQuota", ApsaraStackSdkGoERROR)

	}
	addDebug("GetQuota", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}
	if resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmPasswordPolicy(id string) (response *PasswordPolicy, err error) {
	var requestInfo *ecs.Client
	//	did := strings.Split(id, COLON_SEPARATED)
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "GetPasswordPolicy",
		"Version":         "2019-05-10",
		"id":              id,
	}
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "GetPasswordPolicy"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	request.Domain = s.client.Domain
	var resp = &PasswordPolicy{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorOrganizationNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "GetPasswordPolicy", ApsaraStackSdkGoERROR)

	}
	addDebug("GetPasswordPolicy", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}

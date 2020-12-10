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

func (s *AscmService) ListLoginPolicies(id string) (response *LoginPolicy, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.Scheme = "http"
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
func (s *AscmService) DescribeAscmResourceGroup(id string) (response *ResourceGroup, err error) {
	var requestInfo *ecs.Client
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
	var resp = &ResourceGroup{}
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
	var requestInfo *ecs.Client
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

func (s *AscmService) DescribeAscmOrganization(id string) (response *Organization, err error) {
	var requestInfo *ecs.Client
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

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}

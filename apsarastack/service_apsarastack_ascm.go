package apsarastack

import (
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
)

type AscmService struct {
	client *connectivity.ApsaraStackClient
}

//func (s *AscmService) GetOrganization(id string) (org ascm.Organization, err error) {
//	request := ascm.CreateGetOrganizationsRequest()
//	request.Id = convertListToJsonString([]interface{}{id})
//	request.RegionId = s.client.RegionId
//	request.Headers = map[string]string{"RegionId": s.client.RegionId}
//	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ascm", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
//	raw, err := s.client.WithAscmClient(func(ascmClient *ascm.Client) (interface{}, error) {
//		return ascmClient.GetOrganizations(request)
//	})
//	if err != nil {
//		return org, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
//	}
//	response, _ := raw.(*ascm.GetOrganizationsResponse)
//	if len(response.Organizations.Organization) < 1 || response.Organizations.Organization[0].Data[] != id {
//		err = WrapErrorf(Error(GetNotFoundMessage("Organization", id)), NotFoundMsg, ProviderERROR, response.RequestId)
//		return
//	}
//	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
//	return response.Organizations.Organization[0], nil
//}

//func (s *AscmService) AscmStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
//	return func() (interface{}, string, error) {
//		object, err := s.GetOrganization(id)
//		if err != nil {
//			if NotFoundError(err) {
//				return nil, "", nil
//			}
//			return nil, "", WrapError(err)
//		}
//
//		for _, failState := range failStates {
//			if object.Status == failState {
//				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
//			}
//		}
//
//		return object, object.Status, nil
//	}
//}

//func (s *AscmService) WaitForOrganization(id string, status Status, timeout int) error {
//	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
//	for {
//		object, err := s.GetOrganization(id)
//		if err != nil {
//			if NotFoundError(err) {
//				if status == Deleted {
//					return nil
//				}
//			} else {
//				return WrapError(err)
//			}
//		}
//		if object.Status == string(status) {
//			return nil
//		}
//		if time.Now().After(deadline) {
//			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
//		}
//		time.Sleep(DefaultIntervalShort * time.Second)
//	}
//}

//func (s *AscmService) DescribeResourceGroup(id string) (rg ascm.ResourceGroup, err error) {
//	request := ascm.CreateDescribeResourceGroupsRequest()
//	request.ResourceGroupIds = convertListToJsonString([]interface{}{id})
//	request.RegionId = s.client.RegionId
//	request.Headers = map[string]string{"RegionId": s.client.RegionId}
//	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ascm", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
//	raw, err := s.client.WithAscmClient(func(ascmClient *ascm.Client) (interface{}, error) {
//		return ascmClient.ListResourceGroup(request)
//	})
//	if err != nil {
//		return rg, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
//	}
//	response, _ := raw.(*ascm.DescribeResourceGroupsResponse)
//	if len(response.ResourceGroups.ResourceGroup) < 1 {
//		err = WrapErrorf(Error(GetNotFoundMessage("ResourceGroup", id)), NotFoundMsg, ProviderERROR, response.RequestId)
//		return
//	}
//	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
//	return response.ResourceGroups.ResourceGroup[0], nil
//}

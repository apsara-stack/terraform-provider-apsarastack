package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	//"github.com/coreos/etcd/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"reflect"
	"strconv"
	"time"
)

type DnsService struct {
	client *connectivity.ApsaraStackClient
}

func (dns *DnsService) DescribeDnsRecord(id string) (*alidns.DescribeDomainRecordInfoResponse, error) {
	response := &alidns.DescribeDomainRecordInfoResponse{}
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.Headers = map[string]string{"RegionId": dns.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": dns.client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = dns.client.Department
	request.QueryParams["ResourceGroup"] = dns.client.ResourceGroup
	request.RecordId = id
	request.RegionId = dns.client.RegionId
	raw, err := dns.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.DescribeDomainRecordInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser", "InvalidRR.NoExist"}) {
			return response, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*alidns.DescribeDomainRecordInfoResponse)
	if response.RecordId != id {
		return response, WrapErrorf(Error(GetNotFoundMessage("DnsRecord", id)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}

func (dns *DnsService) DescribeDnsGroup(id string) (alidns.DomainGroup, error) {
	var group alidns.DomainGroup
	request := alidns.CreateDescribeDomainGroupsRequest()
	request.Headers = map[string]string{"RegionId": dns.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": dns.client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = dns.client.Department
	request.QueryParams["ResourceGroup"] = dns.client.ResourceGroup
	request.RegionId = dns.client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := dns.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			return group, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alidns.DescribeDomainGroupsResponse)
		groups := response.DomainGroups.DomainGroup
		for _, domainGroup := range groups {
			if domainGroup.GroupId == id {
				return domainGroup, nil
			}
		}
		if len(groups) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return group, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return group, WrapErrorf(Error(GetNotFoundMessage("DnsGroup", id)), NotFoundMsg, ProviderERROR)
}

func (s *DnsService) ListTagResources(id string) (object alidns.ListTagResourcesResponse, err error) {
	request := alidns.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = s.client.Department
	request.QueryParams["ResourceGroup"] = s.client.ResourceGroup

	request.ResourceType = "DOMAIN"
	request.ResourceId = &[]string{id}

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.ListTagResources(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.ListTagResourcesResponse)
	return *response, nil
}
func (s *DnsService) DescribeDnsDomainAttachment(id string) (object alidns.DescribeInstanceDomainsResponse, err error) {
	request := alidns.CreateDescribeInstanceDomainsRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = s.client.Department
	request.QueryParams["ResourceGroup"] = s.client.ResourceGroup

	request.InstanceId = id

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeInstanceDomains(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DnsDomainAttachment", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeInstanceDomainsResponse)

	if len(response.InstanceDomains) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("DnsDomainAttachment", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return *response, nil
}

func (s *DnsService) WaitForAlidnsDomainAttachment(id string, expected map[string]interface{}, isDelete bool, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDnsDomainAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if isDelete {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		domainNames := make(map[string]interface{}, 0)
		for _, v := range object.InstanceDomains {
			domainNames[v.DomainName] = v.DomainName
		}

		exceptDomainNames := make(map[string]interface{}, 0)
		for _, v := range expected {
			for _, vv := range v.([]interface{}) {
				exceptDomainNames[vv.(string)] = vv.(string)
			}
		}

		if reflect.DeepEqual(domainNames, exceptDomainNames) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", expected, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}
func (s *DnsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]alidns.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, alidns.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := alidns.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.Headers = map[string]string{"RegionId": s.client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "alidns"}
		request.QueryParams["Department"] = s.client.Department
		request.QueryParams["ResourceGroup"] = s.client.ResourceGroup

		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := alidns.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.Headers = map[string]string{"RegionId": s.client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "alidns"}
		request.QueryParams["Department"] = s.client.Department
		request.QueryParams["ResourceGroup"] = s.client.ResourceGroup

		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
	}
	return nil
}

type Domain struct {
	DomainId        string `json:"DomainId" xml:"DomainId"`
	DomainName      string `json:"DomainName" xml:"DomainName"`
	PunyCode        string `json:"PunyCode" xml:"PunyCode"`
	AliDomain       bool   `json:"AliDomain" xml:"AliDomain"`
	RecordCount     int64  `json:"RecordCount" xml:"RecordCount"`
	RegistrantEmail string `json:"RegistrantEmail" xml:"RegistrantEmail"`
	Remark          string `json:"Remark" xml:"Remark"`
	GroupId         string `json:"GroupId" xml:"GroupId"`
	GroupName       string `json:"GroupName" xml:"GroupName"`
	InstanceId      string `json:"InstanceId" xml:"InstanceId"`
	VersionCode     string `json:"VersionCode" xml:"VersionCode"`
	VersionName     string `json:"VersionName" xml:"VersionName"`
	InstanceEndTime string `json:"InstanceEndTime" xml:"InstanceEndTime"`
	InstanceExpired bool   `json:"InstanceExpired" xml:"InstanceExpired"`
	Starmark        bool   `json:"Starmark" xml:"Starmark"`
	CreateTime      string `json:"CreateTime" xml:"CreateTime"`
	CreateTimestamp int64  `json:"CreateTimestamp" xml:"CreateTimestamp"`
	ResourceGroupId string `json:"ResourceGroupId" xml:"ResourceGroupId"`
	//DnsServers      DnsServersInDescribeDomains `json:"DnsServers" xml:"DnsServers"`
	//Tags            TagsInDescribeDomains       `json:"Tags" xml:"Tags"`
}
type DnsDomains struct {
	RequestID       string    `json:"RequestId"`
	PageSize        int       `json:"PageSize"`
	PageNumber      int       `json:"PageNumber"`
	TotalItems      int       `json:"TotalItems"`
	DomainID        int       `json:"DomainId"`
	DomainName      string    `json:"DomainName"`
	VpcNumber       int       `json:"VpcNumber"`
	CreateTime      time.Time `json:"CreateTime"`
	UpdateTime      time.Time `json:"UpdateTime"`
	UpdateTimestamp int64     `json:"UpdateTimestamp"`
	RecordCount     int       `json:"RecordCount"`
	CreateTimestamp int64     `json:"CreateTimestamp"`
	ZoneList        []struct {
		DomainID        int       `json:"DomainId"`
		DomainName      string    `json:"DomainName"`
		VpcNumber       int       `json:"VpcNumber"`
		CreateTime      time.Time `json:"CreateTime"`
		UpdateTime      time.Time `json:"UpdateTime"`
		UpdateTimestamp int64     `json:"UpdateTimestamp"`
		RecordCount     int       `json:"RecordCount"`
		CreateTimestamp int64     `json:"CreateTimestamp"`
	} `json:"ZoneList"`
}

func (s *DnsService) DescribeDnsDomain(id string) (response *DnsDomains, err error) {
	var domain = &responses.CommonResponse{}
	request := requests.NewCommonRequest()
	request.Method = "POST"          // Set request method
	request.Product = "GenesisDns"   // Specify product
	request.Domain = s.client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2018-07-20"   // Specify product version
	request.Scheme = "http"          // Set request scheme. Default: http
	request.ApiName = "ObtainGlobalAuthZoneList"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": s.client.SecretKey,
		"AccessKeyId":     s.client.AccessKey,
		"Product":         "GenesisDns",
		"RegionId":        s.client.RegionId,
		"Action":          "ObtainGlobalAuthZoneList",
		"Version":         "2018-07-20",
		"Id":              id,
	}
	raw, err := s.client.WithEcsClient(func(alidnsClient *ecs.Client) (interface{}, error) {
		return alidnsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DnsDomain", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
		return
	}
	var dnsdomain = DnsDomains{}
	domain, _ = raw.(*responses.CommonResponse)
	_ = json.Unmarshal(domain.GetHttpContentBytes(), &dnsdomain)

	ID, _ := strconv.Atoi(id)
	dm := &DnsDomains{}
	for _, k := range dnsdomain.ZoneList {
		if k.DomainID == ID {
			dm.DomainID = k.DomainID
			dm.DomainName = k.DomainName
			dm.CreateTime = k.CreateTime
			dm.UpdateTime = k.UpdateTime
		}
	}
	addDebug(request.GetActionName(), raw, response, request)
	return dm, nil
}

package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"strings"
)

func resourceApsaraStackDnsDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDnsDomainCreate,
		Read:   resourceApsaraStackDnsDomainRead,
		Update: resourceApsaraStackDnsDomainUpdate,
		Delete: resourceApsaraStackDnsDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"dns_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApsaraStackDnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	DomainName := d.Get("domain_name").(string)
	request := requests.NewCommonRequest()
	request.Method = "POST"        // Set request method
	request.Product = "GenesisDns" // Specify product
	request.Domain = client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2018-07-20" // Specify product version
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "AddGlobalAuthZone"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "GenesisDns",
		"RegionId":        client.RegionId,
		"Action":          "AddGlobalAuthZone",
		"Version":         "2018-07-20",
		"DomainName":      DomainName,
	}
	raw, err := client.WithEcsClient(func(dnsClient *ecs.Client) (interface{}, error) {
		return dnsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_domain", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	dnsresp := DnsDomain{}
	response, _ := raw.(*responses.CommonResponse)
	ok := json.Unmarshal(response.GetHttpContentBytes(), &dnsresp)
	if ok != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_domain", "AddGlobalAuthZone", raw)
	}
	id := strconv.Itoa(dnsresp.ID)
	d.SetId(id)

	return resourceApsaraStackDnsDomainUpdate(d, meta)
}
func resourceApsaraStackDnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}
	object, err := dnsService.DescribeDnsDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
		}
		return WrapError(err)
	}

	d.Set("domain_name", object.ZoneList[0].DomainName)
	d.Set("domain_id", d.Id())
	return nil
}
func resourceApsaraStackDnsDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}
	remarkUpdate := false
	check, err := dnsService.DescribeDnsDomain(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsDomainExist", ApsaraStackSdkGoERROR)
	}

	var desc string

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			desc = v.(string)
		}
		check.ZoneList[0].Remark = desc
		remarkUpdate = true
	} else {
		if v, ok := d.GetOk("remark"); ok {
			desc = v.(string)
		}
		check.ZoneList[0].Remark = desc
	}
	if remarkUpdate {
		request := requests.NewCommonRequest()
		request.Method = "POST"
		request.Product = "GenesisDns"
		request.Domain = client.Domain
		request.Version = "2018-07-20"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "RemarkGlobalAuthZone"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "GenesisDns",
			"RegionId":        client.RegionId,
			"Action":          "RemarkGlobalAuthZone",
			"Version":         "2018-07-20",
			"Id":              d.Id(),
			"Remark":          desc,
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_domain", "RemarkGlobalAuthZone", raw)
		}
		addDebug(request.GetActionName(), raw, request)
	}
	return resourceApsaraStackDnsDomainRead(d, meta)
}
func resourceApsaraStackDnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := requests.NewCommonRequest()
	request.Method = "POST"        // Set request method
	request.Product = "GenesisDns" // Specify product
	request.Domain = client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2018-07-20" // Specify product version
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DeleteGlobalZone"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "GenesisDns",
		"RegionId":        client.RegionId,
		"Action":          "DeleteGlobalZone",
		"Version":         "2018-07-20",
		"Id":              d.Id(),
	}
	raw, err := client.WithEcsClient(func(dnsClient *ecs.Client) (interface{}, error) {
		return dnsClient.ProcessCommonRequest(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return nil
}

package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"tags": tagsSchema(),
		},
	}
}

func resourceApsaraStackDnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := alidns.CreateAddDomainRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = client.Department
	request.QueryParams["ResourceGroup"] = client.ResourceGroup
	request.DomainName = d.Get("domain_name").(string)
	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = v.(string)
	}
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	raw, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.AddDomain(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_domain", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*alidns.AddDomainResponse)
	d.SetId(response.DomainName)

	return resourceApsaraStackDnsDomainUpdate(d, meta)
}
func resourceApsaraStackDnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}
	object, err := dnsService.DescribeDnsDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			//return nil
		}
		return WrapError(err)
	}

	d.Set("domain_name", d.Id())
	d.Set("dns_servers", object.DnsServers.DnsServer)
	d.Set("domain_id", object.DomainId)
	d.Set("group_id", object.GroupId)
	d.Set("remark", object.Remark)

	listTagResourcesObject, err := dnsService.ListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tags := make(map[string]string)
	for _, t := range listTagResourcesObject.TagResources {
		tags[t.TagKey] = t.TagValue
	}
	d.Set("tags", tags)
	return nil
}
func resourceApsaraStackDnsDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}
	d.Partial(true)

	update := false
	request := alidns.CreateChangeDomainGroupRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = client.Department
	request.QueryParams["ResourceGroup"] = client.ResourceGroup
	request.DomainName = d.Id()
	if !d.IsNewResource() && d.HasChange("group_id") {
		update = true
		request.GroupId = d.Get("group_id").(string)
	}
	if !d.IsNewResource() && d.HasChange("lang") {
		update = true
		request.Lang = d.Get("lang").(string)
	}
	if update {
		raw, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.ChangeDomainGroup(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		d.SetPartial("group_id")
		d.SetPartial("lang")
	}
	update = false
	updateDomainRemarkReq := alidns.CreateUpdateDomainRemarkRequest()
	updateDomainRemarkReq.Headers = map[string]string{"RegionId": client.RegionId}
	updateDomainRemarkReq.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "alidns"}
	updateDomainRemarkReq.QueryParams["Department"] = client.Department
	updateDomainRemarkReq.QueryParams["ResourceGroup"] = client.ResourceGroup
	updateDomainRemarkReq.DomainName = d.Id()
	updateDomainRemarkReq.Lang = d.Get("lang").(string)
	if d.HasChange("remark") {
		update = true
		updateDomainRemarkReq.Remark = d.Get("remark").(string)
	}
	if update {
		raw, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UpdateDomainRemark(updateDomainRemarkReq)
		})
		addDebug(updateDomainRemarkReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), updateDomainRemarkReq.GetActionName(), ApsaraStackSdkGoERROR)
		}
		d.SetPartial("remark")
	}
	if d.HasChange("tags") {
		if err := dnsService.SetResourceTags(d, "DOMAIN"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceApsaraStackDnsDomainRead(d, meta)
}
func resourceApsaraStackDnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := alidns.CreateDeleteDomainRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = client.Department
	request.QueryParams["ResourceGroup"] = client.ResourceGroup
	request.DomainName = d.Id()
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	raw, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DeleteDomain(request)
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

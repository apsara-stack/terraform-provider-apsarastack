package apsarastack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceApsaraStackDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDnsRecordCreate,
		Read:   resourceApsaraStackDnsRecordRead,
		Update: resourceApsaraStackDnsRecordUpdate,
		Delete: resourceApsaraStackDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"record_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "CNAME", "MX", "TXT", "PTR", "SRV", "NAPRT", "CAA", "NS"}, false),
			},
			"lba_strategy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL_RR", "RATIO"}, false),
			},
			"rr_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"line_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	//var requestInfo *ecs.Client
	// request := make(map[string]interface{})
	ZoneId := d.Get("zone_id").(string)
	LbaStrategy := d.Get("lba_strategy").(string)
	Type := d.Get("type").(string)
	Name := d.Get("name").(string)
	TTL := d.Get("ttl").(int)
	// var response map[string]interface{}

	// action := "AddGlobalZoneRecord"
	// request["Product"] = "CloudDns"
	// request["product"] = "CloudDns"
	// request["OrganizationId"] = client.Department
	// request["RegionId"] = client.RegionId

	// request["Type"] = Type
	// request["Ttl"] = TTL
	// request["ZoneId"] = ZoneId
	// request["LbaStrategy"] = LbaStrategy
	// request["Name"] = Name
	// conn, err := client.NewDataworkspublicClient()
	// if err != nil {
	// 	return WrapError(err)
	// }
	// request["ClientToken"] = buildClientToken("AddGlobalZoneRecord")
	// runtime := util.RuntimeOptions{}
	// runtime.SetAutoretry(true)
	// wait := incrementalWait(3*time.Second, 3*time.Second)
	// err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
	// 	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-06-24"), StringPointer("AK"), nil, request, &runtime)
	// 	if err != nil {
	// 		if NeedRetry(err) {
	// 			wait()
	// 			return resource.RetryableError(err)
	// 		}
	// 		return resource.NonRetryableError(err)
	// 	}
	// 	return nil
	// })
	request := requests.NewCommonRequest()
	request.Method = "POST"        // Set request method
	request.Product = "CloudDns"   // Specify product
	request.Domain = client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2021-06-24" // Specify product version
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	line_ids := expandStringList(d.Get("line_ids").(*schema.Set).List())
	if len(line_ids) <= 0 {
		line_ids = []string{"default"}
	}

	line_ids_json, _ := json.Marshal(line_ids)
	line_ids_str := string(line_ids_json)
	request.ApiName = "AddGlobalZoneRecord"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "CloudDns",
		"RegionId":        client.RegionId,
		"Action":          "AddGlobalZoneRecord",
		"Version":         "2021-06-24",
		"Name":            Name,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Type":            Type,
		"Ttl":             fmt.Sprintf("%d", TTL),
		"ZoneId":          ZoneId,
		"LbaStrategy":     LbaStrategy,
		"ClientToken":     buildClientToken("AddGlobalZoneRecord"),
		"LineIds":         line_ids_str,
	}
	var rrsets []string
	if v, ok := d.GetOk("rr_set"); ok {
		rrsets = expandStringList(v.(*schema.Set).List())
		for i, key := range rrsets {
			request.QueryParams[fmt.Sprintf("RDatas.%d.Value", i+1)] = key

		}
	}
	raw, err := client.WithEcsClient(func(dnsClient *ecs.Client) (interface{}, error) {
		return dnsClient.ProcessCommonRequest(request)
	})
	addDebug(request.GetActionName(), raw, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ApsaraStack_dns_record", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "ApsaraStack_dns_record", "AddGlobalZoneRecord", ApsaraStackSdkGoERROR)
	}
	addDebug("AddGlobalZoneRecord", raw, request, bresponse.GetHttpContentString())
	var resp map[string]interface{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ApsaraStack_dns_record", "AddGlobalZoneRecord", ApsaraStackSdkGoERROR)
	}
	if resp["asapiSuccess"].(bool) == false {
		return WrapErrorf(err, DefaultErrorMsg, "ApsaraStack_dns_record", "AddGlobalZoneRecord", ApsaraStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(ZoneId))

	return resourceApsaraStackDnsRecordRead(d, meta)
}

func resourceApsaraStackDnsRecordRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)

	dnsService := &DnsService{client: client}
	object, err := dnsService.DescribeDnsRecord(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("ttl", object.Data[0].TTL)
	d.Set("record_id", object.Data[0].Id)
	d.Set("name", object.Data[0].Name)
	d.Set("type", object.Data[0].Type)
	d.Set("remark", object.Data[0].Remark)
	d.Set("zone_id", object.Data[0].ZoneId)
	d.Set("lba_strategy", object.Data[0].LbaStrategy)

	return nil
}
func resourceApsaraStackDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}
	ID := d.Get("record_id").(int)
	ZoneId := d.Get("zone_id").(string)
	Name := d.Get("name").(string)
	LbaStrategy := d.Get("lba_strategy").(string)
	check, err := dnsService.DescribeDnsRecord(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsRecordExist", ApsaraStackSdkGoERROR)
	}
	attributeUpdate := false

	var desc string

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			desc = v.(string)
		}
		check.Data[0].Remark = desc
		request := requests.NewCommonRequest()
		request.Method = "POST"
		request.Product = "CloudDns"
		request.Domain = client.Domain
		request.Version = "2021-06-24"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "UpdateGlobalZoneRecordRemark"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "CloudDns",
			"RegionId":        client.RegionId,
			"Action":          "UpdateGlobalZoneRecordRemark",
			"Version":         "2021-06-24",
			"Id":              fmt.Sprint(ID),
			"Remark":          desc,
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw UpdateGlobalZoneRecordRemark : %s", raw)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "ApsaraStack_dns_record", "UpdateGlobalZoneRecordRemark", raw)
		}
		addDebug(request.GetActionName(), raw, request)
	} else {
		if v, ok := d.GetOk("remark"); ok {
			desc = v.(string)
		}
		check.Data[0].Remark = desc
	}

	var Type string
	var Ttl int

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			Type = v.(string)
		}
		check.Data[0].Type = Type
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("type"); ok {
			Type = v.(string)
		}
		check.Data[0].Type = Type
	}
	if d.HasChange("ttl") {
		if v, ok := d.GetOk("ttl"); ok {
			Ttl = v.(int)
		}
		check.Data[0].TTL = Ttl
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("ttl"); ok {
			Ttl = v.(int)
		}
		check.Data[0].TTL = Ttl
	}

	//var rrset string

	//if v, ok := d.GetOk("rr_set"); ok {
	//	rrsets = expandStringList(v.(*schema.Set).List())
	//	for i, key := range rrsets {
	//		request[fmt.Sprintf("RDatas.%d.Value", i+1)] = key
	//
	//	}
	//}

	if d.HasChange("rr_set") {
		attributeUpdate = true
	}

	if attributeUpdate {
		request := make(map[string]interface{})
		var rrsets []string
		if v, ok := d.GetOk("rr_set"); ok {
			rrsets = expandStringList(v.(*schema.Set).List())
			for i, key := range rrsets {
				request[fmt.Sprintf("RDatas.%d.Value", i+1)] = key

			}
		}
		action := "UpdateGlobalZoneRecord"
		request["Product"] = "CloudDns"
		request["product"] = "CloudDns"
		request["OrganizationId"] = client.Department
		request["RegionId"] = client.RegionId
		request["Type"] = Type
		request["Ttl"] = Ttl
		request["Id"] = ID
		request["ZoneId"] = ZoneId
		request["LbaStrategy"] = LbaStrategy
		request["Name"] = Name
		request["Remark"] = check.Data[0].Remark
		conn, err := client.NewDataworkspublicClient()
		if err != nil {
			return WrapError(err)
		}
		var response map[string]interface{}
		request["ClientToken"] = buildClientToken("UpdateGlobalZoneRecord")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-06-24"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		addDebug("UpdateGlobalZoneRecord", response, request)

	}

	return resourceApsaraStackDnsRecordRead(d, meta)
}

func resourceApsaraStackDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ID := d.Get("record_id").(int)
	ZoneId := d.Get("zone_id").(string)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "CloudDns"
	request.Domain = client.Domain
	request.Version = "2021-06-24"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DeleteGlobalZoneRecord"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "CloudDns",
		"RegionId":        client.RegionId,
		"Action":          "DeleteGlobalZoneRecord",
		"Version":         "2021-06-24",
		"Id":              fmt.Sprint(ID),
		"ZoneId":          fmt.Sprint(ZoneId),
	}
	raw, err := client.WithEcsClient(func(dnsClient *ecs.Client) (interface{}, error) {
		return dnsClient.ProcessCommonRequest(request)
	})
	addDebug(request.GetActionName(), raw)

	if err != nil {
		if IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser"}) {
			return nil
		}
		if IsExpectedErrors(err, []string{"RecordForbidden.DNSChange", "InternalError"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
	}
	return nil
}

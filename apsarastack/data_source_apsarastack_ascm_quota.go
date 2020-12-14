package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceApsaraStackQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackQuotaRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"product_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quota_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"quota_type_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"quota_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_type_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used_vip_public": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allocate_vip_internal": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allocate_vip_public": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_vip_public": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_vip_internal": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceApsaraStackQuotaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.Scheme = "http"
	request.RegionId = client.RegionId
	request.ApiName = "GetQuota"
	productName := d.Get("product_name").(string)
	quotaType := d.Get("quota_type").(string)
	quotaTypeId := d.Get("quota_type_id").(string)
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey,
		"Product":       "ascm",
		"RegionId":      client.RegionId,
		"Department":    client.Department,
		"ResourceGroup": client.ResourceGroup,
		"Action":        "GetQuota",
		"Version":       "2019-05-10",
		"ProductName":   productName,
		"QuotaType":     quotaType,
		"QuotaTypeId":   quotaTypeId,
		"RegionName":    client.RegionId,
	}
	response := QuotaData{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_quota", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.Code == "200" {
			break
		}

	}

	var ids []string
	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":                    response.Data.ID,
		"quota_type":            response.Data.QuotaType,
		"quota_type_id":         string(response.Data.QuotaTypeID),
		"used_vip_public":       response.Data.UsedVipPublic,
		"allocate_vip_internal": response.Data.AllocateVipInternal,
		"allocate_vip_public":   response.Data.AllocateVipPublic,
		"total_vip_public":      response.Data.TotalVipPublic,
		"total_vip_internal":    response.Data.TotalVipInternal,
		"region":                response.Data.Region,
	}

	ids = append(ids, string(rune(response.Data.ID)))
	s = append(s, mapping)

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

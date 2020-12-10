package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceApsaraStackRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackRegionsRead,
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
			"organization": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackRegionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.Scheme = "http"
	request.RegionId = client.RegionId
	request.ApiName = "GetRegionsByProduct"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey, "Product": "ascm", "RegionId": client.RegionId, "Department": client.Department, "ResourceGroup": client.ResourceGroup, "Action": "GetRegionsByProduct", "Version": "2019-05-10"}
	response := RegionsByProduct{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_regions", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.Code == 200 || len(response.Body.RegionList) < 1 {
			break
		}

	}

	var ids []string
	var s []map[string]interface{}
	for _, rg := range response.Body.RegionList {
		mapping := map[string]interface{}{
			"region_id":   rg.RegionID,
			"region_type": rg.RegionType,
		}
		ids = append(ids, rg.RegionID)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("region_list", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

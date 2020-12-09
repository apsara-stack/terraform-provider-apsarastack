package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"regexp"
)

func dataSourceApsaraStackInstanceFamilies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackInstanceFamiliesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"families": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackInstanceFamiliesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.Scheme = "http"
	request.RegionId = client.RegionId
	request.ApiName = "DescribeSeriesIdFamilies"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey, "Product": "ascm", "RegionId": client.RegionId, "Action": "DescribeSeriesIdFamilies", "Version": "2019-05-10"}
	response := InstanceFamily{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_instance_families", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.Success == true || len(response.Data) < 1 {
			break
		}

	}

	var r *regexp.Regexp
	if rt, ok := d.GetOk("resource_type"); ok && rt.(string) != "" {
		r = regexp.MustCompile(rt.(string))
	}
	var ids []string
	var s []map[string]interface{}
	for _, rg := range response.Data {
		if r != nil && !r.MatchString(rg.ResourceType) {
			continue
		}
		mapping := map[string]interface{}{
			"id":            rg.OrderBy.ID,
			"resource_type": rg.ResourceType,
			"status":        rg.Status,
		}
		ids = append(ids, rg.OrderBy.ID)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("families", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

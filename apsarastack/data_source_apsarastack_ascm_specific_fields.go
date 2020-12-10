package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceApsaraStackSpecificFields() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackSpecificFieldsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"group_filed": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"OSS", "ADB", "DRDS", "SLB",
					"NAT", "MAXCOMPUTE", "POSTGRESQL",
					"ECS", "RDS", "IPSIX", "REDIS", "MONGODB", "HITSDB"}, false),
			},
			"label": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"specific_fields": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceApsaraStackSpecificFieldsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.Scheme = "http"
	request.RegionId = client.RegionId
	request.ApiName = "GroupCommonSpec"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	resourceType := d.Get("resource_type").(string)
	groupFiled := d.Get("group_filed").(string)
	request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey, "Product": "ascm", "RegionId": client.RegionId, "Department": client.Department, "ResourceGroup": client.ResourceGroup, "Action": "GroupCommonSpec", "Version": "2019-05-10", "resourceType": resourceType, "groupFiled": groupFiled}
	response := SpecificField{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_specific_fields", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.Code == 200 || len(response.Data) < 1 {
			break
		}

	}
	var ids []string
	var s []map[string]interface{}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("specific_fields", response.Data); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

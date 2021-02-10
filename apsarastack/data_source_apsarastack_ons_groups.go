package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"log"
	"regexp"

	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceApsaraStackOnsGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackOnsGroupsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
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
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"independent_naming": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackOnsGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	namespaceid := d.Get("instance_id").(string)

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "Ons-inner"
	request.Version = "2018-02-05"
	request.Scheme = "http"
	request.RegionId = client.RegionId
	request.ApiName = "ConsoleGroupList"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey,
		"AccessKeySecret": client.SecretKey,
		"Product":         "Ons-inner",
		"RegionId":        client.RegionId,
		"Action":          "ConsoleGroupList",
		"Version":         "2018-02-05",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"OnsRegionId":     client.RegionId,
		"PreventCache":    "",
		"InstanceId":      namespaceid,
	}
	response := OnsGroup{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ConsoleGroupList : %s", raw)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_ascm_ons_groups", request.GetActionName(), ApsaraStackSdkGoERROR)
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
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("group_id_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	var ids []string
	var s []map[string]interface{}
	for _, ons := range response.Data {
		if r != nil && !r.MatchString(ons.GroupID) {
			continue
		}
		mapping := map[string]interface{}{
			"id":                 ons.ID,
			"remark":             ons.Remark,
			"instance_id":        ons.NamespaceID,
			"group_id":           ons.GroupID,
			"owner":              ons.Owner,
			"independent_naming": ons.IndependentNaming,
			"create_time":        ons.CreateTime,
		}
		ids = append(ids, string(rune(ons.ID)))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		_ = writeToFile(output.(string), s)
	}
	return nil

}

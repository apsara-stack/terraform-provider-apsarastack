package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity/ascm"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"regexp"
)

func dataSourceApsaraStackAscmLogonPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackAscmLogonPoliciesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			//"time_ranges": {
			//	Type:     schema.TypeList,
			//	Computed: true,
			//	Elem:     &schema.Schema{Type: schema.TypeString},
			//},
			"rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"login_policy_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_range": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"login_policy_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackAscmLogonPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := requests.NewCommonRequest()
	request.Method = "Get"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.Scheme = "http"
	request.RegionId = client.RegionId
	request.ApiName = "ListLoginPolicies"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey, "Product": "ascm", "RegionId": client.RegionId, "Action": "ListLoginPolicies", "Version": "2019-05-10"}
	response := ascm.LogonPolicy{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_ascm_logon_policies", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.Code == "200" || len(response.Data) < 1 {
			break
		}

	}

	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	var ids []string
	var s []map[string]interface{}
	var t []map[string]interface{}
	for _, u := range response.Data {
		if r != nil && !r.MatchString(u.Name) {
			continue
		}
		for _, time := range response.Data {
			var ipranges []string
			var iprange string
			for _, k := range u.IPRanges {
				ipranges = append(ipranges, k.IPRange)
				if len(ipranges) > 1 {
					iprange = iprange + "," + k.IPRange
				} else {
					iprange = k.IPRange
				}
			}
			timemapping := map[string]interface{}{
				"id":          u.LpID,
				"name":        u.Name,
				"rule":        u.Rule,
				"description": u.Description,
				"ip_range":    iprange,
			}
			s = append(s, timemapping)
			for _, k := range time.TimeRanges {
				allmapping := map[string]interface{}{
					"start_time":      k.StartTime,
					"end_time":        k.EndTime,
					"login_policy_id": k.LoginPolicyID,
				}
				t = append(t, allmapping)
			}
			for k, v := range t {
				s[k] = v
			}

		}

	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}
	//if err := d.Set("time_ranges", t); err != nil {
	//	return WrapError(err)
	//}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

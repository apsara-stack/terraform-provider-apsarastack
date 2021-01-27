package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"regexp"
	"strings"
)

func dataSourceApsarastackCmsAlarms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsarastackCmsAlarmsRead,
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"alarms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"metric_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"no_effective_interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"silence_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"contact_groups": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mail_subject": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"rule_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"period": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"dimensions": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"effective_interval": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enable_state": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"webhook": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"escalations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"critical_comparison_operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"critical_times": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"critical_statistics": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"critical_threshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"info_comparison_operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"info_times": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"info_statistics": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"info_threshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"warn_comparison_operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"warn_times": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"warn_statistics": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"warn_threshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
func dataSourceApsarastackCmsAlarmsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.Method = "GET"
	request.Product = "cms"
	request.Version = "2019-01-01"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.ApiName = "DescribeMetricRuleList"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeyId":     client.AccessKey,
		"AccessKeySecret": client.SecretKey,
		"Product":         "Cms",
		"RegionId":        client.RegionId,
		"Action":          "DescribeMetricRuleList",
		"Version":         "2019-01-01",
	}
	response := AlarmsData{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_cms_alarms", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.Success == true || len(response.Alarms.Alarm) < 1 {
			break
		}

	}
	log.Printf("bhushan123 %v", response)
	var r *regexp.Regexp
	if rt, ok := d.GetOk("name_regex"); ok && rt.(string) != "" {
		r = regexp.MustCompile(rt.(string))
	}
	var ids []string
	var s []map[string]interface{}

	for _, data := range response.Alarms.Alarm {
		if r != nil && !r.MatchString(data.RuleName) {
			continue
		}
		mapping := map[string]interface{}{
			"critical_comparison_operator": data.Escalations.Critical.ComparisonOperator,
			"critical_times":               data.Escalations.Critical.Times,
			"critical_statistics":          data.Escalations.Critical.Statistics,
			"critical_threshold":           data.Escalations.Critical.Threshold,
			"info_comparison_operator":     data.Escalations.Info.ComparisonOperator,
			"info_times":                   data.Escalations.Info.Times,
			"info_statistics":              data.Escalations.Info.Statistics,
			"info_threshold":               data.Escalations.Info.Threshold,
			"warn_comparison_operator":     data.Escalations.Warn.ComparisonOperator,
			"warn_times":                   data.Escalations.Warn.Times,
			"warn_statistics":              data.Escalations.Warn.Statistics,
			"warn_threshold":               data.Escalations.Warn.Threshold,
		}
		mapping1 := map[string]interface{}{
			"resources":          data.Resources,
			"namespace":          data.Namespace,
			"rule_id":            data.RuleID,
			"dimensions":         data.Dimensions,
			"metric_name":        data.MetricName,
			"period":             data.Period,
			"contact_groups":     data.ContactGroups,
			"effective_interval": data.EffectiveInterval,
			"enable_state":       data.EnableState,
			"escalations":        []map[string]interface{}{mapping},
			"rule_name":          data.RuleName,
		}

		ids = append(ids, data.RuleID)
		s = append(s, mapping1)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alarms", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

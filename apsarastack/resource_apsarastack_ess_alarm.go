package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackEssAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackEssAlarmCreate,
		Read:   resourceApsaraStackEssAlarmRead,
		Update: resourceApsaraStackEssAlarmUpdate,
		Delete: resourceApsaraStackEssAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"alarm_actions": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MaxItems: 5,
				MinItems: 1,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metric_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "system",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"system", "custom"}, false),
			},
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 900}),
			},
			"statistics": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  Average,
				ValidateFunc: validation.StringInSlice([]string{
					string(Average),
					string(Minimum),
					string(Maximum),
				}, false),
			},
			"threshold": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comparison_operator": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ">=",
				ValidateFunc: validation.StringInSlice([]string{">", ">=", "<", "<="}, false),
			},
			"evaluation_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"cloud_monitor_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dimensions": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApsaraStackEssAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request, err := buildApsaraStackEssAlarmArgs(d)
	if err != nil {
		return WrapError(err)
	}
	log.Printf("checking built request %v", request)

	request.RegionId = client.RegionId
	request.Domain = client.Domain
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "Ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	log.Printf("checking again built request %v", request)

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.CreateAlarm(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ess_alarm", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request)
	response, _ := raw.(*ess.CreateAlarmResponse)
	log.Printf("checking response %v", response)
	if response == nil {
		return WrapErrorf(err, "Null response found", "apsarastack_ess_alarm", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	d.SetId(response.AlarmTaskId)
	// enable or disable alarm
	enable := d.Get("enable")
	if !enable.(bool) {

		err := enableordisableAlarm(false, d.Id(), meta)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
	}
	return resourceApsaraStackEssAlarmRead(d, meta)
}

func resourceApsaraStackEssAlarmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	essService := EssService{client}

	object, err := essService.DescribeEssAlarm(d.Id())
	if err != nil {

		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("alarm_actions", object.AlarmActions.AlarmAction)
	d.Set("scaling_group_id", object.ScalingGroupId)
	d.Set("metric_type", object.MetricType)
	d.Set("metric_name", object.MetricName)
	d.Set("period", object.Period)
	d.Set("statistics", object.Statistics)
	d.Set("threshold", strconv.FormatFloat(object.Threshold, 'f', -1, 32))
	d.Set("comparison_operator", object.ComparisonOperator)
	d.Set("evaluation_count", object.EvaluationCount)
	d.Set("state", object.State)
	d.Set("enable", object.Enable)

	dims := make([]ess.Dimension, 0, len(object.Dimensions.Dimension))
	for _, dimension := range object.Dimensions.Dimension {
		if dimension.DimensionKey == GroupId {
			d.Set("cloud_monitor_group_id", dimension.DimensionValue)
		} else {
			dims = append(dims, dimension)
		}
	}

	if err := d.Set("dimensions", essService.flattenDimensionsToMap(dims)); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceApsaraStackEssAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := ess.CreateModifyAlarmRequest()
	request.AlarmTaskId = d.Id()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.Domain = client.Domain
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "Ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	d.Partial(true)
	if metricType, ok := d.GetOk("metric_type"); ok && metricType.(string) != "" {
		request.MetricType = metricType.(string)
	}
	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}

	if d.HasChange("alarm_actions") {
		if v, ok := d.GetOk("alarm_actions"); ok {
			alarmActions := expandStringList(v.(*schema.Set).List())
			if len(alarmActions) > 0 {
				request.AlarmAction = &alarmActions
			}
		}
	}
	if d.HasChange("metric_name") {
		request.MetricName = d.Get("metric_name").(string)
	}
	if d.HasChange("statistics") {
		request.Statistics = d.Get("statistics").(string)
	}
	if d.HasChange("threshold") {
		request.Threshold = requests.Float(d.Get("threshold").(string))
	}
	if d.HasChange("comparison_operator") {
		request.ComparisonOperator = d.Get("comparison_operator").(string)
	}
	if d.HasChange("evaluation_count") {
		request.EvaluationCount = requests.NewInteger(d.Get("evaluation_count").(int))
	}
	if v, ok := d.GetOk("cloud_monitor_group_id"); ok {
		request.GroupId = requests.NewInteger(v.(int))
	}

	dimensions := d.Get("dimensions").(map[string]interface{})
	createAlarmDimensions := make([]ess.ModifyAlarmDimension, 0, len(dimensions))
	for k, v := range dimensions {
		if k == UserId || k == ScalingGroup {
			return WrapError(Error("Invalide dimension keys, %s", k))
		}
		if k != "" {
			dimension := ess.ModifyAlarmDimension{
				DimensionKey:   k,
				DimensionValue: v.(string),
			}
			createAlarmDimensions = append(createAlarmDimensions, dimension)
		}
	}
	request.Dimension = &createAlarmDimensions

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyAlarm(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	d.SetPartial("name")
	d.SetPartial("description")
	d.SetPartial("alarm_actions")
	d.SetPartial("metric_name")
	d.SetPartial("statistics")
	d.SetPartial("threshold")
	d.SetPartial("comparison_operator")
	d.SetPartial("evaluation_count")
	d.SetPartial("cloud_monitor_group_id")
	d.SetPartial("dimensions")
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if d.HasChange("enable") {
		enable := d.Get("enable").(bool)
		err := enableordisableAlarm(enable, d.Id(), meta)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		d.SetPartial("enable")
	}
	d.Partial(false)
	return resourceApsaraStackEssAlarmRead(d, meta)
}

func resourceApsaraStackEssAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	essService := EssService{client}
	request := ess.CreateDeleteAlarmRequest()
	request.AlarmTaskId = d.Id()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.Domain = client.Domain
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "Ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteAlarm(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"404"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(essService.WaitForEssAlarm(d.Id(), Deleted, DefaultTimeout))
}

func buildApsaraStackEssAlarmArgs(d *schema.ResourceData) (*ess.CreateAlarmRequest, error) {
	request := ess.CreateCreateAlarmRequest()

	if name, ok := d.GetOk("name"); ok && name.(string) != "" {
		request.Name = name.(string)
	}

	if description, ok := d.GetOk("description"); ok && description.(string) != "" {
		request.Description = description.(string)
	}

	if v, ok := d.GetOk("alarm_actions"); ok {
		alarmActions := expandStringList(v.(*schema.Set).List())
		request.AlarmAction = &alarmActions
	}

	if scalingGroupId := d.Get("scaling_group_id").(string); scalingGroupId != "" {
		request.ScalingGroupId = scalingGroupId
	}

	if metricType, ok := d.GetOk("metric_type"); ok && metricType.(string) != "" {
		request.MetricType = metricType.(string)
	}

	if metricName := d.Get("metric_name").(string); metricName != "" {
		request.MetricName = metricName
	}

	if period, ok := d.GetOk("period"); ok && period.(int) > 0 {
		request.Period = requests.NewInteger(period.(int))
	}

	if statistics, ok := d.GetOk("statistics"); ok && statistics.(string) != "" {
		request.Statistics = statistics.(string)
	}

	if v, ok := d.GetOk("threshold"); ok {
		threshold, err := strconv.ParseFloat(v.(string), 32)
		if err != nil {
			return nil, WrapError(err)
		}
		request.Threshold = requests.NewFloat(threshold)
	}

	if comparisonOperator, ok := d.GetOk("comparison_operator"); ok && comparisonOperator.(string) != "" {
		request.ComparisonOperator = comparisonOperator.(string)
	}

	if evaluationCount, ok := d.GetOk("evaluation_count"); ok && evaluationCount.(int) > 0 {
		request.EvaluationCount = requests.NewInteger(evaluationCount.(int))
	}

	if groupId, ok := d.GetOk("cloud_monitor_group_id"); ok {
		request.GroupId = requests.NewInteger(groupId.(int))
	}

	if v, ok := d.GetOk("dimensions"); ok {
		dimensions := v.(map[string]interface{})
		createAlarmDimensions := make([]ess.CreateAlarmDimension, 0, len(dimensions))
		for k, v := range dimensions {
			if k == UserId || k == ScalingGroup {
				return nil, WrapError(Error("Invalide dimension keys, %s", k))
			}
			if k != "" {
				dimension := ess.CreateAlarmDimension{
					DimensionKey:   k,
					DimensionValue: v.(string),
				}
				createAlarmDimensions = append(createAlarmDimensions, dimension)
			}
		}
		request.Dimension = &createAlarmDimensions
	}

	return request, nil
}

func enableordisableAlarm(check bool, id string, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var apiaction string
	if check {
		apiaction = "EnableAlarm"
	} else {
		apiaction = "DisableAlarm"
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "Ess"
	request.Version = "2014-08-28"
	request.Scheme = "http"
	request.ServiceCode = "ess"
	request.ApiName = apiaction
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.RegionId = client.RegionId
	request.Domain = client.Domain
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"Product":         "Ess",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          apiaction,
		"Version":         "2014-08-28",
		"ProductName":     "ess",
		"RegionId":        client.RegionId,
		"AlarmTaskId":     id,
	}
	raw, err := client.WithEcsClient(func(ess *ecs.Client) (interface{}, error) {
		return ess.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	response := raw.(*responses.CommonResponse)
	if !response.IsSuccess() {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request)
	return nil
}

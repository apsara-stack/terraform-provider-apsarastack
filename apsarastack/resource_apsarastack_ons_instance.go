package apsarastack

import (
	"encoding/json"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"
	"time"
)

func resourceApsaraStackOnsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackOnsInstanceCreate,
		Read:   resourceApsaraStackOnsInstanceRead,
		Update: resourceApsaraStackOnsInstanceUpdate,
		Delete: resourceApsaraStackOnsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},

			"tps_receive_max": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tps_send_max": {
				Type:     schema.TypeString,
				Required: true,
			},
			"topic_capacity": {
				Type:     schema.TypeString,
				Required: true,
			},
			"independent_naming": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},

			// Computed Values
			"instance_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApsaraStackOnsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client
	maxrtps := d.Get("tps_receive_max").(string)
	maxstps := d.Get("tps_send_max").(string)
	topiccapacity := d.Get("topic_capacity").(string)
	independentname := d.Get("independent_naming").(string)
	ins_resp := OnsInstance{}

	cluster := d.Get("cluster").(string)
	name := d.Get("name").(string)
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":          client.RegionId,
		"AccessKeySecret":   client.SecretKey,
		"Department":        client.Department,
		"ResourceGroup":     client.ResourceGroup,
		"Product":           "Ons-inner",
		"Action":            "ConsoleInstanceCreate",
		"Version":           "2018-02-05",
		"ProductName":       "Ons-inner",
		"OnsRegionId":       client.RegionId,
		"InstanceName":      name,
		"MaxReceiveTps":     maxrtps,
		"MaxSendTps":        maxstps,
		"TopicCapacity":     topiccapacity,
		"Cluster":           cluster,
		"IndependentNaming": independentname,
	}
	request.Method = "POST"
	request.Product = "Ons-inner"
	request.Version = "2018-02-05"
	request.ServiceCode = "Ons-inner"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ConsoleInstanceCreate"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ons_instance", "ConsoleInstanceCreate", raw)
	}
	addDebug("ConsoleInstanceCreate", raw, requestInfo, request)
	bresponse, _ := raw.(*responses.CommonResponse)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &ins_resp)
	if ins_resp.Success != true {
		return WrapErrorf(errors.New(ins_resp.Message), DefaultErrorMsg, "apsarastack_ons_instance", "ConsoleInstanceCreate", ApsaraStackSdkGoERROR)
	}

	if err != nil {
		return WrapError(err)
	}
	d.SetId(ins_resp.Data.InstanceID)

	return resourceApsaraStackOnsInstanceUpdate(d, meta)
}

func resourceApsaraStackOnsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	onsService := OnsService{client}

	response, err := onsService.DescribeOnsInstance(d.Id())

	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", response.Data.InstanceName)
	d.Set("instance_type", response.Data.InstanceType)
	d.Set("instance_status", response.Data.InstanceStatus)
	d.Set("create_time", time.Unix(response.Data.CreateTime/1000, 0).Format("2006-01-02 03:04:05"))

	return nil
}

func resourceApsaraStackOnsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	onsService := OnsService{client}

	attributeUpdate := false
	check, err := onsService.DescribeOnsInstance(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsInstanceExist", ApsaraStackSdkGoERROR)
	}
	var name string

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		check.Data.InstanceName = name
		attributeUpdate = true
	}
	var remark string

	if d.HasChange("remark") {

		if v, ok := d.GetOk("remark"); ok {
			remark = v.(string)
		}
		check.Data.Remark = remark
		attributeUpdate = true
	}
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "Ons-inner",
		"Action":          "ConsoleInstanceUpdate",
		"Version":         "2018-02-05",
		"Remark":          remark,
		"InstanceName":    name,
		"OnsRegionId":     client.RegionId,
		"PreventCache":    "",
		"InstanceId":      d.Id(),
	}
	request.Method = "POST"
	request.Product = "Ons-inner"
	request.Version = "2018-02-05"
	request.ServiceCode = "Ons-inner"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.SetHTTPSInsecure(true)
	request.ApiName = "ConsoleInstanceUpdate"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	check.Data.InstanceID = d.Id()

	if attributeUpdate {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ons_instance", "ConsoleInstanceCreate", raw)
		}
		addDebug(request.GetActionName(), raw, request)
	}

	return resourceApsaraStackOnsInstanceRead(d, meta)
}

func resourceApsaraStackOnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	onsService := OnsService{client}
	var requestInfo *ecs.Client
	check, err := onsService.DescribeOnsInstance(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsInstanceExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsInstanceExist", check, requestInfo, map[string]string{"InstanceId": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ons-inner",
			"Action":          "ConsoleInstanceDelete",
			"Version":         "2018-02-05",
			"ProductName":     "Ons-inner",
			"PreventCache":    "",
			"OnsRegionId":     client.RegionId,
			"InstanceId":      d.Id(),
		}

		request.Method = "POST"
		request.Product = "Ons-inner"
		request.Version = "2018-02-05"
		request.ServiceCode = "Ons-inner"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "ConsoleInstanceDelete"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err = onsService.DescribeOnsInstance(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return nil
}

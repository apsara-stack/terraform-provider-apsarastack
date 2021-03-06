package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackSnapshotPolicyCreate,
		Read:   resourceApsaraStackSnapshotPolicyRead,
		Update: resourceApsaraStackSnapshotPolicyUpdate,
		Delete: resourceApsaraStackSnapshotPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"repeat_weekdays": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"time_points": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"disk_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_automated_snapshot_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceApsaraStackSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var err error
	autopolicy := d.Get("enable_automated_snapshot_policy").(bool)
	disks := convertListToJsonString(d.Get("disk_ids").(*schema.Set).List())
	if autopolicy == true && disks == "" {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_snapshot_policy", "AddDiskIds for EnableAutomatedSnapshotPolicy")
	}
	request := ecs.CreateCreateAutoSnapshotPolicyRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.AutoSnapshotPolicyName = d.Get("name").(string)
	request.RepeatWeekdays = convertListToJsonString(d.Get("repeat_weekdays").(*schema.Set).List())
	request.RetentionDays = requests.NewInteger(d.Get("retention_days").(int))
	request.TimePoints = convertListToJsonString(d.Get("time_points").(*schema.Set).List())

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateAutoSnapshotPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_snapshot_policy", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.CreateAutoSnapshotPolicyResponse)
	d.SetId(response.AutoSnapshotPolicyId)

	ecsService := EcsService{client}
	if err := ecsService.WaitForSnapshotPolicy(d.Id(), SnapshotPolicyNormal, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	if d.Get("enable_automated_snapshot_policy").(bool) {
		req := ecs.CreateApplyAutoSnapshotPolicyRequest()
		req.Headers = map[string]string{"RegionId": client.RegionId}
		req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		req.Domain = client.Domain
		req.DiskIds = convertListToJsonString(d.Get("disk_ids").(*schema.Set).List())
		req.AutoSnapshotPolicyId = d.Id()
		raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ApplyAutoSnapshotPolicy(req)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_disk", req.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, req.RpcRequest, req)
	}
	return resourceApsaraStackSnapshotPolicyRead(d, meta)
}

func resourceApsaraStackSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeSnapshotPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.AutoSnapshotPolicyName)
	weekdays, err := convertJsonStringToList(object.RepeatWeekdays)
	if err != nil {
		return WrapError(err)
	}
	d.Set("repeat_weekdays", weekdays)
	d.Set("retention_days", object.RetentionDays)
	d.Set("auto_snapshot_policy_id", object.AutoSnapshotPolicyId)
	timePoints, err := convertJsonStringToList(object.TimePoints)
	if err != nil {
		return WrapError(err)
	}
	d.Set("time_points", timePoints)

	return nil
}

func resourceApsaraStackSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := ecs.CreateModifyAutoSnapshotPolicyExRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.AutoSnapshotPolicyId = d.Id()
	if d.HasChange("name") {
		request.AutoSnapshotPolicyName = d.Get("name").(string)
	}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if d.HasChange("repeat_weekdays") {
		request.RepeatWeekdays = convertListToJsonString(d.Get("repeat_weekdays").(*schema.Set).List())
	}
	if d.HasChange("retention_days") {
		request.RetentionDays = requests.NewInteger(d.Get("retention_days").(int))
	}
	if d.HasChange("time_points") {
		request.TimePoints = convertListToJsonString(d.Get("time_points").(*schema.Set).List())
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyAutoSnapshotPolicyEx(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceApsaraStackSnapshotPolicyRead(d, meta)
}

func resourceApsaraStackSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	//ecsService := EcsService{client}
	log.Printf("autosnapshotpolicy25 %v", d.Id())
	request := ecs.CreateCancelAutoSnapshotPolicyRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.DiskIds = convertListToJsonString(d.Get("disk_ids").(*schema.Set).List())
	err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.CancelAutoSnapshotPolicy(request)
		})
		if err != nil {
			if IsExpectedErrors(err, SnapshotPolicyInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		resp := raw.(*ecs.CancelAutoSnapshotPolicyResponse)
		if resp.GetHttpStatus() != 200 {

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	return nil
	//return WrapError(ecsService.WaitForSnapshotPolicy(d.Id(), Deleted, DefaultTimeout))
}

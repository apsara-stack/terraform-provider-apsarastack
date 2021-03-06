package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackReservedInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackReservedInstanceCreate,
		Read:   resourceApsaraStackReservedInstanceRead,
		Update: resourceApsaraStackReservedInstanceUpdate,
		Delete: resourceApsaraStackReservedInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Region",
				ValidateFunc: validation.StringInSlice([]string{"Region", "Zone"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_amount": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Windows", "Linux"}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Year",
				ValidateFunc: validation.StringInSlice([]string{"Year"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1, 3}),
			},
			"offering_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"No Upfront", "Partial Upfront", "All Upfront"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}
func resourceApsaraStackReservedInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	request := ecs.CreatePurchaseReservedInstancesOfferingRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = v.(string)
	}
	request.RegionId = client.RegionId
	if scope, ok := d.GetOk("scope"); ok {
		request.Scope = scope.(string)
		if v, ok := d.GetOk("zone_id"); ok {
			if scope == "Zone" && v == "" {
				return WrapError(Error("Required when Scope is Zone."))
			}
			request.ZoneId = v.(string)
		}
	}
	if v, ok := d.GetOk("instance_amount"); ok {
		request.InstanceAmount = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("platform"); ok {
		request.Platform = v.(string)
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request.PeriodUnit = v.(string)
	}
	if v, ok := d.GetOk("period"); ok {
		request.Period = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("offering_type"); ok {
		request.OfferingType = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		request.ReservedInstanceName = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.PurchaseReservedInstancesOffering(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_reserved_instance", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.PurchaseReservedInstancesOfferingResponse)
	if len(response.ReservedInstanceIdSets.ReservedInstanceId) != 1 {
		return WrapError(Error("API returned wrong number of collections"))
	}
	d.SetId(response.ReservedInstanceIdSets.ReservedInstanceId[0])

	if err := ecsService.WaitForReservedInstance(d.Id(), Active, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceApsaraStackReservedInstanceRead(d, meta)
}
func resourceApsaraStackReservedInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := ecs.CreateModifyReservedInstanceAttributeRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.ReservedInstanceId = d.Id()
	request.RegionId = client.RegionId
	if d.HasChange("name") || d.HasChange("description") {
		if v, ok := d.GetOk("name"); ok {
			request.ReservedInstanceName = v.(string)
		}
		if v, ok := d.GetOk("description"); ok {
			request.Description = v.(string)
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyReservedInstanceAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceApsaraStackReservedInstanceRead(d, meta)
}
func resourceApsaraStackReservedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	reservedInstances, err := ecsService.DescribeReservedInstance(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_type", reservedInstances.InstanceType)
	d.Set("scope", reservedInstances.Scope)
	d.Set("zone_id", reservedInstances.ZoneId)
	d.Set("instance_amount", reservedInstances.InstanceAmount)
	d.Set("platform", reservedInstances.Platform)
	d.Set("offering_type", reservedInstances.OfferingType)
	d.Set("name", reservedInstances.ReservedInstanceName)
	d.Set("description", reservedInstances.Description)
	d.Set("resource_group_id", reservedInstances.ReservedInstanceId)

	return WrapError(err)
}
func resourceApsaraStackReservedInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	// PurchaseReservedInstancesOffering can not be release.
	return nil
}

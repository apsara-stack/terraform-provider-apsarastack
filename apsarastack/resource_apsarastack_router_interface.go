package apsarastack

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackRouterInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackRouterInterfaceCreate,
		Read:   resourceApsaraStackRouterInterfaceRead,
		Update: resourceApsaraStackRouterInterfaceUpdate,
		Delete: resourceApsaraStackRouterInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"opposite_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"router_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(VRouter), string(VBR)}, false),
				ForceNew:         true,
				DiffSuppressFunc: routerInterfaceAcceptsideDiffSuppressFunc,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(InitiatingSide), string(AcceptingSide)}, false),
				ForceNew: true,
			},
			"specification": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice(GetAllRouterInterfaceSpec(), false),
				DiffSuppressFunc: routerInterfaceAcceptsideDiffSuppressFunc,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"health_check_source_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: routerInterfaceVBRTypeDiffSuppressFunc,
			},
			"health_check_target_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: routerInterfaceVBRTypeDiffSuppressFunc,
			},
			"access_point_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"opposite_access_point_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"opposite_router_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"opposite_router_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"opposite_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"opposite_interface_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackRouterInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}
	request, err := buildApsaraStackRouterInterfaceCreateArgs(d, meta)
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateRouterInterface(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_router_interface", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.CreateRouterInterfaceResponse)
	d.SetId(response.RouterInterfaceId)

	if err := vpcService.WaitForRouterInterface(d.Id(), client.RegionId, Idle, 300); err != nil {
		return WrapError(err)
	}

	return resourceApsaraStackRouterInterfaceUpdate(d, meta)
}

func resourceApsaraStackRouterInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	d.Partial(true)

	request, attributeUpdate, err := buildApsaraStackRouterInterfaceModifyAttrArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	request.RegionId = client.RegionId

	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if attributeUpdate {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyRouterInterfaceAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChange("specification") && !d.IsNewResource() {
		d.SetPartial("specification")
		request := vpc.CreateModifyRouterInterfaceSpecRequest()
		request.RegionId = string(client.Region)

		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.RouterInterfaceId = d.Id()
		request.Spec = d.Get("specification").(string)
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyRouterInterfaceSpec(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.Partial(false)
	return resourceApsaraStackRouterInterfaceRead(d, meta)
}

func resourceApsaraStackRouterInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeRouterInterface(d.Id(), client.RegionId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	d.Set("role", object.Role)
	d.Set("specification", object.Spec)
	d.Set("name", object.Name)
	d.Set("router_id", object.RouterId)
	d.Set("router_type", object.RouterType)
	d.Set("description", object.Description)
	d.Set("access_point_id", object.AccessPointId)
	d.Set("opposite_region", object.OppositeRegionId)
	d.Set("opposite_router_type", object.OppositeRouterType)
	d.Set("opposite_router_id", object.OppositeRouterId)
	d.Set("opposite_interface_id", object.OppositeInterfaceId)
	d.Set("opposite_interface_owner_id", object.OppositeInterfaceOwnerId)
	d.Set("health_check_source_ip", object.HealthCheckSourceIp)
	d.Set("health_check_target_ip", object.HealthCheckTargetIp)
	return nil

}

func resourceApsaraStackRouterInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	if object, err := vpcService.DescribeRouterInterface(d.Id(), client.RegionId); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	} else if object.Status == string(Active) {
		if err := vpcService.DeactivateRouterInterface(d.Id()); err != nil {
			return WrapError(err)
		}
	}

	request := vpc.CreateDeleteRouterInterfaceRequest()
	request.RegionId = string(client.Region)
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RouterInterfaceId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteRouterInterface(&args)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus", "DependencyViolation.RouterInterfaceReferedByRouteEntry"}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return nil
		}
		WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return WrapError(vpcService.WaitForRouterInterface(d.Id(), client.RegionId, Deleted, DefaultTimeoutMedium))
}

func buildApsaraStackRouterInterfaceCreateArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateRouterInterfaceRequest, error) {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	oppositeRegion := d.Get("opposite_region").(string)
	if err := ecsService.JudgeRegionValidation("opposite_region", oppositeRegion); err != nil {
		return nil, WrapError(err)
	}

	request := vpc.CreateCreateRouterInterfaceRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RouterType = d.Get("router_type").(string)
	request.RouterId = d.Get("router_id").(string)
	request.Role = d.Get("role").(string)
	request.Spec = d.Get("specification").(string)
	request.OppositeRegionId = oppositeRegion
	// Accepting side router interface spec only be Negative and router type only be VRouter.
	if request.Role == string(AcceptingSide) {
		request.Spec = string(Negative)
		request.RouterType = string(VRouter)
	} else {
		if request.Spec == "" {
			return request, WrapError(Error("'specification': required field is not set when role is %s.", InitiatingSide))
		}
	}

	// Get VBR access point
	if request.RouterType == string(VBR) {
		describeVirtualBorderRoutersRequest := vpc.CreateDescribeVirtualBorderRoutersRequest()
		values := []string{request.RouterId}
		filters := []vpc.DescribeVirtualBorderRoutersFilter{{
			Key:   "VbrId",
			Value: &values,
		}}
		describeVirtualBorderRoutersRequest.Filter = &filters
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVirtualBorderRouters(describeVirtualBorderRoutersRequest)
		})
		if err != nil {
			return request, WrapErrorf(err, DefaultErrorMsg, "apsarastack_router_interface", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVirtualBorderRoutersResponse)
		if response.TotalCount > 0 {
			request.AccessPointId = response.VirtualBorderRouterSet.VirtualBorderRouterType[0].AccessPointId
		}
	}
	request.OppositeInterfaceId = d.Get("opposite_interface_id").(string)
	request.OppositeRouterType = d.Get("opposite_router_type").(string)
	request.OppositeRouterId = d.Get("opposite_router_id").(string)
	request.OppositeInterfaceOwnerId = d.Get("opposite_interface_owner_id").(string)
	if request.OppositeInterfaceOwnerId == "" {
		owner := request.OppositeInterfaceOwnerId
		owner, err := client.AccountId()
		if err != nil {
			//return WrapError(err.Error()
		}
		request.OppositeInterfaceOwnerId = owner
	}
	request.ClientToken = buildClientToken(request.GetActionName())
	return request, nil
}

func buildApsaraStackRouterInterfaceModifyAttrArgs(d *schema.ResourceData, meta interface{}) (*vpc.ModifyRouterInterfaceAttributeRequest, bool, error) {
	client := meta.(*connectivity.ApsaraStackClient)
	sourceIp, sourceOk := d.GetOk("health_check_source_ip")
	targetIp, targetOk := d.GetOk("health_check_target_ip")
	if sourceOk && !targetOk || !sourceOk && targetOk {
		return nil, false, WrapError(Error("The 'health_check_source_ip' and 'health_check_target_ip' should be specified or not at one time."))
	}

	request := vpc.CreateModifyRouterInterfaceAttributeRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RouterInterfaceId = d.Id()

	attributeUpdate := false

	if d.HasChange("health_check_source_ip") {
		d.SetPartial("health_check_source_ip")
		request.HealthCheckSourceIp = sourceIp.(string)
		request.HealthCheckTargetIp = targetIp.(string)
		attributeUpdate = true
	}

	if d.HasChange("health_check_target_ip") {
		d.SetPartial("health_check_target_ip")
		request.HealthCheckTargetIp = targetIp.(string)
		request.HealthCheckSourceIp = sourceIp.(string)
		attributeUpdate = true
	}

	if d.HasChange("name") {
		d.SetPartial("name")
		request.Name = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		request.Description = d.Get("description").(string)
		attributeUpdate = true
	}

	return request, attributeUpdate, nil
}

package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackEipAssociationCreate,
		Read:   resourceApsaraStackEipAssociationRead,
		Delete: resourceApsaraStackEipAssociationDelete,

		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceApsaraStackEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateAssociateEipAddressRequest()
	request.RegionId = client.RegionId

	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.AllocationId = Trim(d.Get("allocation_id").(string))
	request.InstanceId = Trim(d.Get("instance_id").(string))
	request.InstanceType = EcsInstance
	request.ClientToken = buildClientToken(request.GetActionName())

	if strings.HasPrefix(request.InstanceId, "lb-") {
		request.InstanceType = SlbInstance
	}
	if strings.HasPrefix(request.InstanceId, "ngw-") {
		request.InstanceType = Nat
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = instanceType.(string)
	}
	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.AssociateEipAddress(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_eip_association", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	if err := vpcService.WaitForEip(request.AllocationId, InUse, 60); err != nil {
		return WrapError(err)
	}
	// There is at least 30 seconds delay for ecs instance
	if request.InstanceType == EcsInstance {
		time.Sleep(30 * time.Second)
	}

	d.SetId(request.AllocationId + ":" + request.InstanceId)

	return resourceApsaraStackEipAssociationRead(d, meta)
}

func resourceApsaraStackEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeEipAssociation(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_id", object.InstanceId)
	d.Set("allocation_id", object.AllocationId)
	d.Set("instance_type", object.InstanceType)
	d.Set("force", d.Get("force").(bool))
	return nil
}

func resourceApsaraStackEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	allocationId, instanceId := parts[0], parts[1]
	if err != nil {
		return WrapError(err)
	}

	request := vpc.CreateUnassociateEipAddressRequest()
	request.RegionId = client.RegionId

	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.AllocationId = allocationId
	request.InstanceId = instanceId
	request.Force = requests.NewBoolean(d.Get("force").(bool))
	request.InstanceType = EcsInstance
	request.ClientToken = buildClientToken(request.GetActionName())

	if strings.HasPrefix(instanceId, "lb-") {
		request.InstanceType = SlbInstance
	}
	if strings.HasPrefix(instanceId, "ngw-") {
		request.InstanceType = Nat
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = instanceType.(string)
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateEipAddress(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "IncorrectHaVipStatus", "TaskConflict",
				"InvalidIpStatus.HasBeenUsedBySnatTable", "InvalidIpStatus.HasBeenUsedByForwardEntry"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return WrapError(vpcService.WaitForEipAssociation(d.Id(), Available, DefaultTimeoutMedium))
}

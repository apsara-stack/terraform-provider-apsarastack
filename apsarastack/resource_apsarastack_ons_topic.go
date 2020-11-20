package apsarastack

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackOnsTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackOnsTopicCreate,
		Read:   resourceApsaraStackOnsTopicRead,
		Update: resourceApsaraStackOnsTopicUpdate,
		Delete: resourceApsaraStackOnsTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"message_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 4, 5}),
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"perm": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 4, 6}),
			},
		},
	}
}

func resourceApsaraStackOnsTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	onsService := OnsService{client}

	instanceId := d.Get("instance_id").(string)
	topic := d.Get("topic").(string)

	request := ons.CreateOnsTopicCreateRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ons", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.Topic = topic
	request.InstanceId = instanceId
	request.MessageType = requests.NewInteger(d.Get("message_type").(int))

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = v.(string)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsTopicCreate(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ons_topic", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	d.SetId(instanceId + ":" + topic)

	if err = onsService.WaitForOnsTopic(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceApsaraStackOnsTopicRead(d, meta)
}

func resourceApsaraStackOnsTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	onsService := OnsService{client}

	object, err := onsService.DescribeOnsTopic(d.Id())
	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("topic", object.Topic)
	d.Set("message_type", object.MessageType)
	d.Set("remark", object.Remark)

	return nil
}

func resourceApsaraStackOnsTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	onsService := OnsService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := ons.CreateOnsTopicUpdateRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ons", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.InstanceId = instanceId
	request.Topic = topic

	var perm int
	if d.HasChange("perm") {
		perm = d.Get("perm").(int)
		request.Perm = requests.NewInteger(perm)
		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsTopicUpdate(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceApsaraStackOnsTopicRead(d, meta)
}

func resourceApsaraStackOnsTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	onsService := OnsService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := ons.CreateOnsTopicDeleteRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ons", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.Topic = topic
	request.InstanceId = instanceId

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsTopicDelete(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	return WrapError(onsService.WaitForOnsTopic(d.Id(), Deleted, DefaultTimeoutMedium))
}

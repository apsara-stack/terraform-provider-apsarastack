package apsarastack

import (
	"strings"

	"github.com/denverdino/aliyungo/common"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackSlb() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackSlbCreate,
		Read:   resourceApsaraStackSlbRead,
		Update: resourceApsaraStackSlbUpdate,
		Delete: resourceApsaraStackSlbDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 80),
				Default:      resource.PrefixedUniqueId("tf-lb-"),
			},
			"address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"internet", "intranet"}, false),
			},

			"vswitch_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: slbInternetDiffSuppressFunc,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},

			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36})),
			},
		},
	}
}

func resourceApsaraStackSlbCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	slbService := SlbService{client}
	request := slb.CreateCreateLoadBalancerRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerName = d.Get("name").(string)
	request.AddressType = strings.ToLower(string(Intranet))
	request.InternetChargeType = strings.ToLower(string(PayByTraffic))
	request.ClientToken = buildClientToken(request.GetActionName())

	if v, ok := d.GetOk("address_type"); ok && v.(string) != "" {
		request.AddressType = strings.ToLower(v.(string))
	}

	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		request.VSwitchId = v.(string)
	}

	if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) != "" {
		request.PayType = v.(string)
		if request.PayType == string(PrePaid) {
			request.PayType = "PrePay"
		} else {
			request.PayType = "PayOnDemand"
		}
	}

	if request.PayType == string("PrePay") {
		period := d.Get("period").(int)
		request.Duration = requests.NewInteger(period)
		request.PricingCycle = string(Month)
		if period > 9 {
			request.Duration = requests.NewInteger(period / 12)
			request.PricingCycle = string(Year)
		}
		request.AutoPay = requests.NewBoolean(true)
	}
	var raw interface{}

	invoker := Invoker{}
	invoker.AddCatcher(Catcher{"OperationFailed.TokenIsProcessing", 10, 5})

	if err := invoker.Run(func() error {
		resp, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.CreateLoadBalancer(request)
		})
		raw = resp
		return err
	}); err != nil {
		if IsExpectedErrors(err, []string{"OrderFailed"}) {
			return WrapError(err)
		}
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_slb", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.CreateLoadBalancerResponse)
	d.SetId(response.LoadBalancerId)

	if err := slbService.WaitForSlb(response.LoadBalancerId, Active, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceApsaraStackSlbUpdate(d, meta)
}

func resourceApsaraStackSlbRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlb(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.LoadBalancerName)
	d.Set("address_type", object.AddressType)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("address", object.Address)
	d.Set("instance_charge_type", object.PayType)

	if object.PayType == "PrePay" {
		d.Set("instance_charge_type", PrePaid)
		period, err := computePeriodByUnit(object.CreateTime, object.EndTime, d.Get("period").(int), "Month")
		if err != nil {
			return WrapError(err)
		}
		d.Set("period", period)
	} else {
		d.Set("instance_charge_type", PostPaid)
	}
	tags, _ := slbService.DescribeTags(d.Id(), nil, TagResourceInstance)
	if len(tags) > 0 {
		if err := d.Set("tags", slbService.tagsToMap(tags)); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceApsaraStackSlbUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	slbService := SlbService{client}
	d.Partial(true)

	// set instance tags
	if err := slbService.setInstanceTags(d, TagResourceInstance); err != nil {
		return WrapError(err)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceApsaraStackSlbRead(d, meta)
	}

	if d.HasChange("name") {
		request := slb.CreateSetLoadBalancerNameRequest()
		request.RegionId = client.RegionId
		request.LoadBalancerId = d.Id()
		request.LoadBalancerName = d.Get("name").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetLoadBalancerName(request)
		})
		if err != nil {
			WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("name")
	}
	update := false
	modifyLoadBalancerInternetSpecRequest := slb.CreateModifyLoadBalancerInternetSpecRequest()
	modifyLoadBalancerInternetSpecRequest.RegionId = client.RegionId
	modifyLoadBalancerInternetSpecRequest.LoadBalancerId = d.Id()
	if update {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ModifyLoadBalancerInternetSpec(modifyLoadBalancerInternetSpecRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyLoadBalancerInternetSpecRequest.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(modifyLoadBalancerInternetSpecRequest.GetActionName(), raw, modifyLoadBalancerInternetSpecRequest.RpcRequest, modifyLoadBalancerInternetSpecRequest)
	}

	update = false
	modifyLoadBalancerPayTypeRequest := slb.CreateModifyLoadBalancerPayTypeRequest()
	modifyLoadBalancerPayTypeRequest.RegionId = client.RegionId
	modifyLoadBalancerPayTypeRequest.LoadBalancerId = d.Id()
	if d.HasChange("instance_charge_type") {
		modifyLoadBalancerPayTypeRequest.PayType = d.Get("instance_charge_type").(string)
		if modifyLoadBalancerPayTypeRequest.PayType == string(PrePaid) {
			modifyLoadBalancerPayTypeRequest.PayType = "PrePay"
		} else {
			modifyLoadBalancerPayTypeRequest.PayType = "PayOnDemand"
		}
		if modifyLoadBalancerPayTypeRequest.PayType == "PrePay" {
			period := d.Get("period").(int)
			modifyLoadBalancerPayTypeRequest.Duration = requests.NewInteger(period)
			modifyLoadBalancerPayTypeRequest.PricingCycle = string(Month)
			if period > 9 {
				modifyLoadBalancerPayTypeRequest.Duration = requests.NewInteger(period / 12)
				modifyLoadBalancerPayTypeRequest.PricingCycle = string(Year)
			}
			modifyLoadBalancerPayTypeRequest.AutoPay = requests.NewBoolean(true)
		}
		update = true
		d.SetPartial("instance_charge_type")
	}

	if update {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ModifyLoadBalancerPayType(modifyLoadBalancerPayTypeRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyLoadBalancerPayTypeRequest.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(modifyLoadBalancerPayTypeRequest.GetActionName(), raw, modifyLoadBalancerPayTypeRequest.RpcRequest, modifyLoadBalancerPayTypeRequest)
	}
	d.Partial(false)

	return resourceApsaraStackSlbRead(d, meta)
}

func resourceApsaraStackSlbDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	slbService := SlbService{client}

	request := slb.CreateDeleteLoadBalancerRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerId = d.Id()

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DeleteLoadBalancer(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(slbService.WaitForSlb(d.Id(), Deleted, DefaultTimeoutMedium))
}

package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackOnsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackOnsGroupCreate,
		Read:   resourceApsaraStackOnsGroupRead,
		Update: resourceApsaraStackOnsGroupUpdate,
		Delete: resourceApsaraStackOnsGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOnsGroupId,
			},
			"remark": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"read_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackOnsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	onsService := OnsService{client}
	var requestInfo *ecs.Client

	instanceId := d.Get("instance_id").(string)
	groupId := d.Get("group_id").(string)
	remark := d.Get("remark").(string)

	check, err := onsService.DescribeOnsGroup(groupId, instanceId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ons_group", "Group alreadyExist", ApsaraStackSdkGoERROR)
	}
	if len(check.Data) == 0 {

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ons-inner",
			"Action":          "ConsoleGroupCreate",
			"Version":         "2018-02-05",
			"ProductName":     "Ons-inner",
			"PreventCache":    "",
			"GroupId":         groupId,
			"Remark":          remark,
			"OnsRegionId":     client.RegionId,
			"InstanceId":      instanceId,
		}
		request.Method = "POST"
		request.Product = "Ons-inner"
		request.Version = "2018-02-05"
		request.ServiceCode = "Ons-inner"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.SetHTTPSInsecure(true)
		request.ApiName = "ConsoleGroupCreate"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ons_group", "ConsoleGroupCreate", raw)
		}
		addDebug("ConsoleGroupCreate", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ons_group", "ConsoleGroupCreate", ApsaraStackSdkGoERROR)
		}
		addDebug("ConsoleGroupCreate", raw, requestInfo, bresponse.GetHttpContentString())
	}
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		check, err = onsService.DescribeOnsGroup(groupId, instanceId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(check.Data) != 0 {
			return nil
		}
		return resource.RetryableError(Error("New Group has been created successfully."))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ons_group", "Failed to create ONS Group", ApsaraStackSdkGoERROR)
	}

	d.SetId(groupId)

	return resourceApsaraStackOnsGroupRead(d, meta)
}

func resourceApsaraStackOnsGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	onsService := OnsService{client}
	instanceId := d.Get("instance_id").(string)

	object, err := onsService.DescribeOnsGroup(d.Id(), instanceId)

	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.Data[0].NamespaceID)
	d.Set("group_id", object.Data[0].ID)
	d.Set("remark", object.Data[0].Remark)

	return nil
}

func resourceApsaraStackOnsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceApsaraStackOnsGroupRead(d, meta)
}

func resourceApsaraStackOnsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	onsService := OnsService{client}
	var requestInfo *ecs.Client
	instanceId := d.Get("instance_id").(string)

	check, err := onsService.DescribeOnsGroup(d.Id(), instanceId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsGroupExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsGroupExist", check, requestInfo, map[string]string{"GroupId": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ons-inner",
			"Action":          "ConsoleGroupDelete",
			"Version":         "2018-02-05",
			"ProductName":     "Ons-inner",
			"PreventCache":    "",
			"GroupId":         d.Id(),
			"OnsRegionId":     client.RegionId,
			"InstanceId":      instanceId,
		}

		request.Method = "POST"
		request.Product = "Ons-inner"
		request.Version = "2018-02-05"
		request.ServiceCode = "Ons-inner"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "ConsoleGroupDelete"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err = onsService.DescribeOnsGroup(d.Id(), instanceId)

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return nil
}

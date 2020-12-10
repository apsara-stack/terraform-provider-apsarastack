package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceApsaraStackAscmResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmResourceGroupCreate,
		Read:   resourceApsaraStackAscmResourceGroupRead,
		Update: resourceApsaraStackAscmResourceGroupUpdate,
		Delete: resourceApsaraStackAscmResourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rg_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackAscmResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client

	name := d.Get("name").(string)
	check, err := ascmService.DescribeAscmResourceGroup(name)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_resource_group", "RG alreadyExist", ApsaraStackSdkGoERROR)
	}
	organizationid := d.Get("organization_id").(string)

	if len(check.Data) == 0 {

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":            client.RegionId,
			"AccessKeySecret":     client.SecretKey,
			"Product":             "Ascm",
			"Action":              "CreateResourceGroup",
			"Version":             "2019-05-10",
			"ProductName":         "ascm",
			"resource_group_name": name,
			"organization_id":     organizationid,
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "CreateResourceGroup"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_resource_group", "CreateResourceGroup", raw)
		}
		addDebug("CreateResourceGroup", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_resource_group", "CreateResourceGroup", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateResourceGroup", raw, requestInfo, bresponse.GetHttpContentString())
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		check, err = ascmService.DescribeAscmResourceGroup(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(check.Data) != 0 {
			return nil
		}
		return resource.RetryableError(Error("New Resource Group has been created successfully."))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_resource_group", "Failed to create resource set", ApsaraStackSdkGoERROR)
	}

	d.SetId(check.Data[0].ResourceGroupName)

	return resourceApsaraStackAscmResourceGroupUpdate(d, meta)

}

func resourceApsaraStackAscmResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceApsaraStackAscmResourceGroupRead(d, meta)

}

func resourceApsaraStackAscmResourceGroupRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmResourceGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if len(object.Data) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("rg_id", object.Data[0].ID)
	d.Set("name", object.Data[0].ResourceGroupName)
	d.Set("organization_id", object.Data[0].OrganizationID)

	return nil
}
func resourceApsaraStackAscmResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmResourceGroup(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsResourceGroupExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsResourceGroupExist", check, requestInfo, map[string]string{"resourceGroupName": d.Id()})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":          client.RegionId,
			"AccessKeySecret":   client.SecretKey,
			"Product":           "ascm",
			"Action":            "RemoveResourceGroup",
			"Version":           "2019-05-10",
			"ProductName":       "ascm",
			"resourceGroupName": d.Id(),
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "RemoveResourceGroup"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err := ascmService.DescribeAscmResourceGroup(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		if check.Data[0].ResourceGroupName != "" {
			return resource.RetryableError(Error("Trying to delete Resource Group %#v successfully.", d.Id()))
		}
		return nil
	})
	return nil
}

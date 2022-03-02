package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strconv"
	"strings"
	"time"
)

func resourceApsaraStackAscmUserGroupResourceSetBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmUserGroupResourceSetBindingCreate,
		Read:   resourceApsaraStackAscmUserGroupResourceSetBindingRead,
		Delete: resourceApsaraStackAscmUserGroupResourceSetBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"resource_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApsaraStackAscmUserGroupResourceSetBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client

	resourceSetId := d.Get("resource_set_id").(string)
	userGroupId := d.Get("user_group_id").(string)

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "AddResourceSetToUserGroup"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Product":         "Ascm",
		"Action":          "AddResourceSetToUserGroup",
		"Version":         "2019-05-10",
		"ProductName":     "ascm",
		"ascmRoleId":      "2",
		"userGroupId":     userGroupId,
		"resourceSetId":   resourceSetId,
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw AddResourceSetToUserGroup is : %s", raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user_group_resource_set_binding", "AddResourceSetToUserGroup", raw)
	}
	addDebug("AddResourceSetToUserGroup", raw, requestInfo, request)
	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user_group_resource_set_binding", "AddResourceSetToUserGroup", ApsaraStackSdkGoERROR)
	}
	addDebug("AddResourceSetToUserGroup", raw, requestInfo, bresponse.GetHttpContentString())

	d.SetId(resourceSetId)
	return resourceApsaraStackAscmUserGroupResourceSetBindingRead(d, meta)
}

func resourceApsaraStackAscmUserGroupResourceSetBindingRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)

	ascmService := &AscmService{client: client}
	obj, err := ascmService.DescribeAscmUserGroupResourceSetBinding(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("resource_set_id", strconv.Itoa(obj.Data[0].Id))

	return nil
}
func resourceApsaraStackAscmUserGroupResourceSetBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmUserGroupResourceSetBinding(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBindingExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsBindingExist", check, requestInfo, map[string]string{"resourceGroupId": d.Id()})
	userGroupId := d.Get("user_group_id").(string)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "RemoveResourceSetFromUserGroup",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"userGroupId":     userGroupId,
			"resourceSetId":   d.Id(),
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "RemoveResourceSetFromUserGroup"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}

		addDebug("RemoveResourceSetFromUserGroup", raw, request)
		_, err = ascmService.DescribeAscmUserGroupResourceSetBinding(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}

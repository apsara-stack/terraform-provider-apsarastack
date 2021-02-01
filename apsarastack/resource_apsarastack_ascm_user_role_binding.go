package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
	"time"
)

func resourceApsaraStackAscmUserRoleBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmUserRoleBindingCreate,
		Read:   resourceApsaraStackAscmUserRoleBindingRead,
		Update: resourceApsaraStackAscmUserRoleBindingUpdate,
		Delete: resourceApsaraStackAscmUserRoleBindingDelete,
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApsaraStackAscmUserRoleBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client
	lname := d.Get("login_name").(string)
	roleid := d.Get("role_id")
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Product":         "Ascm",
		"Action":          "AddRoleToUser",
		"Version":         "2019-05-10",
		"ProductName":     "ascm",
		"LoginName":       lname,
		"RoleId":          fmt.Sprint(roleid),
	}
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
	request.ApiName = "AddRoleToUser"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user_role_binding", "AddRoleToUser", raw)
	}

	addDebug("AddRoleToUser", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user_role_binding", "AddRoleToUser", ApsaraStackSdkGoERROR)
	}
	addDebug("AddRoleToUser", raw, requestInfo, bresponse.GetHttpContentString())

	d.SetId(lname)

	return resourceApsaraStackAscmUserRoleBindingRead(d, meta)
}

func resourceApsaraStackAscmUserRoleBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUser(d.Id())
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
	d.Set("login_name", object.Data[0].LoginName)

	return nil
}

func resourceApsaraStackAscmUserRoleBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceApsaraStackAscmUserRoleBindingCreate(d, meta)

}

func resourceApsaraStackAscmUserRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	roleid := d.Get("role_id").(int)
	check, err := ascmService.DescribeAscmUser(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBindingExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsBindingExist", check, requestInfo, map[string]string{"loginName": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "RemoveRoleFromUser",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"LoginName":       d.Id(),
			"RoleId":          fmt.Sprint(roleid),
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
		request.ApiName = "RemoveRoleFromUser"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err = ascmService.DescribeAscmUser(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}

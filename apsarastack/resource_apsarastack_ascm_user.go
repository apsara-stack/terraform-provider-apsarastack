package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
	"time"
)

func resourceApsaraStackAscmUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmUserCreate,
		Read:   resourceApsaraStackAscmUserRead,
		Update: resourceApsaraStackAscmUserUpdate,
		Delete: resourceApsaraStackAscmUserDelete,
		Schema: map[string]*schema.Schema{
			"cellphone_number": {
				Type:     schema.TypeString,
				Required: true,
			},
			"telephone_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mobile_nation_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"login_policy_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackAscmUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	lname := d.Get("login_name").(string)
	check, err := ascmService.DescribeAscmUser(lname)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user", "User alreadyExist", ApsaraStackSdkGoERROR)
	}
	dname := d.Get("display_name").(string)
	email := d.Get("email").(string)
	cellnum := d.Get("cellphone_number").(string)
	mobnationcode := d.Get("mobile_nation_code").(string)
	organizationid := d.Get("organization_id").(string)

	if len(check.Data) == 0 {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":         client.RegionId,
			"AccessKeySecret":  client.SecretKey,
			"Product":          "Ascm",
			"Action":           "AddUser",
			"Version":          "2019-05-10",
			"ProductName":      "ascm",
			"loginName":        lname,
			"displayName":      dname,
			"cellphoneNum":     cellnum,
			"mobileNationCode": mobnationcode,
			"email":            email,
			"organizationId":   organizationid,
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
		request.ApiName = "AddUser"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user", "AddUser", raw)
		}
		addDebug("AddUser", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user", "AddUser", ApsaraStackSdkGoERROR)
		}
		addDebug("AddUser", raw, requestInfo, bresponse.GetHttpContentString())
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		check, err = ascmService.DescribeAscmUser(lname)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(check.Data) != 0 {
			return nil
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user", "Failed to add User", ApsaraStackSdkGoERROR)
	}

	d.SetId(check.Data[0].LoginName)

	return resourceApsaraStackAscmUserUpdate(d, meta)

}

func resourceApsaraStackAscmUserUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	lname := d.Get("login_name").(string)
	organizationid := d.Get("organization_id").(string)

	var dname, cellnum, mobnationcode, email string

	if d.HasChange("display_name") {
		dname = d.Get("display_name").(string)
	}
	if d.HasChange("cellphone_number") {
		cellnum = d.Get("cellphone_number").(string)
	}
	if d.HasChange("mobile_nation_code") {
		mobnationcode = d.Get("mobile_nation_code").(string)
	} else {
		mobnationcode = d.Get("mobile_nation_code").(string)
	}
	if d.HasChange("email") {
		email = d.Get("email").(string)
	}

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}

	request.QueryParams = map[string]string{
		"RegionId":         client.RegionId,
		"AccessKeySecret":  client.SecretKey,
		"Product":          "ascm",
		"Department":       client.Department,
		"ResourceGroup":    client.ResourceGroup,
		"Action":           "ModifyUserInformation",
		"Version":          "2019-05-10",
		"ProductName":      "ascm",
		"loginName":        lname,
		"displayName":      dname,
		"cellphoneNum":     cellnum,
		"mobileNationCode": mobnationcode,
		"email":            email,
		"organization_id":  organizationid,
	}
	request.Domain = client.Domain
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ModifyUserInformation"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user", "ModifyUserInformationRequestFailed", raw)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	if !bresponse.IsSuccess() {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user", "ModifyUserInformationFailed", raw)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), bresponse)
	if err != nil {
		return WrapError(err)
	}
	object, err := ascmService.DescribeAscmUser(lname)
	if err != nil {
		return WrapError(err)
	}
	d.SetId(object.Data[0].LoginName)

	return resourceApsaraStackAscmUserRead(d, meta)

}

func resourceApsaraStackAscmUserRead(d *schema.ResourceData, meta interface{}) error {

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
	d.Set("display_name", object.Data[0].DisplayName)
	d.Set("email", object.Data[0].Email)
	d.Set("mobile_nation_code", object.Data[0].MobileNationCode)
	d.Set("cellphone_number", object.Data[0].CellphoneNum)
	d.Set("organization_id", object.Data[0].Organization.ID)

	return nil
}
func resourceApsaraStackAscmUserDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmUser(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsUserExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsUserExist", check, requestInfo, map[string]string{"loginName": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "RemoveUserByLoginName",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"loginName":       d.Id(),
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
		request.ApiName = "RemoveUserByLoginName"
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
		return resource.RetryableError(err)
	})
	return nil
}

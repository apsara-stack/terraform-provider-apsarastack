package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strings"
	"time"
)

func resourceApsaraStackAscmRamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmRamRoleCreate,
		Read:   resourceApsaraStackAscmRamRoleRead,
		Update: resourceApsaraStackAscmRamRoleUpdate,
		Delete: resourceApsaraStackAscmRamRoleDelete,
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization_visibility": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceApsaraStackAscmRamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client

	name := d.Get("role_name").(string)
	description := d.Get("description").(string)
	organizationvisibility := d.Get("organization_visibility").(string)

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "CreateRole",
		"Version":         "2019-05-10",
		"ProductName":     "ascm",
		"roleName":        name,
		"description":     description,
		"roleRange":       "roleRange.userGroup",
		"roleType":        "ROLETYPE_RAM",
		"organizationVisibility":/*fmt.Sprintf("organizationVisibility.%s", strings.ToLower(*/ organizationvisibility,
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
	request.ApiName = "CreateRole"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("raw %s", raw)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_ram_role", "CreateRole", raw)
	}
	addDebug("CreateRole", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_ram_role", "CreateRole", ApsaraStackSdkGoERROR)
	}
	addDebug("CreateRole", raw, requestInfo, bresponse.GetHttpContentString())
	//}

	log.Printf("rolename %s", name)
	d.SetId(name)
	return resourceApsaraStackAscmRamRoleUpdate(d, meta)

}

func resourceApsaraStackAscmRamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*connectivity.ApsaraStackClient)
	//ascmService := AscmService{client}
	//attributeUpdate := false
	//check, err := ascmService.DescribeAscmRamRole(d.Id())
	//
	//if err != nil {
	//	return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsRamRoleExist", ApsaraStackSdkGoERROR)
	//}
	//var rname, desc string
	//if d.HasChange("role_name") {
	//	if v, ok := d.GetOk("role_name"); ok {
	//		rname = v.(string)
	//	}
	//	check.Data[0].RoleName = rname
	//	check.Data[0].NewRoleName = rname
	//	attributeUpdate = true
	//} else {
	//	if v, ok := d.GetOk("role_name"); ok {
	//		rname = v.(string)
	//	}
	//	check.Data[0].RoleName = rname
	//}
	//if d.HasChange("description") {
	//	if v, ok := d.GetOk("description"); ok {
	//		desc = v.(string)
	//	}
	//	check.Data[0].Description = desc
	//	check.Data[0].NewDescription = desc
	//	attributeUpdate = true
	//} else {
	//	if v, ok := d.GetOk("description"); ok {
	//		desc = v.(string)
	//	}
	//	check.Data[0].Description = desc
	//}
	//
	//request := requests.NewCommonRequest()
	//request.QueryParams = map[string]string{
	//	"RegionId":          client.RegionId,
	//	"AccessKeySecret":   client.SecretKey,
	//	"Department":        client.Department,
	//	"ResourceGroup":     client.ResourceGroup,
	//	"Product":           "ascm",
	//	"Action":            "UpdateRoleInfo",
	//	"Version":           "2019-05-10",
	//	"newRoleName": rname,
	//	//"RoleName": rname,
	//	"newDescription": desc,
	//	//"Description": desc,
	//	"id":          fmt.Sprint(check.Data[0].ID),
	//}
	//request.Method = "POST"
	//request.Product = "ascm"
	//request.Version = "2019-05-10"
	//request.ServiceCode = "ascm"
	//request.Domain = client.Domain
	//if strings.ToLower(client.Config.Protocol) == "https" {
	//	request.Scheme = "https"
	//} else {
	//	request.Scheme = "http"
	//}
	//request.SetHTTPSInsecure(true)
	//request.ApiName = "UpdateRoleInfo"
	//request.RegionId = client.RegionId
	//request.Headers = map[string]string{"RegionId": client.RegionId}
	//
	//if attributeUpdate {
	//	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
	//		return ecsClient.ProcessCommonRequest(request)
	//	})
	//	if err != nil {
	//		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_ram_role", "UpdateRoleInfo", raw)
	//	}
	//	addDebug(request.GetActionName(), raw, request)
	//	//d.Set("role_name", check.Data[0].NewRoleName)
	//	//d.Set("role_id", check.Data[0].ID)
	//	//d.Set("description", check.Data[0].NewDescription)
	//}
	//d.SetId(rname)

	return resourceApsaraStackAscmRamRoleRead(d, meta)

}

func resourceApsaraStackAscmRamRoleRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmRamRole(d.Id())
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
	visibility := d.Get("organization_visibility").(string)

	visibility1 := strings.Trim(object.Data[0].OrganizationVisibility, "organizationVisibility.")
	if visibility1 == visibility {
		d.Set("organization_visibility", visibility)
		d.Set("role_name", object.Data[0].RoleName)
		d.Set("role_id", object.Data[0].ID)
		d.Set("description", object.Data[0].Description)
	} else {
		d.Set("organization_visibility", visibility)
	}
	//log.Printf("suraj visibility %s",visibility)
	//d.Set("role_name", object.Data[0].RoleName)
	//d.Set("role_id", object.Data[0].ID)
	//d.Set("description", object.Data[0].Description)
	return nil
}

func resourceApsaraStackAscmRamRoleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmRamRole(d.Id())

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsRamRoleExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsRamRoleExist", check, requestInfo, map[string]string{"roleName": d.Id()})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "RemoveRole",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"roleName":        d.Id(),
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
		request.ApiName = "RemoveRole"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		_, err = ascmService.DescribeAscmRamRole(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}

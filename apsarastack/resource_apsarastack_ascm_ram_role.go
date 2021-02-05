package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
	ascmService := AscmService{client}
	name := d.Get("role_name").(string)
	description := d.Get("description").(string)
	organizationvisibility := d.Get("organization_visibility").(string)
	check, err := ascmService.DescribeAscmRamRole(name)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_ram_policy", "policy alreadyExist", ApsaraStackSdkGoERROR)
	}
	if len(check.Data) == 0 {
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":               client.RegionId,
			"AccessKeySecret":        client.SecretKey,
			"Department":             client.Department,
			"ResourceGroup":          client.ResourceGroup,
			"Product":                "ascm",
			"Action":                 "CreateRole",
			"Version":                "2019-05-10",
			"ProductName":            "ascm",
			"roleName":               name,
			"description":            description,
			"roleRange":              "roleRange.userGroup",
			"roleType":               "ROLETYPE_RAM",
			"organizationVisibility": organizationvisibility,
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

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_ram_role", "CreateRole", raw)
		}
		addDebug("CreateRole", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_ram_role", "CreateRole", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateRole", raw, requestInfo, bresponse.GetHttpContentString())
	}
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		check, err = ascmService.DescribeAscmRamRole(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	d.SetId(name + COLON_SEPARATED + fmt.Sprint(check.Data[0].ID))
	return resourceApsaraStackAscmRamRoleUpdate(d, meta)

}

func resourceApsaraStackAscmRamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceApsaraStackAscmRamRoleRead(d, meta)
}

func resourceApsaraStackAscmRamRoleRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmRamRole(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("role_name", did[0])
	d.Set("organization_visibility", object.Data[0].OrganizationVisibility)
	d.Set("role_id", object.Data[0].ID)
	d.Set("description", object.Data[0].Description)
	return nil
}

func resourceApsaraStackAscmRamRoleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmRamRole(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsRamRoleExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsRamRoleExist", check, requestInfo, map[string]string{"roleName": did[0]})
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
			"roleName":        did[0],
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

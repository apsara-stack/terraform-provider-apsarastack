package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
	"time"
)

func resourceApsaraStackAscmRamPolicyForRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmRamPolicyForRoleCreate,
		Read:   resourceApsaraStackAscmRamPolicyForRoleRead,
		Update: resourceApsaraStackAscmRamPolicyForRoleUpdate,
		Delete: resourceApsaraStackAscmRamPolicyForRoleDelete,
		Schema: map[string]*schema.Schema{
			"ram_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApsaraStackAscmRamPolicyForRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client
	ram_id := d.Get("ram_policy_id").(string)
	roleid := d.Get("role_id").(string)
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Product":         "Ascm",
		"Action":          "AddRAMPolicyToRole",
		"Version":         "2019-05-10",
		"ProductName":     "ascm",
		"ramPolicyId":     ram_id,
		"roleId":          roleid,
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
	request.ApiName = "AddRAMPolicyToRole"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("Suraj raw %s", raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_ram_policy_for_role", "AddRAMPolicyToRole", raw)
	}

	addDebug("AddRAMPolicyToRole", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_ram_policy_for_role", "AddRAMPolicyToRole", ApsaraStackSdkGoERROR)
	}
	addDebug("AddRAMPolicyToRole", raw, requestInfo, bresponse.GetHttpContentString())

	d.SetId(ram_id)

	return resourceApsaraStackAscmRamPolicyForRoleRead(d, meta)
}

func resourceApsaraStackAscmRamPolicyForRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmRamPolicy(d.Id())
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
	d.Set("ram_policy_id", object.Data[0].RamPolicyId)
	d.Set("role_id", object.Data[0].RoleId)

	return nil
}

func resourceApsaraStackAscmRamPolicyForRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceApsaraStackAscmRamPolicyForRoleCreate(d, meta)

}

func resourceApsaraStackAscmRamPolicyForRoleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	roleid := d.Get("role_id").(string)
	check, err := ascmService.DescribeAscmRamPolicy(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBindingExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsBindingExist", check, requestInfo, map[string]string{"ramPolicyId": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "RemoveRAMPolicyFromRole",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"ramPolicyId":     d.Id(),
			"roleId":          roleid,
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
		request.ApiName = "RemoveRAMPolicyFromRole"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err = ascmService.DescribeAscmRamPolicy(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}

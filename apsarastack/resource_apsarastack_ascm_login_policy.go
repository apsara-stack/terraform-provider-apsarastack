package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceApsaraStackLogInPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackLogInPolicyCreate,
		Read:   resourceApsaraStackLogInPolicyRead,
		Update: resourceApsaraStackLogInPolicyUpdate,
		Delete: resourceApsaraStackLogInPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ALLOW",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW", "DENY"}, false),
			},
		},
	}
}
func resourceApsaraStackLogInPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client
	ascmService := AscmService{client}
	name := d.Get("name").(string)
	descr := d.Get("description").(string)
	rule := d.Get("rule").(string)
	object, err := ascmService.ListLoginPolicies(name)
	if err != nil {

		return WrapError(err)
	}
	if len(object.Data) == 0 {

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":         client.RegionId,
			"AccessKeySecret":  client.SecretKey,
			"Product":          "ascm",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"Action":           "AddLoginPolicy",
			"AccountInfo":      "123456",
			"Version":          "2019-05-10",
			"SignatureVersion": "1.0",
			"ProductName":      "ascm",
			"Name":             name,
			"Description":      descr,
			"Rule":             rule,
		}
		request.Domain = client.Domain
		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Scheme = "http"
		request.ApiName = "AddLoginPolicy"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_login_policy", "AddLoginPolicy", raw)
		}
		addDebug("AddLoginPolicy", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_login_policy", "AddLoginPolicy", ApsaraStackSdkGoERROR)
		}
		addDebug("AddLoginPolicy", raw, requestInfo, bresponse.GetHttpContentString())
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		object, err = ascmService.ListLoginPolicies(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(object.Data) != 0 {
			return nil
		}
		return resource.RetryableError(Error("New Login Policy has been added successfully."))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_login_policy", "Failed to add login Policy", ApsaraStackSdkGoERROR)
	}

	d.SetId(object.Data[0].Name)

	return resourceApsaraStackLogInPolicyUpdate(d, meta)
}

func resourceApsaraStackLogInPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	request := requests.NewCommonRequest()
	var name, rule, desc string
	if d.HasChange("name") {
		name = d.Get("name").(string)
	}
	if d.HasChange("rule") {
		rule = d.Get("rule").(string)
	}
	if d.HasChange("description") {
		desc = d.Get("description").(string)
	}
	policyId := fmt.Sprint(d.Get("policy_id").(int))

	request.QueryParams = map[string]string{
		"RegionId":         client.RegionId,
		"AccessKeySecret":  client.SecretKey,
		"Product":          "ascm",
		"Department":       client.Department,
		"ResourceGroup":    client.ResourceGroup,
		"Action":           "ModifyLoginPolicy",
		"AccountInfo":      "123456",
		"Version":          "2019-05-10",
		"SignatureVersion": "1.0",
		"ProductName":      "ascm",
		"id":               policyId,
		"Name":             name,
		"Rule":             rule,
		"Description":      desc,
	}
	request.Domain = client.Domain
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Scheme = "http"
	request.ApiName = "ModifyLoginPolicy"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_login_policy", "LoginPolicyUpdateRequestFailed", raw)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	if !bresponse.IsSuccess() {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_login_policy", "LoginPolicyUpdateFailed", raw)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), bresponse)
	if err != nil {
		return WrapError(err)
	}
	object, err := ascmService.ListLoginPolicies(name)
	if err != nil {
		return WrapError(err)
	}

	d.SetId(object.Data[0].Name)

	return resourceApsaraStackLogInPolicyRead(d, meta)
}
func resourceApsaraStackLogInPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.ListLoginPolicies(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Data[0].Name)
	d.Set("description", object.Data[0].Description)
	d.Set("policy_id", object.Data[0].ID)
	d.Set("rule", object.Data[0].Rule)
	return nil
}
func resourceApsaraStackLogInPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client

	name := d.Get("name").(string)

	check, err := ascmService.ListLoginPolicies(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsLoginPolicyExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsLoginPolicyExist", check, requestInfo, map[string]string{"loginpolicyName": d.Id()})
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":         client.RegionId,
			"AccessKeySecret":  client.SecretKey,
			"Product":          "ascm",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"Action":           "RemoveLoginPolicyByName",
			"AccountInfo":      "123456",
			"Version":          "2019-05-10",
			"SignatureVersion": "1.0",
			"ProductName":      "ascm",
			"Name":             name,
		}
		request.Domain = client.Domain
		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Scheme = "http"
		request.ApiName = "RemoveLoginPolicyByName"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return resource.RetryableError(err)
		}

		_, err = ascmService.ListLoginPolicies(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	return nil
}

type IPRanges []struct {
	IPRange       string `json:"ipRange"`
	LoginPolicyID int    `json:"loginPolicyId"`
	Protocol      string `json:"protocol"`
}

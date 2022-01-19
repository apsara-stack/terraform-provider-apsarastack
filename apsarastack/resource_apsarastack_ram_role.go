package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strings"
)

func resourceApsaraStackRamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackRamRoleCreate,
		Read:   resourceApsaraStackRamRoleRead,
		Delete: resourceApsaraStackRamRoleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"role_policy_document": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}
func resourceApsaraStackRamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response *CreatePolicyResponse
	name := d.Get("name").(string)

	policyDoc := d.Get("role_policy_document").(string)
	description := d.Get("description").(string)
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":                 client.RegionId,
		"AccessKeySecret":          client.SecretKey,
		"Department":               client.Department,
		"ResourceGroup":            client.ResourceGroup,
		"Product":                  "ram",
		"Action":                   "CreateRole",
		"Version":                  "2015-05-01",
		"RoleName":                 name,
		"Description":              description,
		"AssumeRolePolicyDocument": policyDoc,
	}
	request.Method = "POST"
	request.Product = "Ram"
	request.Version = "2015-05-01"
	request.ServiceCode = "Ram"
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
	log.Printf(" response of CreateRole : %s", raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ram_role", "CreateRole", raw)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	if !bresponse.IsSuccess() || bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ram_role", "CreateRole", raw)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return WrapError(err)
	}
	//d.Set("policy_type", response.Policy.PolicyType)
	d.SetId(name)
	return resourceApsaraStackRamRoleRead(d, meta)
}
func resourceApsaraStackRamRoleRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*connectivity.ApsaraStackClient)
	//ascmService := AscmService{client}
	//policy, err := ascmService.DescribeRamRole(d.Id(), d.Get("policy_type").(string))
	//if err != nil {
	//	return WrapError(err)
	//}

	return nil
}
func resourceApsaraStackRamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	name := d.Get("name").(string)

	//policyType := d.Get("policy_type").(string)
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "ram",
		"Action":          "DeleteRole",
		"Version":         "2015-05-01",
		"RoleName":        name,
	}
	request.Method = "POST"
	request.Product = "Ram"
	request.Version = "2015-05-01"
	request.ServiceCode = "Ram"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DeleteRole"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf(" response of DeleteRole : %s", raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ram_role", "DeleteRole", raw)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	if !bresponse.IsSuccess() || bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ram_role", "DeleteRole", raw)
	}

	return nil
}

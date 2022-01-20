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

func resourceApsaraStackRamPolicyRoleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackRamPolicyRoleAttachmentCreate,
		Read:   resourceApsaraStackRamPolicyRoleAttachmentRead,
		Delete: resourceApsaraStackRamPolicyRoleAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},
			"policy_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},
			"policy_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Custom",
				ForceNew: true,

				//ValidateFunc: validation.StringLenBetween(3, 64),
			},
		},
	}
}
func resourceApsaraStackRamPolicyRoleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response *CreatePolicyResponse
	RoleName := d.Get("role_name").(string)

	PolicyType := d.Get("policy_type").(string)
	PolicyName := d.Get("policy_name").(string)
	//Description := d.Get("description").(string)
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "ram",
		"Action":          "AttachPolicyToRole",
		"Version":         "2015-05-01",
		"RoleName":        RoleName,
		"PolicyType":      PolicyType,
		"PolicyName":      PolicyName,
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
	request.ApiName = "AttachPolicyToRole"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf(" response of AttachPolicyToRole : %s", raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ram_policy_role_attachment", "AttachPolicyToRole", raw)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	if !bresponse.IsSuccess() || bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ram_policy_role_attachment", "AttachPolicyToRole", raw)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return WrapError(err)
	}
	//d.Set("policy_type", response.Policy.PolicyType)
	d.SetId(RoleName)
	return resourceApsaraStackRamRoleRead(d, meta)
}
func resourceApsaraStackRamPolicyRoleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*connectivity.ApsaraStackClient)
	//ascmService := AscmService{client}
	//policy, err := ascmService.DescribeRamRole(d.Id(), d.Get("policy_type").(string))
	//if err != nil {
	//	return WrapError(err)
	//}

	return nil
}
func resourceApsaraStackRamPolicyRoleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	RoleName := d.Get("role_name").(string)

	PolicyType := d.Get("policy_type").(string)
	PolicyName := d.Get("policy_name").(string)

	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "ram",
		"Action":          "DetachPolicyFromRole",
		"Version":         "2015-05-01",
		"RoleName":        RoleName,
		"PolicyType":      PolicyType,
		"PolicyName":      PolicyName,
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
	request.ApiName = "DetachPolicyFromRole"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf(" response of DetachPolicyFromRole : %s", raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ram_policy_role_attachment", "DetachPolicyFromRole", raw)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	if !bresponse.IsSuccess() || bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ram_policy_role_attachment", "DetachPolicyFromRole", raw)
	}

	return nil
}

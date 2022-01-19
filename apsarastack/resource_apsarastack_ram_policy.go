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
	"time"
)

func resourceApsaraStackRamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackRamPolicyCreate,
		Read:   resourceApsaraStackRamPolicyRead,
		Delete: resourceApsaraStackRamPolicyDelete,
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
			"policy_document": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"ram_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceApsaraStackRamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response *CreatePolicyResponse
	name := d.Get("name").(string)

	policyDoc := d.Get("policy_document").(string)
	description := d.Get("description").(string)
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "ram",
		"Action":          "CreatePolicy",
		"Version":         "2015-05-01",
		"PolicyName":      name,
		"Description":     description,
		"PolicyDocument":  policyDoc,
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
	request.ApiName = "CreatePolicy"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf(" response of CreatePolicy : %s", raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ram_policy", "CreatePolicy", raw)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	if !bresponse.IsSuccess() || bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ram_policy", "CreatePolicy", raw)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return WrapError(err)
	}
	d.Set("policy_type", response.Policy.PolicyType)
	d.SetId(name)
	return resourceApsaraStackRamPolicyRead(d, meta)
}
func resourceApsaraStackRamPolicyRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}
func resourceApsaraStackRamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	name := d.Get("name").(string)

	policyType := d.Get("policy_type").(string)
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "ram",
		"Action":          "DeletePolicy",
		"Version":         "2015-05-01",
		"PolicyName":      name,
		"PolicyType":      policyType,
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
	request.ApiName = "DeletePolicy"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf(" response of DeletePolicy : %s", raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ram_policy", "DeletePolicy", raw)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	if !bresponse.IsSuccess() || bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ram_policy", "DeletePolicy", raw)
	}

	return nil
}

type RamPolicy struct {
	Policy struct {
		PolicyType      string    `json:"PolicyType"`
		UpdateDate      time.Time `json:"UpdateDate"`
		Description     string    `json:"Description"`
		AttachmentCount int       `json:"AttachmentCount"`
		PolicyName      string    `json:"PolicyName"`
		DefaultVersion  string    `json:"DefaultVersion"`
		CreateDate      time.Time `json:"CreateDate"`
	} `json:"Policy"`
	ServerRole           string `json:"serverRole"`
	EagleEyeTraceID      string `json:"eagleEyeTraceId"`
	AsapiSuccess         bool   `json:"asapiSuccess"`
	AsapiRequestID       string `json:"asapiRequestId"`
	RequestID            string `json:"RequestId"`
	DefaultPolicyVersion struct {
		VersionID        string    `json:"VersionId"`
		IsDefaultVersion bool      `json:"IsDefaultVersion"`
		PolicyDocument   string    `json:"PolicyDocument"`
		CreateDate       time.Time `json:"CreateDate"`
	} `json:"DefaultPolicyVersion"`
	Domain string `json:"domain"`
	API    string `json:"api"`
}
type CreatePolicyResponse struct {
	Policy struct {
		PolicyType     string    `json:"PolicyType"`
		Description    string    `json:"Description"`
		PolicyName     string    `json:"PolicyName"`
		DefaultVersion string    `json:"DefaultVersion"`
		CreateDate     time.Time `json:"CreateDate"`
	} `json:"Policy"`
	ServerRole      string `json:"serverRole"`
	EagleEyeTraceID string `json:"eagleEyeTraceId"`
	AsapiSuccess    bool   `json:"asapiSuccess"`
	AsapiRequestID  string `json:"asapiRequestId"`
	RequestID       string `json:"RequestId"`
	Domain          string `json:"domain"`
	API             string `json:"api"`
}

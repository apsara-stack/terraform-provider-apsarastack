package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity/ascm"
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
				//ForceNew:     true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				//ForceNew:     true,
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
			//"ip_ranges":{
			//	Type: schema.TypeList,
			//	Optional: true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"ip_range": {
			//				Type:         schema.TypeString,
			//				Optional:     true,
			//			},
			//			"logon_policy_id": {
			//				Type:         schema.TypeString,
			//				Optional:     true,
			//			},
			//			"protocol": {
			//				Type:         schema.TypeString,
			//				Optional:     true,
			//				ValidateFunc: validation.StringInSlice([]string{"IPV4","IPV6"}, false),
			//			},
			//		},
			//	},
			//},
			//"time_ranges":{
			//	Type: schema.TypeList,
			//	Optional: true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"end_time": {
			//				Type:     schema.TypeString,
			//				Optional: true,
			//			},
			//			"login_policy_id": {
			//				Type:     schema.TypeInt,
			//				Optional: true,
			//			},
			//			"start_time": {
			//				Type:     schema.TypeString,
			//				Optional: true,
			//			},
			//		},
			//	},
			//},
		},
	}
}
func resourceApsaraStackLogInPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ascm.Client
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
		//iprangereq:=IPRanges{
		//	IPRange: d.Get("ip_range").(string),
		//
		//}
		//iprange:= d.Get("ip_range").(string)
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
			//"protocol":"IPV4",
			//"ipRange":"2.2.2.5/24",
			//"ipRange",iprangereq.IPRange,

		}
		request.Domain = client.Domain
		request.Method = "POST"        // Set request method
		request.Product = "ascm"       // Specify product
		request.Version = "2019-05-10" // Specify product version
		request.ServiceCode = "ascm"
		request.Scheme = "http" // Set request scheme. Default: http
		request.ApiName = "AddLoginPolicy"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		//	var err error
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_resource_group", "CreateResourceGroup", raw)
		}
		addDebug("CreateResourceGroup", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_resource_group", "CreateResourceGroup", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateResourceGroup", raw, requestInfo, bresponse.GetHttpContentString())
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		object, err = ascmService.ListLoginPolicies(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(object.Data) != 0 {
			return nil
		}
		return resource.RetryableError(Error("New Resource Group has been created successfully."))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_resource_group", "Failed to create resource set", ApsaraStackSdkGoERROR)
	}

	//err = nil
	//raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
	//		return ecsClient.ProcessCommonRequest(request)
	//	})
	//if err != nil {
	//	return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm", "CreateKubernetesCluster", raw)
	//}
	//
	////if debugOn() {
	////	requestMap := make(map[string]interface{})
	////	requestMap["RegionId"] = common.Region(client.RegionId)
	////	requestMap["Params"] = request.GetQueryParams()
	////	addDebug("CreateLogInPolicy", raw, request, requestMap)
	////}
	//
	//resp := responses.BaseResponse{}
	//ok := json.Unmarshal(resp.GetHttpContentBytes(), &resp)
	//if ok != nil {
	//	return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_login_policy", "ParseLoginPolicyResponse", raw)
	//}
	//}

	//d.SetId(d.Id())
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
	//
	//object, err := ascmService.ListLoginPolicies(name)
	//if err != nil {
	//	return WrapError(err)
	//}
	policyId := fmt.Sprint(d.Get("policy_id").(int))

	//id:=strconv.Itoa(object.Data[0].ID)

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
		//"protocol":"IPV4",
		//"ipRange":"2.2.2.5/24",
		//"ipRange",iprangereq.IPRange,

	}
	request.Domain = client.Domain
	request.Method = "POST"        // Set request method
	request.Product = "ascm"       // Specify product
	request.Version = "2019-05-10" // Specify product version
	request.ServiceCode = "ascm"
	request.Scheme = "http" // Set request scheme. Default: http
	request.ApiName = "ModifyLoginPolicy"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	//	err = nil
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm", "LoginPolicyUpdateRequestFailed", raw)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	if !bresponse.IsSuccess() {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm", "LoginPolicyUpdateFailed", raw)
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
	//name:=d.Get("name").(string)
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
	//	d.Set("id", object.Data[0].ID)
	//d.Set("cuser_id",object.CuserID)
	//d.Set("muser_id",object.MuserID)
	//d.Set("organizationVisibility",object.OrganizationVisibility)
	//d.Set("ownerOrganizationId",object.OwnerOrganizationID)
	d.Set("rule", object.Data[0].Rule)
	return nil
}
func resourceApsaraStackLogInPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ascm.Client

	name := d.Get("name").(string)

	check, err := ascmService.ListLoginPolicies(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsResourceGroupExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsLoginPolicyExist", check, requestInfo, map[string]string{"loginpolicyName": d.Id()})

	//object, err := ascmService.ListLoginPolicies(name)
	//if err != nil {
	//	return WrapError(err)
	//}
	//id:=strconv.Itoa(object.Data[0].ID)
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
			//"id":id,
			"Name": name,
			//"protocol":"IPV4",
			//"ipRange":"2.2.2.5/24",
			//"ipRange",iprangereq.IPRange,

		}
		request.Domain = client.Domain
		request.Method = "POST"        // Set request method
		request.Product = "ascm"       // Specify product
		request.Version = "2019-05-10" // Specify product version
		request.ServiceCode = "ascm"
		request.Scheme = "http" // Set request scheme. Default: http
		request.ApiName = "RemoveLoginPolicyByName"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		//err = nil
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
		//if check.Data[0].Name != "" {
		//	return resource.RetryableError(Error("Trying to delete LoginPolicy %#v successfully.", d.Id()))
		//}
		return nil
	})
	//if err != nil {
	//	return err
	//}
	//check, err := ascmService.ListLoginPolicies(d.Id())
	//
	//if err != nil {
	//	return err
	//}
	//if check.Data[0].Name != "" {
	//	return nil
	//}

	return nil
}

type IPRanges []struct {
	IPRange       string `json:"ipRange"`
	LoginPolicyID int    `json:"loginPolicyId"`
	Protocol      string `json:"protocol"`
}

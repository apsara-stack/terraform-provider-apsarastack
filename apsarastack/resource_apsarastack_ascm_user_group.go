package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strconv"
	"strings"
	"time"
)

func resourceApsaraStackAscmUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmUserGroupCreate,
		Read:   resourceApsaraStackAscmUserGroupRead,
		Update: resourceApsaraStackAscmUserGroupUpdate,
		Delete: resourceApsaraStackAscmUserGroupDelete,
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"role_in_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceApsaraStackAscmUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*connectivity.ApsaraStackClient)
	//var requestInfo *ecs.Client
	//
	//groupName := d.Get("group_name").(string)
	//organizationId := d.Get("organization_id").(string)
	//var loginNamesList []string
	//if v, ok := d.GetOk("role_in_ids"); ok {
	//	loginNames := expandStringList(v.(*schema.Set).List())
	//	for _, loginName := range loginNames {
	//		loginNamesList = append(loginNamesList, loginName)
	//	}
	//}
	//request := requests.NewCommonRequest()
	//if client.Config.Insecure {
	//	request.SetHTTPSInsecure(client.Config.Insecure)
	//}
	//request.Headers["x-ascm-product-name"] = "ascm"
	//request.Headers["x-ascm-product-version"] = "2019-05-10"
	//QueryParams := map[string]interface{}{
	//	"groupName":      groupName,
	//	"organizationId": organizationId,
	//	"roleIdList":     loginNamesList,
	//}
	//request.Method = "POST"
	//request.Product = "Ascm"
	//request.Version = "2019-05-10"
	//request.ServiceCode = "ascm"
	//request.Domain = "ascm.inter.env48.shuguang.com"
	//requeststring, err := json.Marshal(QueryParams)
	//if strings.ToLower(client.Config.Protocol) == "https" {
	//	request.Scheme = "https"
	//} else {
	//	request.Scheme = "http"
	//}
	//request.Headers["Content-Type"] = requests.Json
	//request.SetContent(requeststring)
	//request.PathPattern = "/ascm/auth/user/createUserGroup"
	//request.ApiName = "CreateUserGroup"
	//request.RegionId = client.RegionId
	//request.Headers["RegionId"] = client.RegionId
	//
	//raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
	//	return ecsClient.ProcessCommonRequest(request)
	//})
	client1 := meta.(*connectivity.ApsaraStackClient)
	groupName := d.Get("group_name").(string)
	organizationId := d.Get("organization_id").(string)
	var loginNamesList []string
	if v, ok := d.GetOk("role_in_ids"); ok {
		loginNames := expandStringList(v.(*schema.Set).List())
		for _, loginName := range loginNames {
			loginNamesList = append(loginNamesList, loginName)
		}
	}
	QueryParams := map[string]interface{}{
		"groupName":      groupName,
		"organizationId": organizationId,
		"roleIdList":     loginNamesList,
	}

	/*设置请求身份验证*/
	credential := credentials.NewStsTokenCredential(
		client1.AccessKey, // 请替换为您实际的AccessKey ID
		client1.SecretKey, // 请替换为您实际的AccessKey Secret
		"",                // 请替换为您实际的Security Token(非STS调用时为"")
	)
	/*创建请求连接*/
	client, _ := sdk.NewClientWithOptions(client1.RegionId, sdk.NewConfig(), credential)
	/*设置是否忽略证书*/
	client.SetHTTPSInsecure(true)
	/*(可选)设置创建连接超时时间*/
	client.SetConnectTimeout(1 * time.Second)
	/*(可选)设置读取超时时间*/
	//client.SetReadTimeout(10 * time.Second)
	/*（可选）请根据实际情况判断是否设置代理，设置方法如下：*/
	//client.SetHttpProxy("http://" + client1.Config.Proxy)
	//client.SetHttpsProxy("https://" + client1.Config.Proxy)

	/*构造请求对象*/
	request := requests.NewCommonRequest()
	request.Product = "ascm"
	request.ServiceCode = "ascm"
	request.Version = "2019-05-10"
	request.ApiName = "CreateUserGroup"
	request.PathPattern = "/ascm/auth/user/createUserGroup"
	endpoint := client1.Config.ASCMOpenAPIEndpoint
	request.Domain = endpoint
	request.Method = "POST"
	/*设置请求协议,默认http*/
	//request.Scheme = "https" // https | http
	if strings.ToLower(client1.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.SetContentType(requests.Json)
	requeststring, err := json.Marshal(QueryParams)
	//body := `{"groupName": "golangUserGroup","organizationId": 37, "description": "Golang调用示例", "roleIdList":["2","6"]}`
	request.Content = requeststring
	//request.Content = []byte(body)
	raw, err := client.ProcessCommonRequest(request)
	log.Printf("response of raw CreateUserGroup is : %s", raw)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user_group", "CreateUserGroup", raw)
	}

	addDebug("CreateUserGroup", raw, request)

	if raw.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user_group", "CreateUserGroup", ApsaraStackSdkGoERROR)
	}
	addDebug("CreateUserGroup", raw, raw.GetHttpContentString())
	//bresponse, _ := raw.(*responses.CommonResponse)
	//
	//if bresponse.GetHttpStatus() != 200 {
	//	return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user_group", "CreateUserGroup", ApsaraStackSdkGoERROR)
	//}
	//addDebug("CreateUserGroup", raw, bresponse.GetHttpContentString())

	d.SetId(groupName)

	return resourceApsaraStackAscmUserGroupUpdate(d, meta)
}

func resourceApsaraStackAscmUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	return resourceApsaraStackAscmUserGroupRead(d, meta)
}

func resourceApsaraStackAscmUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUserGroup(d.Id())
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

	d.Set("user_group_id", strconv.Itoa(object.Data[0].Id))
	d.Set("group_name", object.Data[0].GroupName)
	d.Set("organization_id", strconv.Itoa(object.Data[0].OrganizationId))

	var roleIds []string
	for _, role := range object.Data[0].Roles {
		roleIds = append(roleIds, strconv.Itoa(role.Id))
	}
	d.Set("role_ids", roleIds)

	return nil
}

func resourceApsaraStackAscmUserGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client

	check, err := ascmService.DescribeAscmUserGroup(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsUserGroupExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsUserGroupExist", check, requestInfo, map[string]string{"groupName": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "DeleteUserGroup",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"userGroupId":     strconv.Itoa(check.Data[0].Id),
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
		request.ApiName = "DeleteUserGroup"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err = ascmService.DescribeAscmUserGroup(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteUserGroup", raw, request)
		return nil
	})
	return nil
}

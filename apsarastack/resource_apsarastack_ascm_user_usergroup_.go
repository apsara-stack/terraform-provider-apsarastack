package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
	"time"
	//"encoding/json"
)

func resourceApsaraStackAscmUserGroupUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmUserGroupUserCreate,
		Read:   resourceApsaraStackAscmUserGroupUserRead,
		Delete: resourceApsaraStackAscmUserGroupUserDelete,
		Schema: map[string]*schema.Schema{
			"user_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"login_names": {
				Type: schema.TypeSet,
				//Computed: true,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			//"login_name": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	ForceNew: true,
			//},
		},
	}
}

func resourceApsaraStackAscmUserGroupUserCreate(d *schema.ResourceData, meta interface{}) error {
	client1 := meta.(*connectivity.ApsaraStackClient)
	//var requestInfo *ecs.Client
	//
	//userGroupId := d.Get("user_group_id").(string)
	//var loginNamesList []string
	//
	//if v, ok := d.GetOk("login_names"); ok {
	//	loginNames := expandStringList(v.(*schema.Set).List())
	//
	//	for _, loginName := range loginNames {
	//		loginNamesList = append(loginNamesList, loginName)
	//	}
	//}
	//
	//request := requests.NewCommonRequest()
	//if client.Config.Insecure {
	//	request.SetHTTPSInsecure(client.Config.Insecure)
	//}
	//
	//request.Headers["x-ascm-product-name"] = "ascm"
	//request.Headers["x-ascm-product-version"] = "2019-05-10"
	//
	//QueryParams := map[string]interface{}{
	//	"userGroupId":      userGroupId,
	//	"LoginNameList": loginNamesList,
	//}
	//
	//request.Method = "POST"
	//request.Product = "Ascm"
	//request.Version = "2019-05-10"
	//request.ServiceCode = "ascm"
	//request.Domain = client.Domain
	//requeststring,err := json.Marshal(QueryParams)
	//
	//if strings.ToLower(client.Config.Protocol) == "https" {
	//	request.Scheme = "https"
	//} else {
	//	request.Scheme = "http"
	//}
	//request.Headers["Content-Type"] = requests.Json
	//request.SetContent(requeststring)
	//request.PathPattern = "/roa/ascm/auth/user/addUsersToUserGroup"
	//request.ApiName = "AddUsersToUserGroup"
	//request.RegionId = client.RegionId
	//request.Headers["RegionId"] = client.RegionId
	//raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
	//	return ecsClient.ProcessCommonRequest(request)
	//})

	userGroupId := d.Get("user_group_id").(string)
	var loginNamesList []string

	if v, ok := d.GetOk("login_names"); ok {
		loginNames := expandStringList(v.(*schema.Set).List())

		for _, loginName := range loginNames {
			loginNamesList = append(loginNamesList, loginName)
		}
	}
	QueryParams := map[string]interface{}{
		"userGroupId":   userGroupId,
		"loginNameList": loginNamesList,
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
	//client.SetConnectTimeout(1 * time.Second)
	/*(可选)设置读取超时时间*/
	//client.SetReadTimeout(10 * time.Second)
	/*（可选）请根据实际情况判断是否设置代理，设置方法如下：*/

	/*构造请求对象*/
	request := requests.NewCommonRequest()
	request.Product = "ascm"
	request.ServiceCode = "ascm"
	request.Version = "2019-05-10"
	request.ApiName = "AddUsersToUserGroup"
	request.PathPattern = "/ascm/auth/user/addUsersToUserGroup"
	request.Domain = client1.Config.ASCMOpenAPIEndpoint
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
	request.Content = requeststring
	response, err := client.ProcessCommonRequest(request)
	log.Printf("response of raw AddUsersToUserGroup is : %s", response)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_usergroup_user", "AddUsersToUserGroup", response)
	}

	addDebug("AddUsersToUserGroup", response, request)

	if response.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_usergroup_user", "AddUsersToUserGroup", ApsaraStackSdkGoERROR)
	}
	addDebug("AddUsersToUserGroup", response, response.GetHttpContentString())

	d.SetId(userGroupId)

	return resourceApsaraStackAscmUserGroupUserRead(d, meta)
}

func resourceApsaraStackAscmUserGroupUserRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUsergroupUser(d.Id())
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

	var loginNames []string
	for _, data := range object.Data {
		loginNames = append(loginNames, data.LoginName)
	}

	d.Set("login_names", loginNames)

	return nil
}

func resourceApsaraStackAscmUserGroupUserDelete(d *schema.ResourceData, meta interface{}) error {
	client1 := meta.(*connectivity.ApsaraStackClient)

	//var login_names []string
	//userGroupId := d.Get("user_group_id").(string)
	//if v, ok := d.GetOk("login_names"); ok {
	//	login_names = expandStringList(v.(*schema.Set).List())
	//}
	//
	//request := requests.NewCommonRequest()
	//if client.Config.Insecure {
	//	request.SetHTTPSInsecure(client.Config.Insecure)
	//}
	//
	//request.Headers["x-ascm-product-name"] = "ascm"
	//request.Headers["x-ascm-product-version"] = "2019-05-10"
	//
	//QueryParams := map[string]interface{}{
	//	"userGroupId":   userGroupId,
	//	"LoginNameList": login_names,
	//}
	//
	//request.Method = "POST"
	//request.Product = "Ascm"
	//request.Version = "2019-05-10"
	//request.ServiceCode = "ascm"
	//request.Domain = client.Domain
	//requeststring, err := json.Marshal(QueryParams)
	//
	//if strings.ToLower(client.Config.Protocol) == "https" {
	//	request.Scheme = "https"
	//} else {
	//	request.Scheme = "http"
	//}
	//request.Headers["Content-Type"] = requests.Json
	//request.SetContent(requeststring)
	//request.PathPattern = "/roa/ascm/auth/user/RemoveUsersFromUserGroup"
	//request.ApiName = "RemoveUsersFromUserGroup"
	//request.RegionId = client.RegionId
	//request.Headers["RegionId"] = client.RegionId
	//raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
	//	return ecsClient.ProcessCommonRequest(request)
	//})
	var login_names []string
	userGroupId := d.Get("user_group_id").(string)
	if v, ok := d.GetOk("login_names"); ok {
		login_names = expandStringList(v.(*schema.Set).List())
	}
	QueryParams := map[string]interface{}{
		"userGroupId":   userGroupId,
		"loginNameList": login_names,
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
	request.ApiName = "RemoveUsersFromUserGroup"
	request.PathPattern = "/ascm/auth/user/removeUsersFromUserGroup"
	request.Domain = client1.Config.ASCMOpenAPIEndpoint
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
	request.Content = requeststring
	raw, err := client.ProcessCommonRequest(request)
	log.Printf("response of raw RemoveUsersFromUserGroup is : %s", raw)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_usergroup_user", "RemoveUsersFromUserGroup", raw)
	}

	return nil
}

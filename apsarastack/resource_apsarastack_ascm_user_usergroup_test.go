package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/go-yaml/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestAccApsaraStackAscm_UserGroup_User_Basic(t *testing.T) {
	var v *User
	resourceId := "apsarastack_ascm_usergroup_user.default"
	ra := resourceAttrInit(resourceId, testAccCheckUserGroupUserBinding)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-ascmusergroup%v", rand)
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		CheckDestroy: testAccCheckAscmUserGroupUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAscmUserGroupUserRoleBinding, name, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAscmUserGroupUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_ascm_usergroup_user" || rs.Type != "apsarastack_ascm_usergroup_user" {
			continue
		}
		ascm, err := ascmService.DescribeAscmUsergroupUser(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if ascm.Message != "" {
			return WrapError(Error("user  still exist"))
		}
	}

	return nil
}

const testAccCheckAscmUserGroupUserRoleBinding = `
resource "apsarastack_ascm_organization" "default" {
 name = "Test_binder"
 parent_id = "1"
}

resource "apsarastack_ascm_user_group" "default" {
 group_name =      "%s"
 organization_id = apsarastack_ascm_organization.default.org_id
}

resource "apsarastack_ascm_user" "default" {
 cellphone_number = "13900000000"
 email = "test@gmail.com"
 display_name = "C2C-DELTA"
 organization_id = apsarastack_ascm_organization.default.org_id
 mobile_nation_code = "91"
 login_name = "User_Role_Test%d"
 login_policy_id = 1
}

resource "apsarastack_ascm_usergroup_user" "default" {
  //login_name = apsarastack_ascm_user.default.login_name
  login_names = ["User_Role_Test6304175127373178963", "User_Role_Test7233024715252325400"]
  //login_names = ["[\"User_Role_Test929636066677054911\"]"]
  user_group_id = apsarastack_ascm_user_group.default.user_group_id
}
`

var testAccCheckUserGroupUserBinding = map[string]string{
	"user_group_id": CHECKSET,
	//"login_name": CHECKSET,
}

func TestCreateUserGroup(t *testing.T) {
	config := GetCloudConfig("48E")
	/*设置请求身份验证*/
	credential := credentials.NewStsTokenCredential(
		"51wQKxCJ2vZz2WqJ",               // 请替换为您实际的AccessKey ID
		"2WJyoVrGlCHlqC7coCgqb9y6TqMNkI", // 请替换为您实际的AccessKey Secret
		"",                               // 请替换为您实际的Security Token(非STS调用时为"")
	)
	/*创建请求连接*/
	client, _ := sdk.NewClientWithOptions("cn-wulan-env48-d01", sdk.NewConfig(), credential)
	/*设置是否忽略证书*/
	client.SetHTTPSInsecure(true)
	/*(可选)设置创建连接超时时间*/
	client.SetConnectTimeout(1 * time.Second)
	/*(可选)设置读取超时时间*/
	//client.SetReadTimeout(10 * time.Second)
	/*（可选）请根据实际情况判断是否设置代理，设置方法如下：*/
	if config.ProxyEnabled {
		client.SetHttpProxy("http://" + "HTTP://100.64.64.132:50646")
		client.SetHttpsProxy("https://" + "HTTP://100.64.64.132:50646")
	}
	/*构造请求对象*/
	request := requests.NewCommonRequest()
	request.Product = "ascm"
	request.ServiceCode = "ascm"
	request.Version = "2019-05-10"
	request.ApiName = "CreateUserGroup"
	request.PathPattern = "/ascm/auth/user/createUserGroup"
	request.Domain = "ascm.inter.env48.shuguang.com"
	request.Method = "POST"
	/*设置请求协议,默认http*/
	//request.Scheme = "https" // https | http
	request.SetContentType(requests.Json)
	body := `{"groupName": "golangUserGroup1","organizationId": 37, "description": "Golang调用示例", "roleIdList":["2","6"]}`
	request.Content = []byte(body)

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		if serverError, ok := err.(*errors.ServerError); ok {
			// 获取错误码
			fmt.Println(serverError.ErrorCode())
			fmt.Println(serverError.RequestId())
			// 获取错误描述
			fmt.Println(serverError.Message())
			// 获取原始http应答
			fmt.Println(response.GetOriginHttpResponse())
			return
		} else if clientError, ok := err.(*errors.ClientError); ok {
			// 获取错误码
			fmt.Println(clientError.ErrorCode())
			// 获取错误描述
			fmt.Println(clientError.Message())
			// 获取原始错误(可能为nil)
			fmt.Println(clientError.OriginError())
			return
		} else {
			panic(err)
		}
	}
	fmt.Print(response.GetHttpContentString())
}

type CloudConfig struct {
	Region         string `yaml:"region"`
	AK             string `yaml:"accessKey"`
	SK             string `yaml:"accessSecret"`
	StsToken       string `yaml:"stsToken"`
	InternetDomain string `yaml:"internetDomain"`
	IntranetDomain string `yaml:"intranetDomain"`
	ProxyEnabled   bool   `yaml:"proxyEnabled"`
	HttpProxy      string `yaml:"httpProxy"`
}

func GetCloudConfig(env string) CloudConfig {
	dir, _ := os.Getwd()
	index := strings.Index(dir, "/openapi-golang-samples")
	if index != -1 {
		dir = dir[0:index]
	}
	yamlfile, _ := ioutil.ReadFile(dir + "/config.yaml")
	resultMap := make(map[string]CloudConfig)
	yaml.Unmarshal(yamlfile, &resultMap)
	return resultMap[env]
}

func GetDomainByRule(rule string, config CloudConfig) string {
	regionPatterm, _ := regexp.Compile("\\$\\{global:region}")
	interDomainPatterm, _ := regexp.Compile("\\$\\{global:internet-domain}")
	intraDomainPatterm, _ := regexp.Compile("\\$\\{global:intranet-domain}")
	rule = regionPatterm.ReplaceAllString(rule, config.Region)
	rule = interDomainPatterm.ReplaceAllString(rule, config.InternetDomain)
	return intraDomainPatterm.ReplaceAllString(rule, config.IntranetDomain)
}

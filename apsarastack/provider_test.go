package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"log"
	"os"
	"testing"
	"time"

	"strings"

	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var defaultRegionToTest = os.Getenv("APSARASTACK_REGION")

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"apsarastack": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("APSARASTACK_ACCESS_KEY"); v == "" {
		t.Fatal("APSARASTACK_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("APSARASTACK_SECRET_KEY"); v == "" {
		t.Fatal("APSARASTACK_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("APSARASTACK_REGION"); v == "" {
		t.Fatal("APSARASTACK_REGION must be set for acceptance tests")
	}
	if v := os.Getenv("APSARASTACK_INSECURE"); v == "" {
		t.Fatal("APSARASTACK_INSECURE must be set for acceptance tests")
	}
	if v := os.Getenv("APSARASTACK_PROXY"); v == "" {
		t.Fatal("APSARASTACK_PROXY must be set for acceptance tests")
	}
	if v := os.Getenv("APSARASTACK_DOMAIN"); v == "" {
		t.Fatal("APSARASTACK_DOMAIN must be set for acceptance tests")
	}
}

func testAccPreCheckWithAccountSiteType(t *testing.T, account AccountSite) {
	defaultAccount := string(DomesticSite)
	if v := strings.TrimSpace(os.Getenv("APSARASTACK_ACCOUNT_SITE")); v != "" {
		defaultAccount = v
	}
	if defaultAccount != string(account) {
		t.Skipf("Skipping unsupported account type %s-Site. It only supports %s-Site.", defaultAccount, account)
		t.Skipped()
	}
}

func testAccPreCheckWithRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	if v := os.Getenv("APSARASTACK_ACCESS_KEY"); v == "" {
		t.Fatal("APSARASTACK_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("APSARASTACK_SECRET_KEY"); v == "" {
		t.Fatal("APSARASTACK_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("APSARASTACK_REGION"); v == "" {
		log.Println("[INFO] Test: Using cn-beijing as test region")
	}
	region := os.Getenv("APSARASTACK_REGION")
	find := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
	}

	if (find && !supported) || (!find && supported) {
		if supported {
			t.Skipf("Skipping unsupported region %s. Supported regions: %s.", region, regions)
		} else {
			t.Skipf("Skipping unsupported region %s. Unsupported regions: %s.", region, regions)
		}
		t.Skipped()
	}
}

// Skip automatically the sweep testcases which does not support some known regions.
// If supported is true, the regions should a list of supporting the service regions.
// If supported is false, the regions should a list of unsupporting the service regions.
func testSweepPreCheckWithRegions(region string, supported bool, regions []connectivity.Region) bool {
	find := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
	}
	return (find && !supported) || (!find && supported)
}

func testAccCheckApsaraStackDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("data source ID not set")
		}
		return nil
	}
}

func testAccPreCheckWithMultipleAccount(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("APSARASTACK_ACCESS_KEY_2")); v == "" {
		t.Skipf("Skipping unsupported test with multiple account")
		t.Skipped()
	}
}

func testAccPreCheckOSSForImageImport(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("APSARASTACK_OSS_BUCKET_FOR_IMAGE")); v == "" {
		t.Skipf("Skipping tests without OSS_Bucket set.")
		t.Skipped()
	}
	if v := strings.TrimSpace(os.Getenv("APSARASTACK_OSS_OBJECT_FOR_IMAGE")); v == "" {
		t.Skipf("Skipping OSS_Object does not exist.")
		t.Skipped()
	}
}

func testAccPreCheckWithCmsContactGroupSetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("APSARASTACK_CMS_CONTACT_GROUP")); v == "" {
		t.Skipf("Skipping the test case with no cms contact group setting")
		t.Skipped()
	}
}

func testAccPreCheckWithSmartAccessGatewaySetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("SAG_INSTANCE_ID")); v == "" {
		t.Skipf("Skipping the test case with no sag instance id setting")
		t.Skipped()
	}
}

func testAccPreCheckWithSmartAccessGatewayAppSetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("SAG_APP_INSTANCE_ID")); v == "" {
		t.Skipf("Skipping the test case with no sag app instance id setting")
		t.Skipped()
	}
}

func testAccPreCheckWithTime(t *testing.T) {
	if time.Now().Day() != 1 {
		t.Skipf("Skipping the test case with not the 1st of every month")
		t.Skipped()
	}
}

func testAccPreCheckWithAlikafkaAclEnable(t *testing.T) {
	aclEnable := os.Getenv("APSARASTACK_ALIKAFKA_ACL_ENABLE")

	if aclEnable != "true" {
		t.Skipf("Skipping the test case because the acl is not enabled.")
		t.Skipped()
	}
}

func testAccPreCheckWithNoDefaultVpc(t *testing.T) {
	region := os.Getenv("APSARASTACK_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	client := rawClient.(*connectivity.ApsaraStackClient)
	request := vpc.CreateDescribeVpcsRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc"}
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.PageNumber = requests.NewInteger(1)
	request.IsDefault = requests.NewBoolean(true)

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeVpcs(request)
	})
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	response, _ := raw.(*vpc.DescribeVpcsResponse)

	if len(response.Vpcs.Vpc) < 1 {
		t.Skipf("Skipping the test case with there is no default vpc")
		t.Skipped()
	}
}

func testAccPreCheckWithNoDefaultVswitch(t *testing.T) {
	region := os.Getenv("REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	client := rawClient.(*connectivity.ApsaraStackClient)
	request := vpc.CreateDescribeVSwitchesRequest()
	request.RegionId = string(client.Region)
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc"}
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.PageNumber = requests.NewInteger(1)
	request.IsDefault = requests.NewBoolean(true)

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeVSwitches(request)
	})
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	response, _ := raw.(*vpc.DescribeVSwitchesResponse)

	if len(response.VSwitches.VSwitch) < 1 {
		t.Skipf("Skipping the test case with there is no default vswitche")
		t.Skipped()
	}
}

var providerCommon = `
provider "apsarastack" {
	assume_role {}
}
`

func TestAccApsaraStackProviderEcs(t *testing.T) {
	var v ecs.Instance

	resourceId := "apsarastack_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAcc%sEcsInstanceConfigVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return providerCommon + resourceInstanceVpcConfigDependence(name)
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":        "${data.apsarastack_images.default.images.0.id}",
					"security_groups": []string{"${apsarastack_security_group.default.0.id}"},
					"instance_type":   "${data.apsarastack_instance_types.default.instance_types.0.id}",

					"availability_zone":             "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}",
					"system_disk_category":          "cloud_efficiency",
					"instance_name":                 "${var.name}",
					"key_name":                      "${apsarastack_key_pair.default.key_name}",
					"spot_strategy":                 "NoSpot",
					"spot_price_limit":              "0",
					"security_enhancement_strategy": "Active",
					"user_data":                     "I_am_user_data",

					"vswitch_id": "${apsarastack_vswitch.default.id}",
					"role_name":  "${apsarastack_ram_role.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"key_name":      name,
						"role_name":     name,
					}),
				),
			},
		},
	})
}

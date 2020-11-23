package apsarastack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("apsarastack_cms_alarm_contact", &resource.Sweeper{
		Name: "apsarastack_cms_alarm_contact",
		F:    testSweepCmsAlarmContact,
	})
}

func testSweepCmsAlarmContact(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting ApsaraStack client.")
	}
	client := rawClient.(*connectivity.ApsaraStackClient)
	cmsService := CmsService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	request := cms.CreateDescribeContactListRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "cms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	raw, err := cmsService.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeContactList(request)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve Cms Alarm in service list: %s", err)
	}

	var response *cms.DescribeContactListResponse
	response, _ = raw.(*cms.DescribeContactListResponse)

	for _, v := range response.Contacts.Contact {
		name := v.Name
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping alarm contact: %s ", name)
			continue
		}
		log.Printf("[INFO] delete alarm contact: %s ", name)

		request := cms.CreateDeleteContactRequest()
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "cms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.ContactName = v.Name
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "cms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteContact(request)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to delete alarm contact (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccApsaraStackCmsAlarmContact_basic(t *testing.T) {
	var v cms.Contact
	resourceId := "apsarastack_cms_alarm_contact.default"
	ra := resourceAttrInit(resourceId, CmsAlarmContactMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeCmsAlarmContact")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCmsAlarmContactzhangsan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, CmsAlarmContactBasicdependence)
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
					"alarm_contact_name": "${var.name}",
					"describe":           "For-test",
					"channels_mail":      "hello.uuuu@aaa.com",
					"lifecycle": []map[string]interface{}{
						{
							"ignore_changes": []string{"channels_mail"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_contact_name": name,
						"describe":           "For-test",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var CmsAlarmContactMap = map[string]string{}

func CmsAlarmContactBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}

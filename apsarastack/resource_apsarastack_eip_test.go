package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"log"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("apsarastack_eip", &resource.Sweeper{
		Name: "apsarastack_eip",
		F:    testSweepEips,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"apsarastack_instance",
			"apsarastack_slb",
			"apsarastack_nat_gateway",
		},
	})
}

func testSweepEips(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting ApsaraStack client: %s", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var eips []vpc.EipAddress
	req := vpc.CreateDescribeEipAddressesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeEipAddresses(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving EIPs: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeEipAddressesResponse)
		if resp == nil || len(resp.EipAddresses.EipAddress) < 1 {
			break
		}
		eips = append(eips, resp.EipAddresses.EipAddress...)

		if len(resp.EipAddresses.EipAddress) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	for _, v := range eips {
		name := v.Name
		id := v.AllocationId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping EIP: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting EIP: %s (%s)", name, id)
		req := vpc.CreateReleaseEipAddressRequest()
		req.AllocationId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ReleaseEipAddress(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete EIP (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckEIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_eip" {
			continue
		}

		_, err := vpcService.DescribeEip(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func TestAccApsaraStackEipBasic_PayByBandwidth(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "apsarastack_eip.default"
	ra := resourceAttrInit(resourceId, testAccCheckEipCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEipConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCheckEipConfig_bandwidth(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAcceEipName%d", rand),
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAcceEipName%d_description", rand),
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcceEipName%d_all", rand),
						"description": fmt.Sprintf("tf-testAcceEipName%d_description_all", rand),
					}),
				),
			},
		},
	})

}

func TestAccApsaraStackEipBasic_PayByTraffic(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "apsarastack_eip.default"
	ra := resourceAttrInit(resourceId, testAccCheckEipCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEipConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccCheckEipConfig_bandwidth(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAcceEipName%d", rand),
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAcceEipName%d_description", rand),
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcceEipName%d_all", rand),
						"description": fmt.Sprintf("tf-testAcceEipName%d_description_all", rand),
					}),
				),
			},
		},
	})

}

func TestAccApsaraStackEipMulti(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "apsarastack_eip.default.9"
	ra := resourceAttrInit(resourceId, testAccCheckEipCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEipConfig_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckEipConfigBasic(rand int) string {
	return fmt.Sprintf(`
resource "apsarastack_eip" "default" {
	bandwidth = "5"
}
`)
}

func testAccCheckEipConfig_bandwidth(rand int) string {
	return fmt.Sprintf(`
resource "apsarastack_eip" "default" {
     bandwidth = "10"
}
`)
}

func testAccCheckEipConfig_name(rand int) string {
	return fmt.Sprintf(`
variable "name"{
	default = "tf-testAcceEipName%d"
}
resource "apsarastack_eip" "default" {
	bandwidth = "10"
	name = "${var.name}"
}
`, rand)
}

func testAccCheckEipConfig_description(rand int) string {
	return fmt.Sprintf(`
variable "name"{
	default = "tf-testAcceEipName%d"
}
resource "apsarastack_eip" "default" {
	bandwidth = "10"
	name = "${var.name}"
    description = "${var.name}_description"
}
`, rand)
}

func testAccCheckEipConfig_all(rand int) string {
	return fmt.Sprintf(`
variable "name"{
	default = "tf-testAcceEipName%d"
}
resource "apsarastack_eip" "default" {	
	bandwidth = "10"
	name = "${var.name}_all"
    description = "${var.name}_description_all"
}
`, rand)
}

func testAccCheckEipConfig_multi(rand int) string {
	return fmt.Sprintf(`
resource "apsarastack_eip" "default" {
    count = 10
	bandwidth = "5"
}
`)
}

var testAccCheckEipCheckMap = map[string]string{
	"name":        "",
	"description": "",
	"bandwidth":   "5",
	// read method does't return a value for the period attribute, so it is not tested
	"ip_address": CHECKSET,
	"status":     CHECKSET,
}

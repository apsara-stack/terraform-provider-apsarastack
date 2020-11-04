package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccApsaraStackEssLifecycleHookBasic(t *testing.T) {
	rand := acctest.RandIntRange(10, 99999)
	var v ess.LifecycleHook
	resourceId := "apsarastack_ess_lifecycle_hook.default"
	basicMap := map[string]string{
		"name":                  fmt.Sprintf("tf-testAccEssLifecycleHook-%d", rand),
		"lifecycle_transition":  "SCALE_OUT",
		"heartbeat_timeout":     "600",
		"notification_metadata": "helloworld",
		"default_result":        "CONTINUE",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEssLifecycleHookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssLifecycleHook(EcsInstanceCommonTestCase, rand),
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
				Config: testAccEssLifecycleHookUpdateLifecycleTransition(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_transition": "SCALE_IN",
					}),
				),
			},
			{
				Config: testAccEssLifecycleHookUpdateHeartbeatTimeout(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"heartbeat_timeout": "400",
					}),
				),
			},
			{
				Config: testAccEssLifecycleHookUpdateNotificationMetadata(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notification_metadata": "helloterraform",
					}),
				),
			},
			{
				Config: testAccEssLifecycleHookUpdateDefaultResult(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_result": "ABANDON",
					}),
				),
			},
			{
				Config: testAccEssLifecycleHookUpdateNotificationArn(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notification_arn": CHECKSET,
					}),
				),
			},
			{
				Config: testAccEssLifecycleHook(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})
}

func testAccCheckEssLifecycleHookDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_ess_lifecycle_hook" {
			continue
		}
		if _, err := essService.DescribeEssLifecycleHook(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("lifecycle hook %s still exists.", rs.Primary.ID)
	}
	return nil
}

func testAccEssLifecycleHook(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "apsarastack_vswitch" "default2" {
		  vpc_id = "${apsarastack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
		  name = "${var.name}"
	}
	
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${apsarastack_vswitch.default.id}","${apsarastack_vswitch.default2.id}"]
	}
	
	resource "apsarastack_ess_lifecycle_hook" "default"{
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_OUT"
		notification_metadata = "helloworld"
	}
	`, common, rand)
}
func testAccEssLifecycleHookUpdateLifecycleTransition(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "apsarastack_vswitch" "default2" {
		  vpc_id = "${apsarastack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
		  name = "${var.name}"
	}
	
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${apsarastack_vswitch.default.id}","${apsarastack_vswitch.default2.id}"]
	}
	
	resource "apsarastack_ess_lifecycle_hook" "default"{
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_IN"
		notification_metadata = "helloworld"
	}
	`, common, rand)
}
func testAccEssLifecycleHookUpdateHeartbeatTimeout(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "apsarastack_vswitch" "default2" {
		  vpc_id = "${apsarastack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
		  name = "${var.name}"
	}
	
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${apsarastack_vswitch.default.id}","${apsarastack_vswitch.default2.id}"]
	}
	
	resource "apsarastack_ess_lifecycle_hook" "default"{
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_IN"
		heartbeat_timeout = 400
		notification_metadata = "helloworld"
	}
	`, common, rand)
}
func testAccEssLifecycleHookUpdateNotificationMetadata(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "apsarastack_vswitch" "default2" {
		  vpc_id = "${apsarastack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
		  name = "${var.name}"
	}
	
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${apsarastack_vswitch.default.id}","${apsarastack_vswitch.default2.id}"]
	}
	
	resource "apsarastack_ess_lifecycle_hook" "default"{
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_IN"
		heartbeat_timeout = 400
		notification_metadata = "helloterraform"
	}
	`, common, rand)
}
func testAccEssLifecycleHookUpdateDefaultResult(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "apsarastack_vswitch" "default2" {
		  vpc_id = "${apsarastack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
		  name = "${var.name}"
	}
	
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${apsarastack_vswitch.default.id}","${apsarastack_vswitch.default2.id}"]
	}
	
	resource "apsarastack_ess_lifecycle_hook" "default"{
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_IN"
		heartbeat_timeout = 400
		notification_metadata = "helloterraform"
		default_result = "ABANDON"
	}
	`, common, rand)
}
func testAccEssLifecycleHookUpdateNotificationArn(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	data "apsarastack_zones" "default2" {
	
	}

	resource "apsarastack_vswitch" "default2" {
		  vpc_id = "${apsarastack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.apsarastack_zones.default2.zones.0.id}"
		  name = "${var.name}"
	}
	
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${apsarastack_vswitch.default.id}","${apsarastack_vswitch.default2.id}"]
	}
	
	resource "apsarastack_ess_lifecycle_hook" "default"{
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_IN"
		heartbeat_timeout = 400
		notification_metadata = "helloterraform"
		default_result = "ABANDON"
		notification_arn = "acs:ess"
	}
	`, common, rand)
}

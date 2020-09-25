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

func TestAccApsaraStackEssScalingRule_basic(t *testing.T) {
	var v ess.ScalingRule
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "apsarastack_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
		"adjustment_type":  "TotalCapacity",
		"adjustment_value": "1",
		"cooldown":         "0",
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

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingRuleConfig(EcsInstanceCommonTestCase, rand),
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
				Config: testAccEssScalingRuleUpdateAdjustmentType(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"adjustment_type": "PercentChangeInCapacity",
					}),
				),
			},
			{
				Config: testAccEssScalingRuleUpdateAdjustmentValue(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"adjustment_value": "2",
					}),
				),
			},
			{
				Config: testAccEssScalingRuleUpdateScalingRuleName(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_name": fmt.Sprintf("tf-testAccEssScalingRuleConfig-%d", rand),
					}),
				),
			},
			{
				Config: testAccEssScalingRuleUpdateCooldown(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cooldown": "200",
					}),
				),
			},
			{
				Config: testAccEssScalingRuleConfig(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})
}

func TestAccApsaraStackEssScalingRule_target_tracking_rule_basic(t *testing.T) {
	var v ess.ScalingRule
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "apsarastack_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
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
		CheckDestroy:  testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssTargetTrackingScalingRuleConfig(EcsInstanceCommonTestCase, rand),
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
				Config: testAccEssTargetTrackingScalingRuleConfig(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})
}

func TestAccApsaraStackEssScalingRule_step_rule_basic(t *testing.T) {
	var v ess.ScalingRule
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "apsarastack_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
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

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssStepScalingRuleConfig(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

func TestAccApsaraStackEssScalingRuleMulti(t *testing.T) {
	var v ess.ScalingRule
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "apsarastack_ess_scaling_rule.default.9"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
		"adjustment_type":  "TotalCapacity",
		"adjustment_value": "1",
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

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingRuleConfigMulti(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckEssScalingRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_ess_scaling_rule" {
			continue
		}
		_, err := essService.DescribeEssScalingRule(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Scaling rule %s still exists.", rs.Primary.ID)
	}

	return nil
}

func testAccEssScalingRuleConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	resource "apsarastack_ess_scaling_configuration" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		image_id = "${data.apsarastack_images.default.images.0.id}"
		instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
		security_group_id = "${apsarastack_security_group.default.id}"
		force_delete = "true"
	}
	resource "apsarastack_ess_scaling_rule" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 1
		cooldown = 0
	}
	`, common, rand)
}

func testAccEssScalingRuleUpdateAdjustmentType(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	resource "apsarastack_ess_scaling_configuration" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		image_id = "${data.apsarastack_images.default.images.0.id}"
		instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
		security_group_id = "${apsarastack_security_group.default.id}"
		force_delete = "true"
	}
	resource "apsarastack_ess_scaling_rule" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		adjustment_type = "PercentChangeInCapacity"
		adjustment_value = 1
		cooldown = 0
	}
	`, common, rand)
}

func testAccEssScalingRuleUpdateAdjustmentValue(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	resource "apsarastack_ess_scaling_configuration" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		image_id = "${data.apsarastack_images.default.images.0.id}"
		instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
		security_group_id = "${apsarastack_security_group.default.id}"
		force_delete = "true"
	}
	resource "apsarastack_ess_scaling_rule" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		adjustment_type = "PercentChangeInCapacity"
		adjustment_value = 2
		cooldown = 0
	}
	`, common, rand)
}

func testAccEssScalingRuleUpdateScalingRuleName(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	resource "apsarastack_ess_scaling_configuration" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		image_id = "${data.apsarastack_images.default.images.0.id}"
		instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
		security_group_id = "${apsarastack_security_group.default.id}"
		force_delete = "true"
	}
	resource "apsarastack_ess_scaling_rule" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		adjustment_type = "PercentChangeInCapacity"
		adjustment_value = 2
		scaling_rule_name = "${var.name}"
		cooldown = 0
	}
	`, common, rand)
}

func testAccEssScalingRuleUpdateCooldown(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	resource "apsarastack_ess_scaling_configuration" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		image_id = "${data.apsarastack_images.default.images.0.id}"
		instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
		security_group_id = "${apsarastack_security_group.default.id}"
		force_delete = "true"
	}
	resource "apsarastack_ess_scaling_rule" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		adjustment_type = "PercentChangeInCapacity"
		adjustment_value = 2
		scaling_rule_name = "${var.name}"
		cooldown = 200
	}
	`, common, rand)
}

func testAccEssScalingRuleConfigMulti(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	resource "apsarastack_ess_scaling_configuration" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		image_id = "${data.apsarastack_images.default.images.0.id}"
		instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
		security_group_id = "${apsarastack_security_group.default.id}"
		force_delete = "true"
	}
	resource "apsarastack_ess_scaling_rule" "default" {
		count = 10
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 1
	}
	`, common, rand)
}

func testAccEssTargetTrackingScalingRuleConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssTargetTrackingScalingRuleConfig-%d"
	}
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
	}
	resource "apsarastack_ess_scaling_rule" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
	}
	`, common, rand)
}

func testAccEssStepScalingRuleConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssStepScalingRuleConfig-%d"
	}
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
	}
	resource "apsarastack_ess_scaling_rule" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
	}
	`, common, rand)
}

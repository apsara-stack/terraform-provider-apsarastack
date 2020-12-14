package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackEssScheduledtasksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackEssScheduledTasksDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_ess_scheduled_tasks.default"),
					resource.TestCheckResourceAttr("data.apsarastack_ess_scheduled_tasks.default", "tasks.#", "1"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ess_scheduled_tasks.default", "tasks.1.id"),
					resource.TestCheckResourceAttrSet("data.apsarastack_ess_scheduled_tasks.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackEssScheduledTasksDataSource = `

variable "name" {
	default = "tf-testAccDataSourceScheduledtas"
}

data "apsarastack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.0.0.0/8"
}

resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "10.1.0.0/16"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "apsarastack_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${apsarastack_vswitch.default.id}"]
}

resource "apsarastack_ess_scaling_rule" "default" {
  scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
  adjustment_type  = "TotalCapacity"
  adjustment_value = 2
  cooldown         = 60
}

resource "apsarastack_ess_scheduled_task" "default" {
  scheduled_action    = "${apsarastack_ess_scaling_rule.default.ari}"
  launch_time         = "2020-12-07T12:50Z"
  scheduled_task_name = "${var.name}"
}

data "apsarastack_ess_scheduled_tasks" "default"{
  ids = ["${apsarastack_ess_scheduled_task.default.id}"]
}
`

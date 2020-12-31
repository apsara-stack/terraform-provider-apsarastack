package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackInstancesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_instances.default"),
					resource.TestCheckResourceAttr("data.apsarastack_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instances.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackInstancesDataSource = `
variable "name" {
  default = "Tf-EcsInstanceDataSource"
}

data "apsarastack_zones" "default" {
   available_disk_category = "cloud_ssd"
   available_resource_creation= "VSwitch"
}
resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_instance" "default" {
  image_id = "wincore_2004_x64_dtc_en-us_40G_alibase_20201015.raw"
  instance_type = "ecs.n4.xlarge"
  instance_name = "${var.name}"
  internet_max_bandwidth_out = "10"
  security_groups = "${apsarastack_security_group.default.*.id}"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  system_disk_category = "cloud_ssd"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}
data "apsarastack_instances" "default" {
  ids = ["${apsarastack_instance.default.id}"]
}
`

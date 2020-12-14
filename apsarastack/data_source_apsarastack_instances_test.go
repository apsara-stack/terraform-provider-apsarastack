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

data "apsarastack_instance_types" "default" {
  eni_amount        = 2
}
data "apsarastack_images" "default" {
  most_recent = true
  owners = "system"
}
resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_instance" "default" {
  image_id = "${data.apsarastack_images.default.images.0.id}"
  instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  internet_max_bandwidth_out = "10"
  security_groups = "${apsarastack_security_group.default.*.id}"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}
data "apsarastack_instances" "default" {
  ids = ["${apsarastack_instance.default.id}"]
}
`

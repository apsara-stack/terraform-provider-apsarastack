package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackSlbBackendServersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSlbBackendServersDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_slb_backend_servers.default"),
					resource.TestCheckResourceAttr("data.apsarastack_slb_backend_servers.default", "load_balancer_id.#", "0"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAccCheckApsaraStackSlbBackendServersDataSource = `
variable "name" {
	default = "tf-slbBackendServersdatasourcebasic"
}

data "apsarastack_zones" "default" {
	available_resource_creation = "VSwitch"
}
data "apsarastack_images" "default" {
  most_recent = true
  owners = "system"
}
data "apsarastack_instance_types" "default" {
	cpu_core_count = 2
	memory_size = 4
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
}

resource "apsarastack_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}

resource "apsarastack_instance" "default" {
  image_id = "${data.apsarastack_images.default.images.0.id}"

  instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  system_disk_category = "cloud_efficiency"

  security_groups = ["${apsarastack_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}

resource "apsarastack_slb_backend_server" "default" {
  load_balancer_id = "${apsarastack_slb.default.id}"

  backend_servers {
    server_id = "${apsarastack_instance.default.id}"
    weight     = 100
  }
}

data "apsarastack_slb_backend_servers" "default" {
 load_balancer_id = "${apsarastack_slb.default.id}"
}
`

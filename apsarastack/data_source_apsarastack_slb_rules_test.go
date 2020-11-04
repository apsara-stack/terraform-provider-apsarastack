package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackSlbRulesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSlbRulesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_slb_rules.default"),
					resource.TestCheckResourceAttr("data.apsarastack_slb_rules.default", "load_balancer_id.#", "0"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

const testAccCheckApsaraStackSlbRulesDataSourceConfig = `
variable "name" {
	default = "tf-testaccslbrulesdatasourcebasic"
}

data "apsarastack_images" "default" {
  most_recent = true
  owners = "system"
}
data "apsarastack_instance_types" "default" {
 	cpu_core_count = 1
	memory_size = 2
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
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

resource "apsarastack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}

resource "apsarastack_slb_listener" "default" {
  load_balancer_id = "${apsarastack_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie = "${var.name}"
  cookie_timeout = 86400
  health_check = "on"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
}

resource "apsarastack_instance" "default" {
  image_id = "${data.apsarastack_images.default.images.0.id}"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  system_disk_category = "cloud_efficiency"
  security_groups = ["${apsarastack_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}

resource "apsarastack_slb_server_group" "default" {
  load_balancer_id = "${apsarastack_slb.default.id}"
  servers {
      server_ids = ["${apsarastack_instance.default.id}"]
      port = 80
      weight = 100
    }
}

resource "apsarastack_slb_rule" "default" {
  load_balancer_id = "${apsarastack_slb.default.id}"
  frontend_port = "${apsarastack_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${apsarastack_slb_server_group.default.id}"
}

data "apsarastack_slb_rules" "default" {
load_balancer_id = "${apsarastack_slb.default.id}"
frontend_port = 80
}
`

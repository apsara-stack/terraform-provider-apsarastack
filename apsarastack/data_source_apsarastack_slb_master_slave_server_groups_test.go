package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackSlbMasterSlaveServerGroupsDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSlbMasterSlaveServerGroupsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_slb_master_slave_server_groups.default"),
					resource.TestCheckResourceAttr("data.apsarastack_slb_master_slave_server_groups.default", "load_balancer_id.#", "0"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAccCheckApsaraStackSlbMasterSlaveServerGroupsDataSourceBasic = `
variable "name" {
  default = "tf-testAccslbmasterslaveservergroupsdatasourcebasic"
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "apsarastack_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
   availability_zone = "cn-qingdao-env66-amtest66001-a"
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
  availability_zone = "cn-qingdao-env66-amtest66001-a"
  security_groups = ["${apsarastack_security_group.default.id}"]
  count                      = "2"
  instance_type              = "ecs.n4.large"
  image_id                   =  "m-9fx0253j413yavd1t8ba"
  instance_name              = "${var.name}"
  vswitch_id                 = apsarastack_vswitch.default.id
  internet_max_bandwidth_out = 10
}

resource "apsarastack_slb_master_slave_server_group" "default" {
  load_balancer_id = "${apsarastack_slb.default.id}"
  name = "${var.name}"
  servers {
      server_id = "${apsarastack_instance.default.0.id}"
      port = 80
      weight = 100
      server_type = "Master"
  }
  servers {
      server_id = "${apsarastack_instance.default.1.id}"
      port = 80
      weight = 100
      server_type = "Slave"
  }
}

data "apsarastack_slb_master_slave_server_groups" "default" {
  load_balancer_id = "${apsarastack_slb.default.id}"
}`

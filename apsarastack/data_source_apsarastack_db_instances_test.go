package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackDBInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackDBInstanceDataSourceConfig_mysql,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_db_instances.default"),
					resource.TestCheckResourceAttr("data.apsarastack_db_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.apsarastack_db_instances.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackDBInstanceDataSourceConfig_mysql = `
variable "name" {
  default = "tf-testAccDBInstanceConfig"
}
data "apsarastack_db_zones" "default" {
  multi = true
}

data "apsarastack_zones" "default" {
  available_resource_creation = "Rds"
}

resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "apsarastack_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "5.7"
  instance_type        = "mysql.x8.medium.2"
  instance_storage     = "50"
  instance_charge_type = "Postpaid"
  instance_name        = "${var.name}"
  monitoring_period    = "60"
  vswitch_id = "${apsarastack_vswitch.default.id}"
  zone_id =  data.apsarastack_db_zones.default.zones[0].id
}
data "apsarastack_db_instances" "default" {
  name_regex = "${apsarastack_db_instance.default.instance_name}"
  ids        = ["${apsarastack_db_instance.default.id}"]
  status     = "Running"
  tags       = {
    "type" = "database",
    "size" = "tiny"
  }
}
`

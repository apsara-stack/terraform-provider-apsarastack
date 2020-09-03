package apsarastack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccApsaraStackSlbMasterSlaveServerGroupsDataSource_basic(t *testing.T) {
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${apsarastack_slb.default.id}"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${apsarastack_slb.default.id}"`,
			"name_regex":       `"${apsarastack_slb_master_slave_server_group.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${apsarastack_slb.default.id}"`,
			"name_regex":       `"${apsarastack_slb_master_slave_server_group.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${apsarastack_slb.default.id}"`,
			"ids":              `["${apsarastack_slb_master_slave_server_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${apsarastack_slb.default.id}"`,
			"ids":              `["${apsarastack_slb_master_slave_server_group.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${apsarastack_slb.default.id}"`,
			"ids":              `["${apsarastack_slb_master_slave_server_group.default.id}"]`,
			"name_regex":       `"${apsarastack_slb_master_slave_server_group.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${apsarastack_slb.default.id}"`,
			"ids":              `["${apsarastack_slb_master_slave_server_group.default.id}_fake"]`,
			"name_regex":       `"${apsarastack_slb_master_slave_server_group.default.name}"`,
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":                       "1",
			"ids.#":                          "1",
			"names.#":                        "1",
			"groups.0.id":                    CHECKSET,
			"groups.0.name":                  "tf-testAccslbmasterslaveservergroupsdatasourcebasic",
			"groups.0.servers.#":             "2",
			"groups.0.servers.0.weight":      "100",
			"groups.0.servers.0.instance_id": CHECKSET,
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"ids.#":    "0",
			"names.#":  "0",
		}
	}

	var slbServerGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_slb_master_slave_server_groups.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbServerGroupsCheckInfo.dataSourceTestCheck(t, -1, allConf, basicConf, nameRegexConf, idsConf)
}

func testAccCheckApsaraStackSlbMasterSlaveServerGroupsDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccslbmasterslaveservergroupsdatasourcebasic"
}

data "apsarastack_zones" "default" {
  available_resource_creation = "VSwitch"
  available_disk_category     = "cloud_efficiency"
}

data "apsarastack_images" "default" {
  name_regex = "^ubuntu_18.*64"
  most_recent = true
  owners = "system"
}

data "apsarastack_instance_types" "default" {
  cpu_core_count = 2
  memory_size = 4
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
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
  specification = "slb.s2.small"
}

resource "apsarastack_slb_listener" "default" {
  load_balancer_id             = "${apsarastack_slb.default.id}"
  master_slave_server_group_id = "${apsarastack_slb_master_slave_server_group.default.id}"
  frontend_port                = "22"
  protocol                     = "tcp"
  bandwidth                    = "10"
  health_check_type            = "tcp"
  persistence_timeout          = 3600
  healthy_threshold            = 8
  unhealthy_threshold          = 8
  health_check_timeout         = 8
  health_check_interval        = 5
  health_check_http_code       = "http_2xx"
  health_check_connect_port    = 20
  health_check_uri             = "/console"
  established_timeout          = 600
}

resource "apsarastack_instance" "default" {
  image_id = "${data.apsarastack_images.default.images.0.id}"

  instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  internet_charge_type = "PayByTraffic"
  system_disk_category = "cloud_efficiency"

  security_groups = ["${apsarastack_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = "${apsarastack_vswitch.default.id}"

  count = "2"
  internet_max_bandwidth_out = "10"
  availability_zone          = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  instance_charge_type       = "PostPaid"
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
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}

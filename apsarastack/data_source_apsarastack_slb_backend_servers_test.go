package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccApsaraStackSlbBackendServersDataSource_basic(t *testing.T) {

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbBackendServersDataSourceConfig(map[string]string{
			"load_balancer_id": `"${apsarastack_slb_backend_server.default.load_balancer_id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbBackendServersDataSourceConfig(map[string]string{
			"load_balancer_id": `"${apsarastack_slb_backend_server.default.load_balancer_id}"`,
			"ids":              `["${apsarastack_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbBackendServersDataSourceConfig(map[string]string{
			"load_balancer_id": `"${apsarastack_slb_backend_server.default.load_balancer_id}"`,
			"ids":              `["${apsarastack_instance.default.id}_fake"]`,
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"backend_servers.#":        "1",
			"backend_servers.0.id":     CHECKSET,
			"backend_servers.0.weight": "100",
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"backend_servers.#": "0",
			"ids.#":             "0",
		}
	}

	var slbServerGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_slb_backend_servers.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbServerGroupsCheckInfo.dataSourceTestCheck(t, acctest.RandInt(), idsConf, allConf)
}

func testAccCheckApsaraStackSlbBackendServersDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccslbbackendserversdatasourcebasic"
}

data "apsarastack_zones" "default" {
	available_resource_creation = "VSwitch"
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
  internet_charge_type = "PayByTraffic"
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
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}

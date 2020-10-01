package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccApsaraStackRouteEntriesDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${apsarastack_route_entry.default.nexthop_id}"`,
			"route_table_id": `"${apsarastack_route_entry.default.route_table_id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${apsarastack_route_entry.default.nexthop_id}"`,
			"route_table_id": `"${apsarastack_route_entry.default.route_table_id}_fake"`,
		}),
	}

	typeConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${apsarastack_route_entry.default.nexthop_id}"`,
			"route_table_id": `"${apsarastack_route_entry.default.route_table_id}"`,
			"type":           `"Custom"`,
		}),
		fakeConfig: testAccCheckApsaraStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${apsarastack_route_entry.default.nexthop_id}"`,
			"route_table_id": `"${apsarastack_route_entry.default.route_table_id}"`,
			"type":           `"System"`,
		}),
	}

	cidrBlockConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"route_table_id": `"${apsarastack_route_entry.default.route_table_id}"`,
			"cidr_block":     `"${apsarastack_route_entry.default.destination_cidrblock}"`,
		}),
		fakeConfig: testAccCheckApsaraStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"route_table_id": `"${apsarastack_route_entry.default.route_table_id}"`,
			"cidr_block":     `"${apsarastack_route_entry.default.destination_cidrblock}_fake"`,
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${apsarastack_instance.default.id}"`,
			"route_table_id": `"${apsarastack_route_entry.default.route_table_id}"`,
			"type":           `"Custom"`,
			"cidr_block":     `"${apsarastack_route_entry.default.destination_cidrblock}"`,
		}),
		fakeConfig: testAccCheckApsaraStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${apsarastack_instance.default.id}"`,
			"route_table_id": `"${apsarastack_route_entry.default.route_table_id}"`,
			"type":           `"Custom"`,
			"cidr_block":     `"${apsarastack_route_entry.default.destination_cidrblock}_fake"`,
		}),
	}

	routeEntriesCheckInfo.dataSourceTestCheck(t, rand, instanceIdConf, typeConfig, cidrBlockConfig, allConfig)
}

func testAccCheckApsaraStackRouteEntriesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "apsarastack_zones" "default" {
  available_resource_creation= "VSwitch"
}
data "apsarastack_instance_types" "default" {
   availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
   cpu_core_count = 1
   memory_size = 2
}
data "apsarastack_images" "default" {
   name_regex = "^ubuntu_18.*64"
   most_recent = true
   owners = "system"
}
variable "name" {
   default = "tf-testAcc-for-route-entries-datasource%d"
}
resource "apsarastack_vpc" "default" {
   name = "${var.name}"
   cidr_block = "10.1.0.0/21"
}
resource "apsarastack_vswitch" "default" {
   vpc_id = "${apsarastack_vpc.default.id}"
   cidr_block = "10.1.1.0/24"
   availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
   name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
   name = "${var.name}"
   description = "${var.name}"
   vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_security_group_rule" "default" {
   type = "ingress"
   ip_protocol = "tcp"
   nic_type = "intranet"
   policy = "accept"
   port_range = "22/22"
   priority = 1
   security_group_id = "${apsarastack_security_group.default.id}"
   cidr_ip = "0.0.0.0/0"
}
resource "apsarastack_instance" "default" {
   # cn-beijing
   security_groups = ["${apsarastack_security_group.default.id}"]
   vswitch_id = "${apsarastack_vswitch.default.id}"
   # series III
   instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
   internet_max_bandwidth_out = 5
   system_disk_category = "cloud_efficiency"
   image_id = "${data.apsarastack_images.default.images.0.id}"
   instance_name = "${var.name}"
}
resource "apsarastack_route_entry" "default" {
   route_table_id = "${apsarastack_vpc.default.route_table_id}"
   destination_cidrblock = "172.11.1.1/32"
   nexthop_type = "Instance"
   nexthop_id = "${apsarastack_instance.default.id}"
}
data "apsarastack_route_entries" "default" {
  %s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRouteEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#":                "1",
		"entries.0.route_table_id": CHECKSET,
		"entries.0.cidr_block":     CHECKSET,
		"entries.0.instance_id":    CHECKSET,
		"entries.0.status":         CHECKSET,
		"entries.0.type":           "Custom",
		"entries.0.next_hop_type":  "Instance",
	}
}

var fakeRouteEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#": "0",
	}
}

var routeEntriesCheckInfo = dataSourceAttr{
	resourceId:   "data.apsarastack_route_entries.default",
	existMapFunc: existRouteEntriesMapFunc,
	fakeMapFunc:  fakeRouteEntriesMapFunc,
}

package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccApsaraStackSnatEntriesDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	snatIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${apsarastack_snat_entry.default.snat_table_id}"`,
			"snat_ip":       `"${apsarastack_snat_entry.default.snat_ip}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${apsarastack_snat_entry.default.snat_table_id}"`,
			"snat_ip":       `"${apsarastack_snat_entry.default.snat_ip}_fake"`,
		}),
	}
	//
	//sourceCidrConf := dataSourceTestAccConfig{
	//	existConfig: testAccCheckApsaraStackSnatEntriesBasicConfig(rand, map[string]string{
	//		"snat_table_id": `"${apsarastack_snat_entry.default.snat_table_id}"`,
	//		"snat_ip":       `"${apsarastack_snat_entry.default.snat_ip}"`,
	//		"source_cidr":   `"172.16.0.0/21"`,
	//	}),
	//	fakeConfig: testAccCheckApsaraStackSnatEntriesBasicConfig(rand, map[string]string{
	//		"snat_table_id": `"${apsarastack_snat_entry.default.snat_table_id}"`,
	//		"snat_ip":       `"${apsarastack_snat_entry.default.snat_ip}"`,
	//		"source_cidr":   `"172.16.0.0/20"`,
	//	}),
	//}
	//allConf := dataSourceTestAccConfig{
	//	existConfig: testAccCheckApsaraStackSnatEntriesBasicConfig(rand, map[string]string{
	//		"snat_table_id": `"${apsarastack_snat_entry.default.snat_table_id}"`,
	//		"source_cidr":   `"172.16.0.0/21"`,
	//	}),
	//	fakeConfig: testAccCheckApsaraStackSnatEntriesBasicConfig(rand, map[string]string{
	//		"snat_table_id": `"${apsarastack_snat_entry.default.snat_table_id}"`,
	//		"source_cidr":   `"172.16.0.0/21"`,
	//	}),
	//}

	snatEntriesCheckInfo.dataSourceTestCheck(t, rand, snatIpConf /*,sourceCidrConf/*,allConf*/)

}

func testAccCheckApsaraStackSnatEntriesBasicConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccForSnatEntriesDatasource%d"
}

data "apsarastack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
	vpc_id = "${apsarastack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "apsarastack_nat_gateway" "default" {
	vpc_id = "${apsarastack_vpc.default.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "apsarastack_eip" "default" {
	name = "${var.name}"
}

resource "apsarastack_eip_association" "default" {
	allocation_id = "${apsarastack_eip.default.id}"
	instance_id = "${apsarastack_nat_gateway.default.id}"
}

resource "apsarastack_snat_entry" "default" {
	snat_table_id = "${apsarastack_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${apsarastack_vswitch.default.id}"
	snat_ip = "${apsarastack_eip.default.ip_address}"
}

data "apsarastack_snat_entries" "default" {
    %s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existSnatEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		//"ids.#":                 "0",
		//"entries.#":             "0",
		"entries.0.id":          CHECKSET,
		"entries.0.status":      "Available",
		"entries.0.source_cidr": "172.16.0.0/21",
	}
}

var fakeSnatEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":     "0",
		"entries.#": "0",
	}
}

var snatEntriesCheckInfo = dataSourceAttr{
	resourceId:   "data.apsarastack_snat_entries.default",
	existMapFunc: existSnatEntriesMapFunc,
	fakeMapFunc:  fakeSnatEntriesMapFunc,
}

package apsarastack

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccApsaraStackSlbsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_slb.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_slb.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_slb.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_slb.default.id}_fake"]`,
		}),
	}

	vpcIDConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${apsarastack_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${apsarastack_vpc.default.id}_fake"`,
		}),
	}

	vswitchConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${apsarastack_slb.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${apsarastack_slb.default.vswitch_id}_fake"`,
		}),
	}

	netWorkTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"name_regex":   `"${apsarastack_slb.default.name}"`,
			"network_type": `"vpc"`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"name_regex":   `"${apsarastack_slb.default.name}"`,
			"network_type": `"classic"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_slb.default.name}"`,
			"tags":       `{tag_f = 6}`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_slb.default.name}"`,
			"tags":       `{tag_f = 0}`,
		}),
	}

	masterZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"name_regex":               `"${apsarastack_slb.default.name}"`,
			"master_availability_zone": `"${data.apsarastack_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"name_regex":               `"${apsarastack_slb.default.name}"`,
			"master_availability_zone": `"${data.apsarastack_zones.default.zones.0.id}_fake"`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${apsarastack_slb.default.name}"`,
			"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("APSARASTACK_RESOURCE_GROUP_ID")),
		}),
		fakeConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${apsarastack_slb.default.name}"`,
			"resource_group_id": fmt.Sprintf(`"%s_fake"`, os.Getenv("APSARASTACK_RESOURCE_GROUP_ID")),
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"name_regex":               `"${apsarastack_slb.default.name}"`,
			"ids":                      `["${apsarastack_slb.default.id}"]`,
			"vswitch_id":               `"${apsarastack_slb.default.vswitch_id}"`,
			"vpc_id":                   `"${apsarastack_vpc.default.id}"`,
			"network_type":             `"vpc"`,
			"tags":                     `{tag_f = 6}`,
			"master_availability_zone": `"${data.apsarastack_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbDataSourceConfig(rand, map[string]string{
			"name_regex":               `"${apsarastack_slb.default.name}_fake"`,
			"ids":                      `["${apsarastack_slb.default.id}"]`,
			"vswitch_id":               `"${apsarastack_slb.default.vswitch_id}"`,
			"vpc_id":                   `"${apsarastack_vpc.default.id}"`,
			"network_type":             `"vpc"`,
			"tags":                     `{tag_f = 6}`,
			"master_availability_zone": `"${data.apsarastack_zones.default.zones.0.id}"`,
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slbs.#":                          "1",
			"ids.#":                           "1",
			"names.#":                         "1",
			"slbs.0.id":                       CHECKSET,
			"slbs.0.name":                     fmt.Sprintf("tf-testAccCheckApsaraStackSlbsDataSourceBasic-%d", rand),
			"slbs.0.region_id":                CHECKSET,
			"slbs.0.master_availability_zone": CHECKSET,
			"slbs.0.slave_availability_zone":  CHECKSET,
			"slbs.0.network_type":             "vpc",
			"slbs.0.vpc_id":                   CHECKSET,
			"slbs.0.vswitch_id":               CHECKSET,
			"slbs.0.address":                  CHECKSET,
			"slbs.0.creation_time":            CHECKSET,

			"slbs.0.tags.%":     "8",
			"slbs.0.tags.tag_a": "1",
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slbs.#":  "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var slbsRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_slbs.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbsRecordsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, vpcIDConf, vswitchConf, netWorkTypeConf, tagsConf, masterZoneConf, resourceGroupIdConf, allConf)

}

func testAccCheckApsaraStackSlbDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckApsaraStackSlbsDataSourceBasic-%d"
}

data "apsarastack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "apsarastack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
  master_zone_id = "${data.apsarastack_zones.default.zones.0.id}"
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
  }
}

data "apsarastack_slbs" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

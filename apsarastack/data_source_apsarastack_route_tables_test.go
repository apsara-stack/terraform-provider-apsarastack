package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccApsaraStackRouteTablesDataSourceBasic(t *testing.T) {
	preCheck := func() {
		testAccPreCheck(t)
	}
	rand := acctest.RandInt()

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${apsarastack_route_table.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${apsarastack_route_table.default.name}_fake"`,
		}),
	}

	vpcIdConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${apsarastack_route_table.default.name}"`,
			"vpc_id":     `"${apsarastack_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${apsarastack_route_table.default.name}"`,
			"vpc_id":     `"${apsarastack_vpc.default.id}_fake"`,
		}),
	}

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"ids": `[ "${apsarastack_route_table.default.id}" ]`,
		}),
		fakeConfig: testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"ids": `[ "${apsarastack_route_table.default.id}_fake" ]`,
		}),
	}

	tagsConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${apsarastack_route_table.default.name}"`,
			"tags": `{
							Created = "TF"
							For 	= "acceptance test"
					  }`,
		}),
		fakeConfig: testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${apsarastack_route_table.default.name}"`,
			"tags": `{
							Created = "TF-fake"
							For 	= "acceptance test-fake"
					  }`,
		}),
	}

	resourceGroupIdConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${apsarastack_route_table.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${apsarastack_route_table.default.name}"`,
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${apsarastack_route_table.default.name}"`,
			"vpc_id":     `"${apsarastack_vpc.default.id}"`,
			"ids":        `[ "${apsarastack_route_table.default.id}" ]`,
		}),
		fakeConfig: testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${apsarastack_route_table.default.name}_fake"`,
			"vpc_id":     `"${apsarastack_vpc.default.id}"`,
			"ids":        `[ "${apsarastack_route_table.default.id}" ]`,
		}),
	}

	routeTablesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConfig, vpcIdConfig, idsConfig, tagsConfig, resourceGroupIdConfig, allConfig)
}

func testAccCheckApsaraStackRouteTablesDataSourceConfigBaisc(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccRouteTablesDatasource%d"
}

resource "apsarastack_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "apsarastack_route_table" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  name = "${var.name}"
  description = "${var.name}_description"
  tags 		= {
		Created = "TF"
		For 	= "acceptance test"
  }
}

data "apsarastack_route_tables" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRouteTablesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                     "0",
		"names.#":                   "0",
		"tables.#":                  "0",
		"tables.0.id":               CHECKSET,
		"tables.0.route_table_type": CHECKSET,
		"tables.0.creation_time":    CHECKSET,
		"tables.0.router_id":        CHECKSET,
		"tables.0.name":             fmt.Sprintf("tf-testAccRouteTablesDatasource%d", rand),
		"tables.0.description":      fmt.Sprintf("tf-testAccRouteTablesDatasource%d_description", rand),
	}
}

var fakeRouteTablesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":    "0",
		"names.#":  "0",
		"tables.#": "0",
	}
}

var routeTablesCheckInfo = dataSourceAttr{
	resourceId:   "data.apsarastack_route_tables.default",
	existMapFunc: existRouteTablesMapFunc,
	fakeMapFunc:  fakeRouteTablesMapFunc,
}

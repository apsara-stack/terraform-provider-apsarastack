package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"strings"
	"testing"
)

func TestAccApsaraStackEssScalingconfigurationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${apsarastack_ess_scaling_configuration.default.scaling_group_id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${apsarastack_ess_scaling_configuration.default.scaling_group_id}_fake"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_ess_scaling_configuration.default.scaling_configuration_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_ess_scaling_configuration.default.scaling_configuration_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_ess_scaling_configuration.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_ess_scaling_configuration.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${apsarastack_ess_scaling_configuration.default.scaling_group_id}"`,
			"ids":              `["${apsarastack_ess_scaling_configuration.default.id}"]`,
			"name_regex":       `"${apsarastack_ess_scaling_configuration.default.scaling_configuration_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${apsarastack_ess_scaling_configuration.default.scaling_group_id}"`,
			"ids":              `["${apsarastack_ess_scaling_configuration.default.id}_fake"]`,
			"name_regex":       `"${apsarastack_ess_scaling_configuration.default.scaling_configuration_name}"`,
		}),
	}

	var existEssScalingconfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                       "1",
			"names.#":                                     "1",
			"configurations.#":                            "1",
			"configurations.0.name":                       fmt.Sprintf("tf-testAccDataSourceEssScalingRules-%d", rand),
			"configurations.0.scaling_group_id":           CHECKSET,
			"configurations.0.image_id":                   CHECKSET,
			"configurations.0.instance_type":              CHECKSET,
			"configurations.0.security_group_id":          CHECKSET,
			"configurations.0.creation_time":              CHECKSET,
			"configurations.0.system_disk_category":       CHECKSET,
			"configurations.0.system_disk_size":           CHECKSET,
			"configurations.0.internet_max_bandwidth_in":  CHECKSET,
			"configurations.0.internet_max_bandwidth_out": CHECKSET,
			"configurations.0.data_disks.#":               "0",
		}
	}

	var fakeEssScalingconfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"configurations.#": "0",
			"ids.#":            "0",
			"names.#":          "0",
		}
	}

	var essScalingconfigurationsCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_ess_scaling_configurations.default",
		existMapFunc: existEssScalingconfigurationsMapFunc,
		fakeMapFunc:  fakeEssScalingconfigurationsMapFunc,
	}

	essScalingconfigurationsCheckInfo.dataSourceTestCheck(t, rand, scalingGroupIdConf, nameRegexConf, idsConf, allConf)
}

func testAccCheckApsaraStackEssScalingconfigurationsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s
variable "name" {
	default = "tf-testAccDataSourceEssScalingRules-%d"
}
resource "apsarastack_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${apsarastack_vswitch.default.id}"]
}
resource "apsarastack_ess_scaling_configuration" "default"{
	scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
	scaling_configuration_name = "${var.name}"
	image_id = "${data.apsarastack_images.default.images.0.id}"
	instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
	security_group_id = "${apsarastack_security_group.default.id}"
	force_delete = true
}
data "apsarastack_ess_scaling_configurations" "default"{
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}

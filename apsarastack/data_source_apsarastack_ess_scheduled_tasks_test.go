package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccApsaraStackEssScheduledtasksDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"scheduled_task_id": `"${apsarastack_ess_scheduled_task.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"scheduled_task_id": `"${apsarastack_ess_scheduled_task.default.id}_fake"`,
		}),
	}
	actionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"scheduled_action": `"${apsarastack_ess_scheduled_task.default.scheduled_action}"`,
		}),
		fakeConfig: testAccCheckApsaraStackEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"scheduled_action": `"${apsarastack_ess_scheduled_task.default.scheduled_action}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_ess_scheduled_task.default.scheduled_task_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_ess_scheduled_task.default.scheduled_task_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_ess_scheduled_task.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_ess_scheduled_task.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"scheduled_action":  `"${apsarastack_ess_scheduled_task.default.scheduled_action}"`,
			"ids":               `["${apsarastack_ess_scheduled_task.default.id}"]`,
			"name_regex":        `"${apsarastack_ess_scheduled_task.default.scheduled_task_name}"`,
			"scheduled_task_id": `"${apsarastack_ess_scheduled_task.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"scheduled_action":  `"${apsarastack_ess_scheduled_task.default.scheduled_action}_fake"`,
			"ids":               `["${apsarastack_ess_scheduled_task.default.id}"]`,
			"name_regex":        `"${apsarastack_ess_scheduled_task.default.scheduled_task_name}"`,
			"scheduled_task_id": `"${apsarastack_ess_scheduled_task.default.id}"`,
		}),
	}

	var existEssScheduledTasksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"tasks.#":                        "1",
			"tasks.0.name":                   fmt.Sprintf("tf-testAccDataSourceEssScheduledTasks-%d", rand),
			"tasks.0.id":                     CHECKSET,
			"tasks.0.scheduled_action":       CHECKSET,
			"tasks.0.launch_expiration_time": CHECKSET,
			"tasks.0.launch_time":            "2020-02-21T11:37Z",
			"tasks.0.max_value":              CHECKSET,
			"tasks.0.min_value":              CHECKSET,
			"tasks.0.task_enabled":           CHECKSET,
		}
	}

	var fakeEssScheduledTasksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"tasks.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var essScheduledTasksCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_ess_scheduled_tasks.default",
		existMapFunc: existEssScheduledTasksMapFunc,
		fakeMapFunc:  fakeEssScheduledTasksMapFunc,
	}

	essScheduledTasksCheckInfo.dataSourceTestCheck(t, rand, idConf, actionConf, nameRegexConf, idsConf, allConf)
}

func testAccCheckApsaraStackEssScheduledTasksDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssScheduledTasks-%d"
}

resource "apsarastack_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${apsarastack_vswitch.default.id}"]
}
resource "apsarastack_ess_scaling_rule" "default" {
  scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
  adjustment_type  = "TotalCapacity"
  adjustment_value = 2
  cooldown         = 60
}

resource "apsarastack_ess_scheduled_task" "default" {
  scheduled_action    = "${apsarastack_ess_scaling_rule.default.ari}"
  launch_time         = "2020-02-21T11:37Z"
  scheduled_task_name = "${var.name}"
}

data "apsarastack_ess_scheduled_tasks" "default"{
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}

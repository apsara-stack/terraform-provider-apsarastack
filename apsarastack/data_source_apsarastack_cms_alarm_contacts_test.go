package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccApsaraStackCmsAlarmContacts_basic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_cms_alarm_contact.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_cms_alarm_contact.default.id}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_cms_alarm_contact.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_cms_alarm_contact.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_cms_alarm_contact.default.id}"`,
			"ids":        `["${apsarastack_cms_alarm_contact.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_cms_alarm_contact.default.id}_fake"`,
			"ids":        `["${apsarastack_cms_alarm_contact.default.id}_fake"]`,
		}),
	}

	var existcmsAlarmContactsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"contacts.#":                    "1",
			"contacts.0.id":                 CHECKSET,
			"contacts.0.alarm_contact_name": CHECKSET,
			"contacts.0.describe":           "For Test",
		}
	}

	var fakecmsAlarmContactsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var cmsAlarmContactsCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_cms_alarm_contacts.default",
		existMapFunc: existcmsAlarmContactsMapFunc,
		fakeMapFunc:  fakecmsAlarmContactsMapFunc,
	}

	cmsAlarmContactsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckApsaraStackCmsAlarmContactsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
		variable "name" {
			default = "tf-testAccCmsAlarmContactBisic-%d"
		}
		resource "apsarastack_cms_alarm_contact" "default" {
			alarm_contact_name = var.name
		    describe           = "For Test"
		    channels_mail      = "hello.uuuu@aaa.com"
			lifecycle {
				ignore_changes = [channels_mail]
  			}	
		}

		data "apsarastack_cms_alarm_contacts" "default" {
		  %s
		}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

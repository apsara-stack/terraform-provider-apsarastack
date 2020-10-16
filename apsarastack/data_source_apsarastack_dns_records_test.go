package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccApsaraStackDnsRecordsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.apsarastack_dns_records.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc%sdns%v.abc", defaultRegionToTest, rand),
		dataSourceDnsRecordsConfigDependence)

	domainNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
		}),
	}

	hostRecordRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name":       "${apsarastack_dns_record.default.name}",
			"host_record_regex": "^ali",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name":       "${apsarastack_dns_record.default.name}",
			"host_record_regex": "anyother",
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
			"type":        "CNAME",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
			"type":        "TXT",
		}),
	}

	valueRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
			"value_regex": "^mail",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
			"value_regex": "anyother",
		}),
	}

	lineConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
			"line":        "default",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
			"line":        "telecom",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
			"status":      "enable",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
			"status":      "disable",
		}),
	}

	isLockConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
			"is_locked":   "false",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
			"is_locked":   "true",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
			"ids":         []string{"${apsarastack_dns_record.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${apsarastack_dns_record.default.name}",
			"ids":         []string{"${apsarastack_dns_record.default.id}-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name":       "${apsarastack_dns_record.default.name}",
			"host_record_regex": "^ali",
			"value_regex":       "^mail",
			"type":              "CNAME",
			"line":              "default",
			"status":            "enable",
			"is_locked":         "false",
			"ids":               []string{"${apsarastack_dns_record.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name":       "${apsarastack_dns_record.default.name}",
			"host_record_regex": "^ali",
			"value_regex":       "^mail",
			"type":              "CNAME",
			"line":              "default",
			"status":            "enable",
			"is_locked":         "true",
			"ids":               []string{"${apsarastack_dns_record.default.id}"},
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"ids.0":                 CHECKSET,
			"urls.#":                "1",
			"urls.0":                fmt.Sprintf("alimail.tf-testacc%sdns%d.abc", defaultRegionToTest, rand),
			"records.#":             "1",
			"domain_name":           fmt.Sprintf("tf-testacc%sdns%d.abc", defaultRegionToTest, rand),
			"records.0.locked":      "false",
			"records.0.host_record": "alimail",
			"records.0.type":        "CNAME",
			"records.0.value":       "mail.mxhichin.com",
			"records.0.record_id":   CHECKSET,
			"records.0.ttl":         "600",
			"records.0.priority":    "0",
			"records.0.line":        "default",
			"records.0.status":      "enable",
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"urls.#":    "0",
			"records.#": "0",
		}
	}

	var dnsRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	dnsRecordsCheckInfo.dataSourceTestCheck(t, rand, domainNameConf, hostRecordRegexConf, typeConf, valueRegexConf, valueRegexConf,
		lineConf, statusConf, isLockConf, idsConf, allConf)
}

func dataSourceDnsRecordsConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "apsarastack_dns" "default" {
  name = "%s"
}

resource "apsarastack_dns_record" "default" {
  name = "${apsarastack_dns.default.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
}
`, name)
}

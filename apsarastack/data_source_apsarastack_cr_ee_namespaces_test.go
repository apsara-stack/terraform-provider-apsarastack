package apsarastack

//import (
//	"fmt"
//	"testing"
//
//	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
//)
//
//func TestAccApsaraStackCrEENamespacesDataSource(t *testing.T) {
//	rand := acctest.RandIntRange(1000000, 9999999)
//	namespaceName := fmt.Sprintf("tf-testacc-cr-ee-ns-%d", rand)
//	resourceId := "data.apsarastack_cr_ee_namespaces.default"
//
//	testAccConfig := dataSourceTestAccConfigFunc(resourceId, namespaceName,
//		dataSourceCrEENamespacesConfigDependence)
//
//	nameRegexConf := dataSourceTestAccConfig{
//		existConfig: testAccConfig(map[string]interface{}{
//			"instance_id": "${data.apsarastack_cr_ee_instances.default.ids.0}",
//			"name_regex":  "${apsarastack_cr_ee_namespace.default.name}",
//		}),
//		fakeConfig: testAccConfig(map[string]interface{}{
//			"instance_id": "${data.apsarastack_cr_ee_instances.default.ids.0}",
//			"name_regex":  "${apsarastack_cr_ee_namespace.default.name}-fake",
//		}),
//	}
//
//	idsConf := dataSourceTestAccConfig{
//		existConfig: testAccConfig(map[string]interface{}{
//			"instance_id": "${data.apsarastack_cr_ee_instances.default.ids.0}",
//		}),
//		fakeConfig: testAccConfig(map[string]interface{}{
//			"instance_id": "${data.apsarastack_cr_ee_instances.default.ids.0}",
//			"ids":         []string{"test-id-fake"},
//		}),
//	}
//
//	allConf := dataSourceTestAccConfig{
//		existConfig: testAccConfig(map[string]interface{}{
//			"instance_id": "${data.apsarastack_cr_ee_instances.default.ids.0}",
//			"name_regex":  "${apsarastack_cr_ee_namespace.default.name}",
//		}),
//		fakeConfig: testAccConfig(map[string]interface{}{
//			"instance_id": "${data.apsarastack_cr_ee_instances.default.ids.0}",
//			"name_regex":  "${apsarastack_cr_ee_namespace.default.name}-fake",
//			"ids":         []string{"test-id-fake"},
//		}),
//	}
//
//	var existCrEENamespacesMapFunc = func(rand int) map[string]string {
//		return map[string]string{
//			"ids.#":                           "1",
//			"ids.0":                           CHECKSET,
//			"names.#":                         "1",
//			"names.0":                         namespaceName,
//			"namespaces.#":                    "1",
//			"namespaces.0.name":               namespaceName,
//			"namespaces.0.default_visibility": "PRIVATE",
//			"namespaces.0.auto_create":        "true",
//			"namespaces.0.instance_id":        CHECKSET,
//		}
//	}
//
//	var fakeCrEENamespacesMapFunc = func(rand int) map[string]string {
//		return map[string]string{
//			"ids.#":        "0",
//			"names.#":      "0",
//			"namespaces.#": "0",
//		}
//	}
//
//	var crEENamespacesCheckInfo = dataSourceAttr{
//		resourceId:   resourceId,
//		existMapFunc: existCrEENamespacesMapFunc,
//		fakeMapFunc:  fakeCrEENamespacesMapFunc,
//	}
//	preCheck := func() {
//		testAccPreCheckWithCrEE(t)
//	}
//	crEENamespacesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
//}
//
//func dataSourceCrEENamespacesConfigDependence(name string) string {
//	return fmt.Sprintf(`
//	variable "name" {
//		default = "%s"
//	}
//
//	data "apsarastack_cr_ee_instances" "default" {
//	}
//
//	resource "apsarastack_cr_ee_namespace" "default" {
//		instance_id = "${data.apsarastack_cr_ee_instances.default.ids.0}"
//		name = "${var.name}"
//		auto_create	= true
//		default_visibility = "PRIVATE"
//	}
//	`, name)
//}

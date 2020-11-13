package apsarastack

import (
	"strings"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccApsaraStackKVStoreInstancesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)

	KvstoreNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex": `"${apsarastack_kvstore_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex": `"${apsarastack_kvstore_instance.default.instance_name}-fake"`,
		}),
	}

	KvstoreIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"ids": `["${apsarastack_kvstore_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"ids": `["${apsarastack_kvstore_instance.default.id}-fake"]`,
		}),
	}

	KvstoreStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex": `"${apsarastack_kvstore_instance.default.instance_name}"`,
			"status":     `"Normal"`,
		}),
		fakeConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex": `"${apsarastack_kvstore_instance.default.instance_name}"`,
			"status":     `"Creating"`,
		}),
	}

	KvstoreInstanceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex":    `"${apsarastack_kvstore_instance.default.instance_name}"`,
			"instance_type": `"Redis"`,
		}),
		fakeConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex":    `"${apsarastack_kvstore_instance.default.instance_name}"`,
			"instance_type": `"Memcache"`,
		}),
	}

	KvstoreVpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex": `"${apsarastack_kvstore_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex": `"${apsarastack_kvstore_instance.default.instance_name}"`,
		}),
	}

	KvstoreVswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex": `"${apsarastack_kvstore_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex": `"${apsarastack_kvstore_instance.default.instance_name}"`,
		}),
	}

	KvstoreTagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex": `"${apsarastack_kvstore_instance.default.instance_name}"`,
			"tags": `{
				Created = "TF"
				For 	= "acceptance test"
			}`,
		}),
		fakeConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex": `"${apsarastack_kvstore_instance.default.instance_name}"`,
			"tags": `{
				Created = "TF"
				For 	= "acceptance test fake"
			}`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex":    `"${apsarastack_kvstore_instance.default.instance_name}"`,
			"ids":           `["${apsarastack_kvstore_instance.default.id}"]`,
			"status":        `"Normal"`,
			"instance_type": `"Redis"`,
			"tags": `{
				Created = "TF"
				For 	= "acceptance test"
			}`}),
		fakeConfig: testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8), map[string]string{
			"name_regex":    `"${apsarastack_kvstore_instance.default.instance_name}-fake"`,
			"ids":           `["${apsarastack_kvstore_instance.default.id}"]`,
			"status":        `"Normal"`,
			"instance_type": `"Redis"`,
			"tags": `{
				Created = "TF"
				For 	= "acceptance test"
			}`}),
	}

	kvstoreInstanceCheckInfo.dataSourceTestCheck(t, rand, KvstoreNameConf, KvstoreIdsConf, KvstoreStatusConf, KvstoreInstanceTypeConf, KvstoreVpcIdConf, KvstoreVswitchIdConf, KvstoreTagsConf, allConf)
}

func testAccCheckApsaraStackKVStoreInstanceDataSourceConfig(rand int, instanceType, instanceClass, engineVersion string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccCheckApsaraStackRKVInstancesDataSourceConfig-%d"
	}
	resource "apsarastack_kvstore_instance" "default" {
		//instance_class = "%s"
		instance_name  = "${var.name}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
		tags 			= {
			Created = "TF"
			For 	= "acceptance test"
		}
	}
	data "apsarastack_kvstore_instances" "default" {
	  %s
	}
`, KVStoreCommonTestCase, rand, instanceClass, instanceType, engineVersion, strings.Join(pairs, "\n  "))
	return config
}

const testAccCheckApsaraStackRKVInstancesDataSourceEmpty = `
data "apsarastack_kvstore_instances" "default" {
  name_regex = "^tf-testacc-fake-name"
}
`

var existKVstoreRecordsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                         "1",
		"names.#":                       "1",
		"instances.#":                   "1",
		"instances.0.id":                CHECKSET,
		"instances.0.name":              fmt.Sprintf("tf-testAccCheckApsaraStackRKVInstancesDataSourceConfig-%d", rand),
		"instances.0.instance_type":     string(KVStoreRedis),
		"instances.0.charge_type":       string(PostPaid),
		"instances.0.region_id":         CHECKSET,
		"instances.0.create_time":       CHECKSET,
		"instances.0.expire_time":       "",
		"instances.0.status":            string(Normal),
		"instances.0.availability_zone": CHECKSET,
		"instances.0.private_ip":        CHECKSET,
		"instances.0.port":              CHECKSET,
		"instances.0.user_name":         CHECKSET,
		"instances.0.capacity":          CHECKSET,
		"instances.0.bandwidth":         CHECKSET,
		"instances.0.connections":       CHECKSET,
		"instances.0.connection_domain": CHECKSET,
	}
}

var fakeKVstoreRecordsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"names.#":     "0",
		"instances.#": "0",
	}
}

var kvstoreInstanceCheckInfo = dataSourceAttr{
	resourceId:   "data.apsarastack_kvstore_instances.default",
	existMapFunc: existKVstoreRecordsMapFunc,
	fakeMapFunc:  fakeKVstoreRecordsMapFunc,
}

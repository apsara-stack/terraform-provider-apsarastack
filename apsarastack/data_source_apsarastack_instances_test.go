package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"strings"
	"testing"

	"fmt"
)

func TestAccApsaraStackInstancesDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_instance.default.instance_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_instance.default.id}" ]`,
		}),
		fakeConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_instance.default.id}_fake" ]`,
		}),
	}

	imageIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckApsaraStackInstancesDataSource%d"`, rand),
			"image_id":   `"${data.apsarastack_images.default.images.0.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckApsaraStackInstancesDataSource%d"`, rand),
			"image_id":   `"${data.apsarastack_images.default.images.0.id}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckApsaraStackInstancesDataSource%d"`, rand),
			"status":     `"Running"`,
		}),
		fakeConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckApsaraStackInstancesDataSource%d"`, rand),
			"status":     `"Stopped"`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${apsarastack_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${apsarastack_vpc.default.id}_fake"`,
		}),
	}

	vSwitchConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${apsarastack_vswitch.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${apsarastack_vswitch.default.id}_fake"`,
		}),
	}

	availabilityZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"name_regex":        fmt.Sprintf(`"tf-testAccCheckApsaraStackInstancesDataSource%d"`, rand),
			"availability_zone": `"${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"`,
		}),
		fakeConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"name_regex":        fmt.Sprintf(`"tf-testAccCheckApsaraStackInstancesDataSource%d"`, rand),
			"availability_zone": `"${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}_fake"`,
		}),
	}

	ramRoleNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"ram_role_name": `"${apsarastack_instance.default.role_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackInstancesDataSourceConfig(rand, map[string]string{
			"ram_role_name": `"${apsarastack_instance.default.role_name}_fake"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackInstancesDataSourceConfigWithTag(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckApsaraStackInstancesDataSource%d"`, rand),
		},
			`tags = {
				from = "datasource"
				usage1 = "test"
				usage2 = "test"
				usage3 = "test"
				usage4 = "test"
				usage5 = "test"
				usage6 = "test"
			}`,
		),
		fakeConfig: testAccCheckApsaraStackInstancesDataSourceConfigWithTag(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckApsaraStackInstancesDataSource%d"`, rand),
		},
			`tags = {
				from = "datasource_fake"
				usage1 = "test"
				usage2 = "test"
				usage3 = "test"
				usage4 = "test"
				usage5 = "test"
				usage6 = "test"
			}`,
		),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackInstancesDataSourceConfigWithTag(rand, map[string]string{
			"ids":               `[ "${apsarastack_instance.default.id}" ]`,
			"name_regex":        fmt.Sprintf(`"tf-testAccCheckApsaraStackInstancesDataSource%d"`, rand),
			"image_id":          `"${data.apsarastack_images.default.images.0.id}"`,
			"status":            `"Running"`,
			"vswitch_id":        `"${apsarastack_vswitch.default.id}"`,
			"availability_zone": `"${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"`,
		},
			`tags = {
				from = "datasource"
				usage1 = "test"
				usage2 = "test"
				usage3 = "test"
				usage4 = "test"
				usage5 = "test"
				usage6 = "test"
			}`,
		),
		fakeConfig: testAccCheckApsaraStackInstancesDataSourceConfigWithTag(rand, map[string]string{
			"ids":               `[ "${apsarastack_instance.default.id}_fake" ]`,
			"name_regex":        fmt.Sprintf(`"tf-testAccCheckApsaraStackInstancesDataSource%d"`, rand),
			"image_id":          `"${data.apsarastack_images.default.images.0.id}"`,
			"status":            `"Running"`,
			"vswitch_id":        `"${apsarastack_vswitch.default.id}"`,
			"availability_zone": `"${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"`,
		},
			`tags = {
				from = "datasource_fake"
				usage1 = "test"
				usage2 = "test"
				usage3 = "test"
				usage4 = "test"
				usage5 = "test"
				usage6 = "test"
			}`,
		),
	}

	instancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, imageIdConf, statusConf,
		vpcIdConf, vSwitchConf, availabilityZoneConf, ramRoleNameConf, tagsConf, allConf)
}

func testAccCheckApsaraStackInstancesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckApsaraStackInstancesDataSource%d"
	}
	resource "apsarastack_instance" "default" {
		availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
		vswitch_id = "${apsarastack_vswitch.default.id}"
		private_ip = "172.16.0.10"
		image_id = "${data.apsarastack_images.default.images.0.id}"
		instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		security_groups = ["${apsarastack_security_group.default.id}"]
		role_name = "${apsarastack_ram_role.default.name}"
		data_disks {
				name  = "${var.name}-disk1"
				size =        "20"
				category =  "cloud_efficiency"
				description = "disk1"
		}
		data_disks {
				name  = "${var.name}-disk2"
				size =        "20"
				category =  "cloud_efficiency"
				description = "disk2"
		}
        tags = {
			from = "datasource"
			usage1 = "test"
			usage2 = "test"
			usage3 = "test"
			usage4 = "test"
			usage5 = "test"
			usage6 = "test"
		}
	}
	
	resource "apsarastack_ram_role" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
	  force = true
	}
	data "apsarastack_instances" "default" {
		%s
	}`, EcsInstanceCommonNoZonesTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckApsaraStackInstancesDataSourceConfigWithTag(rand int, attrMap map[string]string, tags string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckApsaraStackInstancesDataSource%d"
	}
	resource "apsarastack_instance" "default" {
		availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
		vswitch_id = "${apsarastack_vswitch.default.id}"
		private_ip = "172.16.0.10"
		image_id = "${data.apsarastack_images.default.images.0.id}"
		instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		security_groups = ["${apsarastack_security_group.default.id}"]
		data_disks {
				name  = "${var.name}-disk1"
				size =        "20"
				category =  "cloud_efficiency"
				description = "disk1"
		}
		data_disks {
				name  = "${var.name}-disk2"
				size =        "20"
				category =  "cloud_efficiency"
				description = "disk2"
		}
        tags = {
			from = "datasource"
			usage1 = "test"
			usage2 = "test"
			usage3 = "test"
			usage4 = "test"
			usage5 = "test"
			usage6 = "test"
		}
	}
	data "apsarastack_instances" "default" {
		%s
		%s
	}`, EcsInstanceCommonNoZonesTestCase, rand, strings.Join(pairs, "\n  "), tags)
	return config
}

var existInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                  "1",
		"names.#":                                "1",
		"instances.#":                            "1",
		"instances.0.id":                         CHECKSET,
		"instances.0.region_id":                  CHECKSET,
		"instances.0.availability_zone":          CHECKSET,
		"instances.0.private_ip":                 "172.16.0.10",
		"instances.0.status":                     string(Running),
		"instances.0.name":                       fmt.Sprintf("tf-testAccCheckApsaraStackInstancesDataSource%d", rand),
		"instances.0.instance_type":              CHECKSET,
		"instances.0.vswitch_id":                 CHECKSET,
		"instances.0.image_id":                   CHECKSET,
		"instances.0.eip":                        "",
		"instances.0.description":                "",
		"instances.0.security_groups.#":          "1",
		"instances.0.key_name":                   "",
		"instances.0.creation_time":              CHECKSET,
		"instances.0.internet_max_bandwidth_out": "0",
		"instances.0.disk_device_mappings.#":     "3",
	}
}

var fakeInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"names.#":     "0",
		"instances.#": "0",
	}
}

var instancesCheckInfo = dataSourceAttr{
	resourceId:   "data.apsarastack_instances.default",
	existMapFunc: existInstancesMapFunc,
	fakeMapFunc:  fakeInstancesMapFunc,
}

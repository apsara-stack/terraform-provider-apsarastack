package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccApsaraStackNetworkInterfacesDataSourceBasic(t *testing.T) {

	rand := acctest.RandInt()

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
		}),
		fakeConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_network_interface_attachment.default.network_interface_id}_fake" ]`,
		}),
	}

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":         `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"instance_id": `"${apsarastack_instance.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":         `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"instance_id": `"${apsarastack_instance.default.id}_fake"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccNetworkInterfacesBasic%d"`, rand),
		}),
		fakeConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccNetworkInterfacesBasic%d_fake"`, rand),
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":    `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"vpc_id": `"${apsarastack_vpc.default.id}"`,
		}),
	}

	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":        `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"vswitch_id": `"${apsarastack_vswitch.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":        `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"vswitch_id": `"${apsarastack_vswitch.default.id}_fake"`,
		}),
	}

	privateIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":        `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"private_ip": `"192.168.0.2"`,
		}),
		fakeConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":        `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"private_ip": `"192.168.0.1"`,
		}),
	}

	securityGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":               `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"security_group_id": `"${apsarastack_security_group.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":               `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"security_group_id": `"${apsarastack_security_group.default.id}_fake"`,
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":  `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"type": `"Secondary"`,
		}),
		fakeConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":  `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"type": `"Primary"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"tags": fmt.Sprintf(`{
							 TF-VER = "0.11.3%d"
						   }`, rand),
		}),
		fakeConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"tags": fmt.Sprintf(`{
							 TF-VER = "0.11.3%d_fake"
						   }`, rand),
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":        `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"name_regex": fmt.Sprintf(`"tf-testAccNetworkInterfacesBasic%d"`, rand),
			"tags": fmt.Sprintf(`{
							 TF-VER = "0.11.3%d"
						   }`, rand),
			"vswitch_id":        `"${apsarastack_vswitch.default.id}"`,
			"vpc_id":            `"${apsarastack_vpc.default.id}"`,
			"private_ip":        `"192.168.0.2"`,
			"security_group_id": `"${apsarastack_security_group.default.id}"`,
			"type":              `"Secondary"`,
			"instance_id":       `"${apsarastack_instance.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":        `[ "${apsarastack_network_interface_attachment.default.network_interface_id}" ]`,
			"name_regex": fmt.Sprintf(`"tf-testAccNetworkInterfacesBasic%d"`, rand),
			"tags": fmt.Sprintf(`{
							 TF-VER = "0.11.3%d_fake"
						   }`, rand),
			"vpc_id":            `"${apsarastack_vpc.default.id}"`,
			"vswitch_id":        `"${apsarastack_vswitch.default.id}"`,
			"private_ip":        `"192.168.0.2"`,
			"security_group_id": `"${apsarastack_security_group.default.id}"`,
			"type":              `"Primary"`,
			"instance_id":       `"${apsarastack_instance.default.id}"`,
		}),
	}

	networkInterfacesCheckInfo.dataSourceTestCheck(t, rand, idsConf, instanceIdConf, nameRegexConf, vpcIdConf, vswitchIdConf, privateIpConf,
		securityGroupIdConf, typeConf, tagsConf, allConf)
}

func testAccCheckApsaraStackNetworkInterfacesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`


variable "name" {
 default = "tf-testAccNetworkInterfacesBasic"
}

resource "apsarastack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "apsarastack_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "apsarastack_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_network_interface" "default" {
    name = "${var.name}%d"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    security_groups = [ "${apsarastack_security_group.default.id}" ]
	description = "Basic test"
	private_ip = "192.168.0.2"
	tags = {
		TF-VER = "0.11.3%d"
	}
	
}

data "apsarastack_instance_types" "default" {
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  eni_amount = 2
}

data "apsarastack_images" "default" {
  	name_regex  = "^ubuntu_18.*64"
  	most_recent = true
	owners = "system"
}

resource "apsarastack_instance" "default" {
    availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
    security_groups = ["${apsarastack_security_group.default.id}"]
    instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
    system_disk_category = "cloud_efficiency"
    image_id             = "${data.apsarastack_images.default.images.0.image_id}"
    instance_name        = "${var.name}"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    
}

resource "apsarastack_network_interface_attachment" "default" {
    instance_id = "${apsarastack_instance.default.id}"
    network_interface_id = "${apsarastack_network_interface.default.id}"
}

data "apsarastack_network_interfaces" "default"  {
	%s
}`, rand, rand, strings.Join(pairs, "\n  "))
	return config
}

var existNetworkInterfacesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                          "1",
		"names.#":                        "1",
		"interfaces.#":                   "1",
		"interfaces.0.id":                CHECKSET,
		"interfaces.0.name":              fmt.Sprintf("tf-testAccNetworkInterfacesBasic%d", rand),
		"interfaces.0.status":            CHECKSET,
		"interfaces.0.zone_id":           CHECKSET,
		"interfaces.0.public_ip":         "",
		"interfaces.0.private_ip":        "192.168.0.2",
		"interfaces.0.private_ips.#":     "0",
		"interfaces.0.security_groups.#": "1",
		"interfaces.0.description":       "Basic test",
		"interfaces.0.instance_id":       CHECKSET,
		"interfaces.0.creation_time":     CHECKSET,
		"interfaces.0.tags.%":            "1",
		"interfaces.0.tags.TF-VER":       fmt.Sprintf("0.11.3%d", rand),
	}
}

var fakeNetworkInterfacesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"interfaces.#": "0",
		"names.#":      "0",
		"ids.#":        "0",
	}
}

var networkInterfacesCheckInfo = dataSourceAttr{
	resourceId:   "data.apsarastack_network_interfaces.default",
	existMapFunc: existNetworkInterfacesMapFunc,
	fakeMapFunc:  fakeNetworkInterfacesMapFunc,
}

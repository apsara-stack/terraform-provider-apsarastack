package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccApsaraStackNetworkInterfaceAttachmentBasic(t *testing.T) {
	var v ecs.NetworkInterfaceSet
	resourceId := "apsarastack_network_interface_attachment.default"
	ra := resourceAttrInit(resourceId, testAccCheckNetworkInterfaceAttachmentCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkInterfaceAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceAttachmentConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApsaraStackNetworkInterfaceAttachmentMulti(t *testing.T) {
	var v ecs.NetworkInterfaceSet
	resourceId := "apsarastack_network_interface_attachment.default.1"
	ra := resourceAttrInit(resourceId, testAccCheckNetworkInterfaceAttachmentCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkInterfaceAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceAttachmentConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckNetworkInterfaceAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_network_interface_Attachment" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NetworkInterface ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeNetworkInterfaceAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

const testAccNetworkInterfaceAttachmentConfigBasic = `
variable "name" {
  default = "tf-testAccNetworkInterfaceAttachment"
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
    availability_zone = "${reverse(data.apsarastack_zones.default.zones).0.id}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

data "apsarastack_instance_types" "default" {
    availability_zone = "${reverse(data.apsarastack_zones.default.zones).0.id}"
    eni_amount = 2
}

data "apsarastack_images" "default" {
	name_regex  = "^ubuntu_18.*64"
  	most_recent = true
	owners = "system"
}

resource "apsarastack_instance" "default" {
    availability_zone = "${reverse(data.apsarastack_zones.default.zones).0.id}"
    security_groups = ["${apsarastack_security_group.default.id}"]

    instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
    system_disk_category = "cloud_efficiency"
    image_id             = "${data.apsarastack_images.default.images.0.id}"
    instance_name        = "${var.name}"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    internet_max_bandwidth_out = 10
}

resource "apsarastack_network_interface" "default" {
    name = "${var.name}"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    security_groups = [ "${apsarastack_security_group.default.id}" ]
}

resource "apsarastack_network_interface_attachment" "default" {
    instance_id = "${apsarastack_instance.default.id}"
    network_interface_id = "${apsarastack_network_interface.default.id}"
}
`

const testAccNetworkInterfaceAttachmentConfigMulti = `
variable "name" {
  default = "tf-testAccNetworkInterfaceAttachment"
}

variable "number" {
		default = "2"
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
    availability_zone = "${reverse(data.apsarastack_zones.default.zones).0.id}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

data "apsarastack_instance_types" "default" {
    availability_zone = "${reverse(data.apsarastack_zones.default.zones).0.id}"
    eni_amount = 2
}

data "apsarastack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "apsarastack_instance" "default" {
	count = "${var.number}"
    availability_zone = "${reverse(data.apsarastack_zones.default.zones).0.id}"
    security_groups = ["${apsarastack_security_group.default.id}"]

    instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
    system_disk_category = "cloud_efficiency"
    image_id             = "${data.apsarastack_images.default.images.0.id}"
    instance_name        = "${var.name}"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    internet_max_bandwidth_out = 10
}

resource "apsarastack_network_interface" "default" {
    count = "${var.number}"
    name = "${var.name}"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    security_groups = [ "${apsarastack_security_group.default.id}" ]
}

resource "apsarastack_network_interface_attachment" "default" {
	count = "${var.number}"
    instance_id = "${element(apsarastack_instance.default.*.id, count.index)}"
    network_interface_id = "${element(apsarastack_network_interface.default.*.id, count.index)}"
}
`

var testAccCheckNetworkInterfaceAttachmentCheckMap = map[string]string{
	"instance_id":          CHECKSET,
	"network_interface_id": CHECKSET,
}

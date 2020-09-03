package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckKeyPairAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_key_pair_attachment" {
			continue
		}
		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
		ecsService := EcsService{client}

		instanceIds := rs.Primary.Attributes["instance_ids"]

		for _, inst := range instanceIds {
			response, err := ecsService.DescribeInstance(string(inst))
			if err != nil {
				return err
			}

			if response.KeyPairName != "" {
				return fmt.Errorf("Error Key Pair Attachment still exist")
			}

		}
	}

	return nil
}

func TestAccApsaraStackKeyPairAttachmentBasic(t *testing.T) {
	var v ecs.KeyPair
	resourceId := "apsarastack_key_pair_attachment.default"
	ra := resourceAttrInit(resourceId, testAccCheckKeyPairAttachmentBasicMap)
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

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeyPairAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairAttachmentConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

const testAccKeyPairAttachmentConfigBasic = `
data "apsarastack_zones" "default" {
	available_disk_category = "cloud_ssd"
	available_resource_creation= "VSwitch"
}
data "apsarastack_instance_types" "default" {
 	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
}
data "apsarastack_images" "default" {
	name_regex = "^ubuntu_18.*64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "tf-testAccKeyPairAttachmentConfig"
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name = "${var.name}"
  description = "New security group"
  vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_instance" "default" {
  instance_name = "${var.name}-${count.index+1}"
  image_id = "${data.apsarastack_images.default.images.0.id}"
  instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  count = 2
  security_groups = ["${apsarastack_security_group.default.id}"]
  vswitch_id = "${apsarastack_vswitch.default.id}"

  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = 5
  password = "Yourpassword1234"

  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_ssd"
}

resource "apsarastack_key_pair" "default" {
  key_name = "${var.name}"
}

resource "apsarastack_key_pair_attachment" "default" {
  key_name = "${apsarastack_key_pair.default.id}"
  instance_ids = "${apsarastack_instance.default.*.id}"
}
`

var testAccCheckKeyPairAttachmentBasicMap = map[string]string{
	"key_name":       CHECKSET,
	"instance_ids.#": "2",
}

package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"testing"

	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccApsaraStackRamRoleAttachment_basic(t *testing.T) {
	var v *ecs.DescribeInstanceRamRoleResponse
	resourceId := "apsarastack_ram_role_attachment.default"
	ra := resourceAttrInit(resourceId, ramRoleAttachmentMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
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
		//CheckDestroy:  rac.checkResourceDestroy(),
		CheckDestroy: testAccCheckRamRoleAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAscm_RamRoleAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckRamRoleAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	ascmService := RamService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_ram_role_attachment" || rs.Type != "apsarastack_ram_role_attachment" {
			continue
		}
		ascm, err := ascmService.DescribeRamRoleAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if ascm.InstanceRamRoleSets.InstanceRamRoleSet[0].RamRoleName != "" {
			return WrapError(Error("resource  still exist"))
		}
	}

	return nil
}

const testAccCheckAscm_RamRoleAttachment = `
variable "name" {
  default = "Test_ram_role_attachment"
}
data "apsarastack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "apsarastack_images" "default" {
  most_recent = true
  owners = "system"
}
resource "apsarastack_vpc" "default" {
  name = var.name
  cidr_block = "192.168.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id = apsarastack_vpc.default.id
  cidr_block = "192.168.0.0/16"
  availability_zone = data.apsarastack_zones.default.zones[0].id
  name = var.name
}
resource "apsarastack_security_group" "default" {
  name = var.name
  vpc_id = apsarastack_vpc.default.id
}
resource "apsarastack_instance" "default" {
  image_id = data.apsarastack_images.default.images.0.id
  instance_type = "ecs.n4.small"
  instance_name = var.name
  security_groups = [apsarastack_security_group.default.id]
  availability_zone = data.apsarastack_zones.default.zones[0].id
  system_disk_category = "cloud_efficiency"
  system_disk_size = 100
  vswitch_id = apsarastack_vswitch.default.id
}

data "apsarastack_ascm_ram_service_roles" "role" {
  product = "ecs"
}
resource "apsarastack_ram_role_attachment" "default" {
   role_name    = data.apsarastack_ascm_ram_service_roles.role.roles.0.name
   instance_ids = [apsarastack_instance.default.id]
}
`

var ramRoleAttachmentMap = map[string]string{
	"role_name": CHECKSET,
}

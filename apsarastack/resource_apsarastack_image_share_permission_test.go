package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApsaraStackImageSharePermission(t *testing.T) {
	var v *ecs.DescribeImageSharePermissionResponse
	resourceId := "apsarastack_image_share_permission.default"
	ra := resourceAttrInit(resourceId, testAccImageSharePermissionCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageShareByImageId")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccEcsImageShareConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageSharePermissionConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithMultipleAccount(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":   "${apsarastack_image.default.id}",
					"account_id": os.Getenv("APSARASTACK_ACCESS_KEY_2"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_id": CHECKSET,
					}),
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

var testAccImageSharePermissionCheckMap = map[string]string{
	"image_id": CHECKSET,
}

func resourceImageSharePermissionConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "apsarastack_instance_types" "default" {
 	cpu_core_count    = 1
	memory_size       = 2
}
data "apsarastack_images" "default" {
  name_regex  = "^ubuntu*"
  owners      = "system"
}
resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  name              = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_instance" "default" {
  image_id = "${data.apsarastack_images.default.ids[0]}"
  instance_type = "${data.apsarastack_instance_types.default.ids[0]}"
  security_groups = "${[apsarastack_security_group.default.id]}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
  instance_name = "${var.name}"
}
resource "apsarastack_image" "default" {
  instance_id = "${apsarastack_instance.default.id}"
  image_name        = "${var.name}"
}
`, name)
}

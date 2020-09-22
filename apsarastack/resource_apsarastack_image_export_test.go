package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApsaraStackImageExport(t *testing.T) {
	var v ecs.Image
	resourceId := "apsarastack_image_export.default"
	ra := resourceAttrInit(resourceId, testAccExportImageCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}

	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testaccecsimageexportconfigbasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageExportBasicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":   "${apsarastack_image.default.id}",
					"oss_bucket": "${apsarastack_oss_bucket.default.bucket}",
					"oss_prefix": "ecsExport",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_prefix": "ecsExport",
					}),
				),
			},
		},
	})
}

var testAccExportImageCheckMap = map[string]string{
	"image_id":   CHECKSET,
	"oss_bucket": CHECKSET,
}

func resourceImageExportBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "apsarastack_instance_types" "default" {
	cpu_core_count    = 1
	memory_size       = 2
}
data "apsarastack_images" "default" {
 name_regex  = "^ubuntu_18.*64"
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
resource "apsarastack_oss_bucket" "default" {
  bucket = "${var.name}"
}
`, name)
}

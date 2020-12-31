package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func (rc *resourceCheck) testAccCheckKeyPairAttachmentDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		strs := strings.Split(rc.resourceId, ":")
		var resourceType string
		for _, str := range strs {
			if strings.Contains(str, "apsarastack_") {
				resourceType = strings.Trim(str, " ")
				break
			}
		}

		if resourceType == "" {
			return WrapError(Error("The resourceId %s is not correct and it should prefix with apsarastack_", rc.resourceId))
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			outValue, err := rc.callDescribeMethod(rs)
			errorValue := outValue[1]
			if !errorValue.IsNil() {
				err = errorValue.Interface().(error)
				if err != nil {
					if NotFoundError(err) {
						continue
					}
					return WrapError(err)
				}
			} else {
				return WrapError(Error("the resource %s %s was not destroyed ! ", rc.resourceId, rs.Primary.ID))
			}
		}
		return nil
	}
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
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAccKeyPairAttachment%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccKeyPairAttachmentConfigBasic)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.testAccCheckKeyPairAttachmentDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"key_name":     name,
					"instance_ids": "${apsarastack_instance.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}
func testAccKeyPairAttachmentConfigBasic(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}
data "apsarastack_zones" "default" {
	available_disk_category = "cloud_ssd"
	available_resource_creation= "VSwitch"
}
resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}
resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = apsarastack_vpc.default.cidr_block
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
  image_id = "wincore_2004_x64_dtc_en-us_40G_alibase_20201015.raw"
  instance_type = "ecs.n4.xlarge"
  count = 2
  security_groups = ["${apsarastack_security_group.default.id}"]
  vswitch_id = "${apsarastack_vswitch.default.id}"
  internet_max_bandwidth_out = 5
  password = "Yourpassword1234"
  system_disk_category = "cloud_ssd"
}
resource "apsarastack_key_pair" "default" {
  key_name = "${var.name}"
}
`, name)
}

var testAccCheckKeyPairAttachmentBasicMap = map[string]string{
	"key_name":       CHECKSET,
	"instance_ids.#": "2",
}

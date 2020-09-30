package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApsaraStackEssScalingConfigurationUpdate(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "apsarastack_ess_scaling_configuration.default"
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^ubuntu_18",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${apsarastack_ess_scaling_group.default.id}",
					"image_id":          "${data.apsarastack_images.default.images.0.id}",
					"instance_type":     "${data.apsarastack_instance_types.default.instance_types.0.id}",
					"security_group_id": "${apsarastack_security_group.default.id}",
					"force_delete":      "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id", "password", "kms_encrypted_password", "kms_encryption_context"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"active": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"active": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_configuration_name": fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_configuration_name": fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_ssd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]string{{
						"size":                 "20",
						"category":             "cloud_ssd",
						"delete_with_instance": "false",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "1",
						"data_disk.0.size":                 "20",
						"data_disk.0.category":             "cloud_ssd",
						"data_disk.0.delete_with_instance": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_data": `#!/bin/bash\necho \"hello\"\n`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data": "#!/bin/bash\necho \"hello\"\n",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name": "${apsarastack_ram_role.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_name": "${apsarastack_key_pair.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"name": "tf-test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.name": "tf-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":  REMOVEKEY,
					"instance_types": []string{"${data.apsarastack_instance_types.default.instance_types.0.id}", "${data.apsarastack_instance_types.default.instance_types.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":    REMOVEKEY,
						"instance_types.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id": "${data.apsarastack_images.default1.images.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id": REGEXMATCH + "^centos",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":              REMOVEKEY,
					"instance_types":             REMOVEKEY,
					"system_disk_category":       REMOVEKEY,
					"key_name":                   REMOVEKEY,
					"role_name":                  REMOVEKEY,
					"data_disk":                  REMOVEKEY,
					"tags":                       REMOVEKEY,
					"user_data":                  REMOVEKEY,
					"scaling_configuration_name": REMOVEKEY,
					"scaling_group_id":           "${apsarastack_ess_scaling_group.default.id}",
					"image_id":                   "${data.apsarastack_images.default.images.0.id}",
					"instance_type":              "${data.apsarastack_instance_types.default.instance_types.0.id}",
					"security_group_id":          "${apsarastack_security_group.default.id}",
					"force_delete":               "true",
					"override":                   "true",
					"system_disk_size":           "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":                    REMOVEKEY,
						"instance_types.#":                 REMOVEKEY,
						"system_disk_category":             REMOVEKEY,
						"key_name":                         REMOVEKEY,
						"role_name":                        REMOVEKEY,
						"data_disk.#":                      REMOVEKEY,
						"data_disk.0.delete_with_instance": REMOVEKEY,
						"data_disk.0.size":                 REMOVEKEY,
						"data_disk.0.category":             REMOVEKEY,
						"tags.name":                        REMOVEKEY,
						"user_data":                        REMOVEKEY,
						"scaling_configuration_name":       REMOVEKEY,
						"system_disk_size":                 "100",
						"scaling_group_id":                 CHECKSET,
						"instance_type":                    CHECKSET,
						"security_group_id":                CHECKSET,
						"image_id":                         REGEXMATCH + "^ubuntu_18",
						"override":                         "true",
					}),
				),
			},
		},
	})
}

func TestAccApsaraStackEssScalingConfigurationMulti(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "apsarastack_ess_scaling_configuration.default.9"
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^ubuntu_18",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":             "10",
					"scaling_group_id":  "${apsarastack_ess_scaling_group.default.id}",
					"image_id":          "${data.apsarastack_images.default.images.0.id}",
					"instance_type":     "${data.apsarastack_instance_types.default.instance_types.0.id}",
					"security_group_id": "${apsarastack_security_group.default.id}",
					"force_delete":      "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceEssScalingConfigurationConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}
	
	resource "apsarastack_key_pair" "default" {
  		key_name = "${var.name}"
	}
	
	resource "apsarastack_security_group" "default1" {
	  name   = "${var.name}"
	  vpc_id = "${apsarastack_vpc.default.id}"
	}
	data "apsarastack_images" "default1" {
		name_regex  = "^centos.*_64"
  		most_recent = true
  		owners      = "system"
	}
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
	}`, EcsInstanceCommonTestCase, name)
}

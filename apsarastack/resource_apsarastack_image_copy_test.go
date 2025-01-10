package apsarastack

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApsaraStackImageCopyBasic(t *testing.T) {
	var v ecs.Image

	resourceId := "apsarastack_image_copy.default"
	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"apsarastack": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, testAccCopyImageCheckMap)
	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsCopyImageConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageCopyBasicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckImageDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_image_id": "${apsarastack_image.default.id}",
					"description":     fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"image_name":      name,
					"kms_key_id":      "3852c3cd-3ace-468d-8b9b-c301c33a32b2",
					"encrypted":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"image_name":  name,
						"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescriptionChange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescriptionChange", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": fmt.Sprintf("tf-testAccEcsImageConfigBasic%dchange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"image_name": fmt.Sprintf("tf-testAccEcsImageConfigBasic%dchange", rand),
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"image_name":  name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
						"image_name":  name,
					}),
				),
			},
		},
	})
}

func TestAccApsaraStackImageCopyEncrypted(t *testing.T) {
	// var v ecs.Image

	resourceId := "apsarastack_image_copy.default"
	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"apsarastack": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, testAccCopyImageCheckMap)
	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsCopyImageConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageCopyBasicConfigDependenceEncrypted)
	region := os.Getenv("APSARASTACK_REGION")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckImageDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_image_id":       "m-gj001q4w0gr14n74athd",
					"description":           fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"destination_region_id": region,
					"image_name":            name,
					"kms_key_id":            "2208243d-faff-484a-be0b-94ae65d9963e",
					"encrypted":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckImageExistsWithProviders(n string, image *ecs.Image, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No image  ID is set")
		}
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.ApsaraStackClient)
			ecsService := EcsService{client}

			resp, err := ecsService.DescribeImageById(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err
			}

			*image = resp
			return nil
		}
		return fmt.Errorf("image not found")
	}
}

func testAccCheckImageDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckImageDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckImageDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {

	client := provider.Meta().(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_copy_image" {
			continue
		}

		resp, err := ecsService.DescribeImageById(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("image still exist,  ID %s ", resp.ImageId)
		}
	}

	return nil
}

var testAccCopyImageCheckMap = map[string]string{}

func resourceImageCopyBasicConfigDependenceEncrypted(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
}`, name)
}

func resourceImageCopyBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
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
  availability_zone = "cn-wulan-env212-amtest212001-a"
  name              = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_instance" "default" {
  image_id = "${data.apsarastack_images.default.ids[0]}"
  instance_type        =  "ecs.n4.large"
  system_disk_category = "cloud_ssd"
  system_disk_size     = 40
  system_disk_name     = "test_sys_disk"
  security_groups      = "${[apsarastack_security_group.default.id]}"
  instance_name        = "${var.name}_ecs"
  vswitch_id           = "${apsarastack_vswitch.default.id}"
  availability_zone    =  "cn-wulan-env212-amtest212001-a"
  is_outdated          = false
}


resource "apsarastack_image" "default" {

  instance_id = "${apsarastack_instance.default.id}"
  image_name        = "${var.name}"
}
`, name)
}

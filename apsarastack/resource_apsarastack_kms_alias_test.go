package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApsaraStackKmsAlias_basic(t *testing.T) {
	var v kms.KeyMetadata

	resourceId := "apsarastack_kms_alias.default"
	ra := resourceAttrInit(resourceId, kmsAliasBasicMap)

	serviceFunc := func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("alias/tf-testKmsAlias_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceKmsAliadConfigDependence)

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
					"alias_name": name,
					"key_id":     "${apsarastack_kms_key.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alias_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_id": "${apsarastack_kms_key.default1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_id": CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceKmsAliadConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "apsarastack_kms_key" "default" {}

resource "apsarastack_kms_key" "default1" {}
`)
}

var kmsAliasBasicMap = map[string]string{
	"key_id": CHECKSET,
}

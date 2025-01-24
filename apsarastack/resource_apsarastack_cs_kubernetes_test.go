package apsarastack

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccCheckCsK8sDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_cs_kubernetes" {
			continue
		}

		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
		csService := CsService{client}
		log.Printf("repo ID %s", rs.Primary.ID)
		_, err := csService.DescribeCsKubernetes(rs.Primary.ID)

		if err == nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func TestAccApsaraStackCsK8s_Basic(t *testing.T) {
	var v cs.KubernetesClusterDetail
	resourceId := "apsarastack_cs_kubernetes.k8s"
	ra := resourceAttrInit(resourceId, CsK8sMap)
	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCsK8sConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCsK8sConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCsK8sDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  "${var.name}",
					"version":               "1.20.11-aliyun.1",
					"os_type":               "linux",
					"platform":              "AliyunLinux",
					"timeout_mins":          "60",
					"vpc_id":                "${apsarastack_vpc.default.id}",
					"image_id":              "${data.apsarastack_images.default.images.0.id}",
					"master_count":          "3",
					"master_disk_category":  "cloud_efficiency",
					"master_disk_size":      "40",
					"master_instance_types": []string{"ecs.se1ne.large", "ecs.se1ne.large", "ecs.se1ne.large"},
					"master_vswitch_ids":    []string{"${apsarastack_vswitch.default.id},${apsarastack_vswitch.default.id},${apsarastack_vswitch.default.id}"},

					"num_of_nodes":         "1",
					"worker_disk_category": "cloud_efficiency",
					"worker_disk_size":     "40",
					"runtime": []map[string]interface{}{
						{
							"name":    "Docker",
							"version": "19.03.15",
						},
					},
					"worker_instance_types": []string{"ecs.se1ne.large"},
					"worker_vswitch_ids":    []string{"${apsarastack_vswitch.default.id}"},
					"security_group_id":     "${apsarastack_security_group.default.id}",
					"password":              "Alibaba@1688",
					"delete_protection":     "false",
					"pod_cidr":              "172.20.0.0/16",
					"service_cidr":          "172.21.0.0/20",
					"node_cidr_mask":        "24",
					"new_nat_gateway":       "true",
					"slb_internet_enabled":  "true",
					"proxy_mode":            "ipvs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"security_group_id": CHECKSET,
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

func resourceCsK8sConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "apsarastack_zones" default {
  available_resource_creation = "VSwitch"
}

// data "apsarastack_instance_types" "default" {
// 	cpu_core_count    = 1
// 	memory_size       = 1
// }

data "apsarastack_images" "default" {
	name_regex  = "^ubuntu*"
	owners      = "system"
}

resource "apsarastack_vpc" "default" {
cidr_block = "172.16.0.0/16"
name = "${var.name}"
}

resource "apsarastack_vswitch" "default" {
vpc_id            = "${apsarastack_vpc.default.id}"
cidr_block        = "172.16.0.0/24"
availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
name              = "${var.name}"
}

resource "apsarastack_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${apsarastack_vpc.default.id}"
  }

  
variable "cluster_addons" {
  description = "Addon components in kubernetes cluster"
  type = list(object({
    name      = string
    config    = string
  }))
  default = [
    {
      "name"     = "terway",
      "config"   = "",
    },
    {
      "name"     = "csi-plugin",
      "config"   = "",
    },
    {
      "name"     = "csi-provisioner",
      "config"   = "",
    },
    {
      "name"     = "logtail-ds",
      "config"   = "{\"IngressDashboardEnabled\":\"true\",\"sls_project_name\":\"alibaba-test\"}",
    },
    {
      "name"     = "nginx-ingress-controller",
      "config"   = "{\"IngressSlbNetworkType\":\"internet\"}",
    }
  ]
}
`, name)
}

var CsK8sMap = map[string]string{}

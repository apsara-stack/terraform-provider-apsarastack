package apsarastack

import (
	"fmt"
	"log"
	"testing"

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
	var v Cluster
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
					"name": "${var.name}",
					//"count":        "${var.k8s_number}",
					"version":               "1.20.11-aliyun.1",
					"os_type":               "linux",
					"platform":              "CentOS",
					"timeout_mins":          "60",
					"vpc_id":                "${apsarastack_vpc.default.id}",
					"image_id":              "${var.image_id}",
					"master_count":          "3",
					"master_disk_category":  "cloud_ssd",
					"master_disk_size":      "45",
					"master_instance_types": "${var.master_instance_types}",
					"master_vswitch_ids":    "${var.vswitch_ids}",
					//"master_vswitch_ids":    []string{"${apsarastack_vswitch.default.id}，${apsarastack_vswitch.default.id}，${apsarastack_vswitch.default.id}"},

					//"num_of_nodes":         "${var.worker_number}",
					"num_of_nodes":         "1",
					"worker_disk_category": "cloud_ssd",
					"worker_disk_size":     "30",
					"runtime": []map[string]interface{}{
						{
							"name":    "Docker",
							"version": "19.03.15",
						},
					},
					"worker_instance_types": "${var.worker_instance_types}",
					"worker_vswitch_ids":    "${var.vswitch_ids}",
					"security_group_id":     "${apsarastack_security_group.default.id}",
					// "enable_ssh":            "${var.enable_ssh}",
					"password":             "${var.password}",
					"delete_protection":    "false",
					"pod_cidr":             "${var.pod_cidr}",
					"service_cidr":         "${var.service_cidr}",
					"node_cidr_mask":       "${var.node_cidr_mask}",
					"new_nat_gateway":      "true",
					"slb_internet_enabled": "true",
					"proxy_mode":           "ipvs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
// variable "k8s_number" {
//   description = "The number of kubernetes cluster."
//   default     = 1
// }
variable "image_id" {
  default     = "centos_7_9_x64_20G_alibase_20220322.vhd"
}
# leave it to empty then terraform will create several vswitches
variable "vswitch_ids" {
 description = "List of existing vswitch id."
 type        = list(string)
 default     = ["${apsarastack_vswitch.default.id}","${apsarastack_vswitch.default.id}","${apsarastack_vswitch.default.id}"]
}
variable "new_nat_gateway" {
  description = "Whether to create a new nat gateway. In this template, a new nat gateway will create a nat gateway, eip and server snat entries."
  default     = "true"
}
# 3 masters is default settings,so choose three appropriate instance types in the availability zones above.
variable "master_instance_types" {
  description = "The ecs instance types used to launch master nodes."
  default     = ["ecs.s7-k-c1m2.2xlarge","ecs.s7-k-c1m2.2xlarge","ecs.s7-k-c1m2.2xlarge"]
}
variable "worker_instance_types" {
  description = "The ecs instance types used to launch worker nodes."
  default     = ["ecs.s7-k-c1m2.2xlarge"]
}
# options: between 24-28
variable "node_cidr_mask" {
  description = "The node cidr block to specific how many pods can run on single node."
  default     = 24
}
// variable "enable_ssh" {
//   description = "Enable login to the node through SSH."
//   default     = true
// }
variable "password" {
  description = "The password of ECS instance."
  default     = "Alibaba@1688"
}
variable "worker_number" {
  description = "The number of worker nodes in kubernetes cluster."
  default     = 3
}
# k8s_pod_cidr is only for flannel network
variable "pod_cidr" {
  description = "The kubernetes pod cidr block. It cannot be equals to vpc's or vswitch's and cannot be in them."
  default     = "172.20.0.0/16"
}
variable "service_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or pod's and cannot be in them."
  default     = "172.21.0.0/20"
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

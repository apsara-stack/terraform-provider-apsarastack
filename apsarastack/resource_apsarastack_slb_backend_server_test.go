package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApsaraStackSlbBackendServers_vpc(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "apsarastack_slb_backend_server.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerVpcCountConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${apsarastack_slb.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${apsarastack_instance.instance.0.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "1",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

func TestAccApsaraStackSlbBackendServers_multi_vpc(t *testing.T) {

	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "apsarastack_slb_backend_server.default.1"
	ra := resourceAttrInit(resourceId, slb_vpc)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc_multi%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":            "2",
					"load_balancer_id": "${apsarastack_slb.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${apsarastack_instance.instance.0.id}",
							"weight":    "80",
						},
						{
							"server_id": "${apsarastack_instance.instance.1.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "2",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccApsaraStackSlbBackendServers_classic(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "apsarastack_slb_backend_server.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc_multi%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.SlbClassicNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(), //testAccCheckSlbBackendServersDestroy,
		Steps: []resource.TestStep{
			{
				//Config: testAccSlbBackendServersClassic,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${apsarastack_slb.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${apsarastack_instance.instance.0.id}",
							"weight":    "80",
						},
						{
							"server_id": "${apsarastack_instance.instance.1.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "2",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				//Config: testAccSlbBackendServersClassicUpdateServer,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${apsarastack_slb.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${apsarastack_instance.instance.0.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "1",
					}),
				),
			},
			{
				//Config: testAccSlbBackendServersClassic,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${apsarastack_slb.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${apsarastack_instance.instance.0.id}",
							"weight":    "80",
						},
						{
							"server_id": "${apsarastack_instance.instance.1.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "2",
					}),
				),
			},
		},
	})
}

func buildBackendServersMap(count int) []map[string]interface{} {
	var result []map[string]interface{}

	str := `${apsarastack_instance.instance.%d.id}`
	for i := 0; i < count; i++ {
		tmp := make(map[string]interface{}, 2)
		tmp["server_id"] = fmt.Sprintf(str, i)
		tmp["weight"] = "10"
		result = append(result, tmp)
	}
	return result
}

func resourceBackendServerVpcCountConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSlbBackendServersVpc"
}
data "apsarastack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "apsarastack_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}
data "apsarastack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}
resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/16"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  name              = "${var.name}"
}
resource "apsarastack_security_group" "group" {
  name   = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_instance" "instance" {
  image_id                   = "${data.apsarastack_images.default.images.0.id}"
  instance_type              = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  instance_name              = "${var.name}"
  count                      = "2"
  security_groups            = "${apsarastack_security_group.group.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone          = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = "${apsarastack_vswitch.default.id}"
}
resource "apsarastack_slb" "default" {
  name          = "${var.name}"
  vswitch_id    = "${apsarastack_vswitch.default.id}"
}


data "apsarastack_instance_types" "new" {
 	availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
	eni_amount = 2
}
resource "apsarastack_network_interface" "default" {
    count = 1
    name = "${var.name}"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    security_groups = [ "${apsarastack_security_group.group.id}" ]
}
resource "apsarastack_instance" "new" {
  image_id = "${data.apsarastack_images.default.images.0.id}"
  instance_type = "${data.apsarastack_instance_types.new.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "1"
  security_groups = "${apsarastack_security_group.group.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}
resource "apsarastack_network_interface_attachment" "default" {
	count = 1
    instance_id = "${apsarastack_instance.new.0.id}"
    network_interface_id = "${element(apsarastack_network_interface.default.*.id, count.index)}"
}
`)
}

func resourceBackendServerConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSlbBackendServersVpc"
}
data "apsarastack_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "apsarastack_instance_types" "default" {
	cpu_core_count = 1
	memory_size = 2
}
data "apsarastack_images" "default" {
    name_regex = "^ubuntu_18.*64"
	most_recent = true
	owners = "system"
}
resource "apsarastack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
    vpc_id = "${apsarastack_vpc.default.id}"
    cidr_block = "172.16.0.0/16"
    availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
    name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  	name = "${var.name}"
	vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_instance" "instance" {
  	image_id = "${data.apsarastack_images.default.images.0.id}"
  	instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	count = "2"
  	security_groups = "${apsarastack_security_group.default.*.id}"
  	internet_max_bandwidth_out = "10"
  	availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  	system_disk_category = "cloud_efficiency"
  	vswitch_id = "${apsarastack_vswitch.default.id}"
}
resource "apsarastack_slb" "default" {
  	name = "${var.name}"
  	vswitch_id = "${apsarastack_vswitch.default.id}"
}

`)
}

var slb_vpc = map[string]string{
	"backend_servers.#": "2",
}

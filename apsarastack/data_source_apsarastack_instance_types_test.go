package apsarastack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApsaraStackInstanceTypesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackInstanceTypesDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_instance_types.c4g8"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c4g8", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c4g8", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c4g8", "instance_types.0.memory_size", "8"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c4g8", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c4g8", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c4g8", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c4g8", "instance_types.0.gpu.%", "2"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c4g8", "instance_types.0.gpu.amount"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c4g8", "instance_types.0.gpu.category", ""),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c4g8", "instance_types.0.burstable_instance.%", "2"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c4g8", "instance_types.0.burstable_instance.initial_credit"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c4g8", "instance_types.0.burstable_instance.baseline_credit"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c4g8", "instance_types.0.local_storage.%", "3"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c4g8", "instance_types.0.local_storage.capacity"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c4g8", "instance_types.0.local_storage.amount"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c4g8", "instance_types.0.local_storage.category", ""),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c4g8", "ids.#"),
				),
			},
		},
	})
}

func TestAccApsaraStackInstanceTypesDataSource_gpu(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackInstanceTypesDataSourceGpu,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_instance_types.gpu"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.price"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.id"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.cpu_core_count"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.memory_size"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.family", "ecs.gn5"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.gpu.%", "2"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.gpu.amount"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.gpu.category"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.burstable_instance.%", "2"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.burstable_instance.initial_credit"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.burstable_instance.baseline_credit"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.local_storage.%", "3"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.local_storage.capacity"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.local_storage.amount"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.local_storage.category"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "ids.#"),
				),
			},
			{
				Config: testAccCheckApsaraStackInstanceTypesDataSourceGpuK8SMaster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_instance_types.gpu"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.price"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.gpu.%"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.burstable_instance.%"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.local_storage.%"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "ids.#"),
				),
			},
			{
				Config: testAccCheckApsaraStackInstanceTypesDataSourceGpuK8SWorker,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_instance_types.gpu"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.id"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.price"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.cpu_core_count"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.memory_size"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.family", "ecs.gn5"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.gpu.%", "2"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.gpu.amount"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.gpu.category"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.burstable_instance.%", "2"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.burstable_instance.initial_credit"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.burstable_instance.baseline_credit"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.gpu", "instance_types.0.local_storage.%", "3"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.local_storage.capacity"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.local_storage.amount"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "instance_types.0.local_storage.category"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.gpu", "ids.#"),
				),
			},
		},
	})
}

func TestAccApsaraStackInstanceTypesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackInstanceTypesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_instance_types.empty"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.empty", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.empty", "instance_types.0.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.empty", "instance_types.0.price"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.empty", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.empty", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.empty", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.empty", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.empty", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.empty", "instance_types.0.gpu.%"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.empty", "instance_types.0.burstable_instance.%"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.empty", "instance_types.0.local_storage.%"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.empty", "ids.#"),
				),
			},
		},
	})
}

func TestAccApsaraStackInstanceTypesDataSource_k8sSpec(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackInstanceTypesDataSourceK8Sc1g2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_instance_types.c1g2"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c1g2", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.c1g2", "instance_types.0.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.c1g2", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.c1g2", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.c1g2", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.c1g2", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.c1g2", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.c1g2", "instance_types.0.gpu.%"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.c1g2", "instance_types.0.burstable_instance.%"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.c1g2", "instance_types.0.local_storage.%"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c1g2", "ids.#"),
				),
			},
			{
				Config: testAccCheckApsaraStackInstanceTypesDataSourceK8Sc2g4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_instance_types.c2g4"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c2g4", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c2g4", "instance_types.0.cpu_core_count", "2"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c2g4", "instance_types.0.memory_size", "4"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c2g4", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c2g4", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c2g4", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c2g4", "instance_types.0.gpu.%", "2"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c2g4", "instance_types.0.gpu.amount"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c2g4", "instance_types.0.gpu.category", ""),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c2g4", "instance_types.0.burstable_instance.%", "2"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c2g4", "instance_types.0.burstable_instance.initial_credit"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c2g4", "instance_types.0.burstable_instance.baseline_credit"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c2g4", "instance_types.0.local_storage.%", "3"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c2g4", "instance_types.0.local_storage.capacity"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c2g4", "instance_types.0.local_storage.amount"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.c2g4", "instance_types.0.local_storage.category", ""),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.c2g4", "ids.#"),
				),
			},
		},
	})
}
func TestAccApsaraStackInstanceTypesDataSource_k8sFamily(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackInstanceTypesDataSourceK8ST5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_instance_types.t5"),
					resource.TestCheckResourceAttr("data.apsarastack_instance_types.t5", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.t5", "instance_types.0.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.t5", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.t5", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.t5", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.t5", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.t5", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.t5", "instance_types.0.gpu.%"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.t5", "instance_types.0.burstable_instance.%"),
					resource.TestCheckNoResourceAttr("data.apsarastack_instance_types.t5", "instance_types.0.local_storage.%"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instance_types.t5", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackInstanceTypesDataSourceBasicConfig = `
data "apsarastack_instance_types" "c4g8" {
	cpu_core_count = 4
	memory_size = 8
}
`

const testAccCheckApsaraStackInstanceTypesDataSourceGpu = `
provider "apsarastack" {
	region = "cn-hangzhou"
}
data "apsarastack_instance_types" "gpu" {
	sorted_by = "Price"
	instance_type_family = "ecs.gn5"
}
`
const testAccCheckApsaraStackInstanceTypesDataSourceGpuK8SMaster = `
provider "apsarastack" {
	region = "cn-hangzhou"
}
data "apsarastack_instance_types" "gpu" {
	sorted_by = "Price"
	kubernetes_node_role = "Master"
	instance_type_family = "ecs.gn5"
}
`
const testAccCheckApsaraStackInstanceTypesDataSourceGpuK8SWorker = `
provider "apsarastack" {
	region = "cn-hangzhou"
}
data "apsarastack_instance_types" "gpu" {
	sorted_by = "Price"
	kubernetes_node_role = "Worker"
	instance_type_family = "ecs.gn5"
}
`

const testAccCheckApsaraStackInstanceTypesDataSourceEmpty = `
data "apsarastack_instance_types" "empty" {
	instance_type_family = "ecs.fake"
}
`

const testAccCheckApsaraStackInstanceTypesDataSourceK8Sc1g2 = `
data "apsarastack_instance_types" "c1g2" {
	cpu_core_count = 1
	memory_size = 2
	kubernetes_node_role = "Master"
}
`
const testAccCheckApsaraStackInstanceTypesDataSourceK8Sc2g4 = `
data "apsarastack_instance_types" "c2g4" {
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}
`
const testAccCheckApsaraStackInstanceTypesDataSourceK8ST5 = `
data "apsarastack_instance_types" "t5" {
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Master"
	instance_type_family = "ecs.t5"
}
`

---
subcategory: "ECS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ecs_hpc_cluster"
sidebar_current: "docs-apsarastack-resource-ecs-hpc-cluster"
description: |-
  Provides a Apsarastack ECS Hpc Cluster resource.
---

# apsarastack\_ecs\_hpc\_cluster

Provides a ECS Hpc Cluster resource.

For information about ECS Hpc Cluster and how to use it, see [What is Hpc Cluster](https://help.aliyun.com/document_detail/109138.html?spm=5176.21213303.J_6704733920.7.21d953c9oW34ti&scm=20140722.S_help%40%40%E6%96%87%E6%A1%A3%40%40109138._.ID_help%40%40%E6%96%87%E6%A1%A3%40%40109138-RL_CreateHpcCluster-LOC_main-OR_ser-V_2-P0_0).

-> **NOTE:** Available in v1.116.0+.

## Example Usage

Basic Usage

```
resource "apsarastack_ecs_hpc_cluster" "default" {
  name = "tf-testAcccn-qingdao-env17-d01ApsaraStackEcsHpcCluster21597Update"
  description = "Test For Terraform"
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of ECS Hpc Cluster.
* `name` - (Required) The name of ECS Hpc Cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Hpc Cluster.

## Import

ECS Hpc Cluster can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_hpc_cluster.example <id>
```

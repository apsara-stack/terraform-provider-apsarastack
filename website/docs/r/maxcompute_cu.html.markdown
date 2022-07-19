---
subcategory: "MaxCompute"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_maxcompute_cu"
sidebar_current: "docs-apsarastack-resource-maxcompute-cu"
description: |-
  Provides a Apsarastack maxcompute cu resource.
---

# apsarastack\_maxcompute\_cu

The cu is the basic unit of operation in maxcompute. It is similar to the concept of Database or Schema in traditional databases, and sets the boundary for maxcompute multi-user isolation and access control.
->**NOTE:** Available in 1.0.18+.

## Example Usage

Basic Usage

```terraform
resource "apsarastack_maxcompute_cu" "example" {
   cu_name      = "testcu"
  cu_num       = "1"
  cluster_name = "HYBRIDODPSCLUSTER-A-20210520-07B0"
}
```
## Argument Reference

The following arguments are supported:
* `cu_name` - (Required, ForceNew) The name of the cu that you want to create.
* `cluster_name` - (Required, ForceNew) The name of the cluster that you want to create.
* `cu_num` - (Required, ForceNew) The num of the maxcompute cu. 

## Attributes Reference



## Import

MaxCompute project can be imported using the *name* or ID, e.g.

```
$ terraform import apsarastack_maxcompute_cu.example tf_maxcompute_cu
```

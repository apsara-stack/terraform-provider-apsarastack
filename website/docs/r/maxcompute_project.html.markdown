---
subcategory: "MaxCompute"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_maxcompute_project"
sidebar_current: "docs-apsarastack-resource-maxcompute-project"
description: |-
  Provides a Apsarastack maxcompute project resource.
---

# apsarastack\_maxcompute\_project

The project is the basic unit of operation in maxcompute. It is similar to the concept of Database or Schema in traditional databases, and sets the boundary for maxcompute multi-user isolation and access control.
->**NOTE:** Available in 1.77.0+.

## Example Usage

Basic Usage

```terraform
resource "apsarastack_maxcompute_project" "example" {
    cluster        = "HYBRIDODPSCLUSTER-A-20210520-07B0"
	external_table = "false"
	quota_id       = "38"
	disk           = "5"
	name           = "tf_testAccApsaraStack3011"
    pk = "1075451910171540"
}
```
## Argument Reference

The following arguments are supported:
* `name` - (Required, ForceNew) It has been deprecated from provider version 1.110.0 and `project_name` instead.
* `cluster` - (Required, ForceNew) The name of the cluster that you want to create..
* `external_table` - (Optional, ForceNew) Determines whether to automatically split a shard. Default to `false`. 
* `quota_id` - (Required, ForceNew)  `quota_id` - ID of the quota.
* `disk` - (Required, ForceNew)  User-defined instance one core node's storage. space.Unit: GB. Value range:
* `pk` - (Required, ForceNew)  `pk` - ID of the TaskPk.
## Attributes Reference



## Import

MaxCompute project can be imported using the *name* or ID, e.g.

```
$ terraform import apsarastack_maxcompute_project.example tf_maxcompute_project
```

---
subcategory: "ECS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_instance_type_families"
sidebar_current: "docs-apsarastack-datasource-instance-type-families"
description: |-
    Provides a list of ECS Instance Type Families to be used by the apsarastack_instance resource.
---

# apsarastack\_instance\_type\_families

This data source provides the ECS instance type families of ApsaraStack.

-> **NOTE:** Available in 1.54.0+

## Example Usage

```
data "apsarastack_instance_type_families" "default" {
  instance_charge_type = "PrePaid"
}

output "first_instance_type_family_id" {
  value = "${data.apsarastack_instance_type_families.default.instance_type_families.0.id}"
}

output "instance_ids" {
  value = "${data.apsarastack_instance_type_families.default.ids}"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional, ForceNew) The Zone to launch the instance.
* `generation` - (Optional) The generation of the instance type family, Valid values: `ecs-1`, `ecs-2`, `ecs-3` and `ecs-4`. For more information, see [Instance type families](https://www.alibabacloud.com/help/doc-detail/25378.htm). 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance type family IDs.
* `id` - ID of the instance type family.
* `generation` - The generation of the instance type family.
* `zone_ids` - A list of Zone to launch the instance.
 
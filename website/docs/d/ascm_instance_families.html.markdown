---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_instance_families"
sidebar_current: "docs-apsarastack-datasource-ascm-instance-families"
description: |-
    Provides a list of instance families to the user.
---

# apsarastack\_ascm_instance_families

This data source provides the instance families of the current Apsara Stack Cloud user.

## Example Usage

```
data "apsarastack_ascm_instance_families" "default" {
    name_regex = "AutoTest"
    output_file = "instance_families"
    resource_type = "DRDS"
}
output "families" {
    value = data.apsarastack_ascm_instance_families.default.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance family IDs.
* `name_regex` - (Optional) A regex string to filter the resulting instance families by their series_names.
* `resource_type` - (Optional) Filter the results by the specified resource type.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `families` - A list of instance families. Each element contains the following attributes:
    * `id` - ID of the instance families.
    * `series_name` - Series name for instance families.
    * `modifier` - Modifier name.
    * `series_name` - Series name for instance families.
    * `resource_type` - Specified resource type.
    * `is_deleted` - Specify the state in "Y" or "N" form.
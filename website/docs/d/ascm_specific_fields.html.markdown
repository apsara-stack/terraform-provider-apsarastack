---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "ApsaraStack: apsarastack_ascm_specific_fields"
sidebar_current: "docs-apsarastack-datasource-specific-fields"
description: |-
    Provides a list of specific fields to the user.
---

# apsarastack\_specific_fields

This data source provides the specific fields of the current Apsara Stack Cloud user.

## Example Usage

```
data "apsarastack_ascm_specific_fields" "specifields" {
  group_filed ="storageType"
  resource_type ="OSS"
  output_file = "fields"
}
output "specifields" {
  value = data.apsarastack_ascm_specific_fields.specifields.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of specific fields IDs.
* `group_filed` - (Required) Filter the results by specified group filed.
* `resource_type` - (Required) Filter the results by the specified resource type.
* `label` - (Optional) Filter the results by the specified label. Takes Bool Value.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `specific_fields` - A list of specific fields.
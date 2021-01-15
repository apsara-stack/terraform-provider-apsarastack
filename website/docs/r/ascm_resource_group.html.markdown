---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_resource_group"
sidebar_current: "docs-apsarastack-resource-ascm-resource-group"
description: |-
  Provides a Ascm resource group resource.
---

# apsarastack\_ascm_resource_group

Provides a Ascm resource group resource.

## Example Usage

```
resource "apsarastack_ascm_resource_group" "default" {
  name = "apsara_resource_group"
  organization_id = "437"
}
output "org" {
  value = apsarastack_ascm_resource_group.default.*
}
```
## Argument Reference

The following arguments are supported:

* `rg_id` - The ID of the resource group.
* `name` - (Required) The name of the resource group. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
* `organization_id` - (Required) ID of an Organization.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource group.
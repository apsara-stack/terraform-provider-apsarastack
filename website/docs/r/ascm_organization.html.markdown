---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "ApsaraStack: apsarastack_ascm_organization"
sidebar_current: "docs-apsarastack-resource-ascm-organization"
description: |-
  Provides an Ascm organization resource.
---

# apsarastack\_ascm_organization

Provides an Ascm organization resource.

## Example Usage

```
resource "apsarastack_ascm_organization" "default" {
  name = "apsara_Organization"
  parent_id = "19"
}
output "org" {
  value = apsarastack_ascm_organization.default.*
}
```
## Argument Reference

The following arguments are supported:

* `org_id` - (Computed) The ID of the organization.
* `name` - (Required) The name of the organization. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
* `parent_id` - (Required) Parent ID.
* `person_num` - (Optional) A reserved parameter.
* `resource_group_num` - (Optional) A reserved parameter.

## Attributes Reference

The following attributes are exported:

* `id` - Name and ID of the organization. The value is in format `Name/ID`
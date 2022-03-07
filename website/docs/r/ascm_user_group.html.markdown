---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_user_group"
sidebar_current: "docs-apsarastack-resource-ascm-user_group"
description: |-
  Provides a Ascm user group resource.
---

# apsarastack\_ascm_user_group

Provides a Ascm user group resource.

## Example Usage

```
resource "apsarastack_ascm_organization" "default" {
  name = "Dummy_Test_1"
}

resource "apsarastack_ascm_user_group" "default" {
   group_name = "test"
   organization_id = apsarastack_ascm_organization.default.org_id
}

output "org" {
  value = apsarastack_ascm_user_group.default.*
}
```
## Argument Reference

The following arguments are supported:

* `group_name` - (Required) group name. 
* `organization_id` - (Required) User Organization ID.

## Attributes Reference

The following attributes are exported:

* `id` - Login Name of the user group.
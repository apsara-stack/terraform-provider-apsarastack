---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "ApsaraStack: apsarastack_ascm_roles"
sidebar_current: "docs-apsarastack-datasource-ascm-roles"
description: |-
    Provides a list of roles to the user.
---

# apsarastack\_ascm_roles

This data source provides the roles of the current Apsara Stack Cloud user.

## Example Usage

```
resource "apsarastack_ascm_ram_role" "default" {
  role_name = "DELTA1"
  description = "Testing Complete"
  organization_visibility = "global"
}

data "apsarastack_ascm_roles" "default" {
  name_regex = apsarastack_ascm_ram_role.default.role_name
}

output "roles" {
  value = data.apsarastack_ascm_roles.default.*
}


```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of roles IDs.
* `name_regex` - (Optional) A regex string to filter results by role name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `roles` - A list of roles. Each element contains the following attributes:
    * `id` - ID of the role.
    * `name` - role name.
    * `description` - Description about the role.
    * `role_level` - role level.
    * `role_type` - types of role.
    * `ram_role` - ram authorized role.
    * `role_range` - specific range for a role.
    * `user_count` - user count.
     

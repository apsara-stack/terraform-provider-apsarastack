---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_user_groups"
sidebar_current: "docs-apsarastack-datasource-ascm-user-groups"
description: |-
    Provides a list of users to the user.
---

# apsarastack\_ascm_user_groups

This data source provides the users of the current Apsara Stack Cloud user.

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
  value = apsarastack_ascm_user_groups.default.*
}
data "apsarastack_ascm_user_groups" "default" {
 ids = [apsarastack_ascm_user_group.default.user_id]
}
output "groups" {
 value = data.apsarastack_ascm_user_groups.default.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of groups IDs.
     

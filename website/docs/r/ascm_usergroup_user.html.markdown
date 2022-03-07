---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_usergroup_user"
sidebar_current: "docs-apsarastack-resource-ascm-usergroup_user"
description: |-
  Provides a Ascm usergroup_user resource.
---

# apsarastack\_ascm_usergroup_user

Provides a Ascm usergroup_user resource.

## Example Usage

```
resource "apsarastack_ascm_organization" "default" {
 name = "Test_binder"
 parent_id = "1"
}

resource "apsarastack_ascm_user_group" "default" {
 group_name =      "%s"
 organization_id = apsarastack_ascm_organization.default.org_id
}

resource "apsarastack_ascm_user" "default" {
 cellphone_number = "13900000000"
 email = "test@gmail.com"
 display_name = "C2C-DELTA"
 organization_id = apsarastack_ascm_organization.default.org_id
 mobile_nation_code = "91"
 login_name = "User_Role_Test%d"
 login_policy_id = 1
}

resource "apsarastack_ascm_usergroup_user" "default" {
  login_names = ["${apsarastack_ascm_user.default.login_name}"]
  user_group_id = apsarastack_ascm_user_group.default.user_group_id
}

output "org" {
  value = apsarastack_ascm_usergroup_user.default.*
}
```
## Argument Reference

The following arguments are supported:

* `user_group_id` - (Required) group name. 
* `login_names` - (Required) List of user login name.

## Attributes Reference

The following attributes are exported:

* `id` - Login Name of the usergroup_user.
---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "ApsaraStack: apsarastack_ascm_users"
sidebar_current: "docs-apsarastack-datasource-ascm-users"
description: |-
    Provides a list of users to the user.
---

# apsarastack\_ascm_users

This data source provides the users of the current Apsara Stack Cloud user.

## Example Usage

```
data "apsarastack_ascm_users" "users" {
 ids = [apsarastack_ascm_user.user.id]
}
output "users" {
 value = data.apsarastack_ascm_users.users.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of users IDs.
* `name_regex` - (Optional) A regex string to filter results by user login name.
* `cell_phone_number` - (Optional) Filter the results by the Cellphone Number of a user.
* `display_name` - (Optional) Filter the results by the Display name of a user.
* `email` - (Optional) Filter the results by the Email ID of a user.
* `mobile_nation_code` - (Optional) Filter the results by the Mobile Nation Code of a user, where user belongs to.
* `organization_id` - (Optional) Filter the results by the specified user Organization ID.
* `login_policy_id` - (Optional) Filter the results by the specified user login policy ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `users` - A list of users. Each element contains the following attributes:
    * `id` - ID of the user.
    * `name` - User login name.
    * `cell_phone_number` - Cellphone Number of a user.
    * `display_name` - Display name of a user.
    * `email` - Email ID of a user.
    * `mobile_nation_code` - Mobile Nation Code of a user, where user belongs to.
    * `organization_id` - User Organization ID.
    * `login_policy_id` - User login policy ID.
     

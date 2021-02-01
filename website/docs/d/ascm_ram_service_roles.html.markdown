---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_ram_service_roles"
sidebar_current: "docs-apsarastack-datasource-ascm-ram-service-roles"
description: |-
    Provides a list of ram roles to the user.
---

# apsarastack\_ascm_ram_service_roles

This data source provides the ram roles of the current Apsara Stack Cloud user.

## Example Usage

```
data "apsarastack_ascm_ram_service_roles" "role" {
  product = "ECS"
}
output "role" {
  value = data.apsarastack_ascm_ram_service_roles.role.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ram roles IDs.
* `product` - (Optional) A regex string to filter results by their product. valid values - "ECS".
* `description` - (Optional) Description about the ram role.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `roles` - A list of roles. Each element contains the following attributes:
    * `id` - ID of the role.
    * `name` - role name.
    * `description` - Description about the role.
    * `role_type` - types of role.
    * `product` - types of role.
    * `organization_name` - Name of an Organization.
    * `aliyun_user_id` - Aliyun User Id.
     

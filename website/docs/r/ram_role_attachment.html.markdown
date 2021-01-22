---
subcategory: "RAM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ram_role_attachment"
sidebar_current: "docs-apsarastack-resource-ram-role-attachment"
description: |-
  Provides a RAM role attachment resource to bind role for several ECS instances.
---

# apsarastack\_ram\_role\_attachment

Provides a RAM role attachment resource to bind role for several ECS instances.

## Example Usage

```
resource "apsarastack_ram_role_attachment" "default" {
  role_name    = apsarastack_ascm_ram_role.ramrole.role_name
  instance_ids = ["i-2sdfdasd3423g48dhsa"]
}

output "ramrole" {
  value = apsarastack_ram_role_attachment.default.role_name
}

```

## Argument Reference

The following arguments are supported:

* `role_name` - (Required, ForceNew) The name of role used to bind. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-", "_", and must not begin with a hyphen.
* `instance_ids` - (Required, ForceNew) The list of ECS instance's IDs.

## Attributes Reference

The following attributes are exported:

* `role_name` - The name of the role.
* `instance_ids` The list of ECS instance's IDs.

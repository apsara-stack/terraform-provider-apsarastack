---
subcategory: "RocketMQ"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ons_group"
sidebar_current: "docs-apsarastack-resource-ons-group"
description: |-
  Provides a apsarastack ONS Group resource.
---

# apsarastack\_ons\_group

Provides an ONS group resource.


## Example Usage

Basic Usage

```
variable "name" {
  default = "onsInstanceName"
}

variable "group_id" {
  default = "GID-onsGroupDatasourceName"
}

resource "apsarastack_ons_instance" "default" {
  name = "${var.name}"
  remark = "default_ons_instance_remark"
}

resource "apsarastack_ons_group" "default" {
  group_id = "${var.group_id}"
  instance_id = "${apsarastack_ons_instance.default.id}"
  remark = "dafault_ons_group_remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ONS Instance that owns the groups.
* `group_id` - (Required) Name of the group. Two groups on a single instance cannot have the same name. A `group_id` starts with "GID_" or "GID-", and contains letters, numbers, hyphens (-), and underscores (_).
* `remark` - (Optional) This attribute is a concise description of group. The length cannot exceed 256.
* `read_enable` - (Optional) This attribute is used to set the message reading enabled or disabled. It can only be set after the group is used by the client.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<group_id>`.



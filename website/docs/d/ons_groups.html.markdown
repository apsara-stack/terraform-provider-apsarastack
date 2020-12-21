---
subcategory: "RocketMQ"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ons_groups"
sidebar_current: "docs-apsarastack-datasource-ons-groups"
description: |-
    Provides a list of ons groups available to the user.
---

# apsarastack\_ons\_groups

This data source provides a list of ONS Groups in an Apsara Stack Cloud account according to the specified filters.


## Example Usage

```
variable "name" {
  default = "onsInstanceName"
}

variable "group_id" {
  default = "GID-onsGroupDatasourceName"
}

resource "apsarastack_ons_instance" "default" {
  tps_receive_max = "500"
  tps_send_max = "500"
  topic_capacity = "50"
  cluster = "cluster1"
  independent_naming = "true"
  name = "Ons_Apsara_instance"
  remark = "Ons Instance"
}

resource "apsarastack_ons_group" "default" {
  group_id = var.group_id
  instance_id = "${apsarastack_ons_instance.default.id}"
  remark = "dafault_ons_group_remark"
}

data "apsarastack_ons_groups" "default" {
  instance_id = apsarastack_ons_group.default.instance_id

}
output "onsgroups" {
  value = data.apsarastack_ons_groups.default.*
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ONS Instance that owns the groups.
* `group_id_regex` - (Optional) A regex string to filter results by the group name. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of group names.
* `groups` - A list of groups. Each element contains the following attributes:
  * `id` - The name of the group.
  * `owner` - The ID of the group owner, which is the Apsara Stack Cloud UID.
  * `independent_naming` - Indicates whether namespaces are available.
  * `remark` - Remark of the group.

---
subcategory: "ECS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_security_groups"
sidebar_current: "docs-apsarastack-datasource-security-groups"
description: |-
    Provides a list of Security Groups available to the user.
---

# apsarastack\_security\_groups

This data source provides a list of Security Groups in an ApsaraStack account according to the specified filters.

## Example Usage

```
# Filter security groups and print the results into a file
data "apsarastack_security_groups" "sec_groups_ds" {
  name_regex  = "^web-"
}

# In conjunction with a VPC
resource "apsarastack_vpc" "primary_vpc_ds" {
  # ...
}

data "apsarastack_security_groups" "primary_sec_groups_ds" {
  vpc_id = "${apsarastack_vpc.primary_vpc_ds.id}"
}

output "first_group_id" {
  value = "${data.apsarastack_security_groups.primary_sec_groups_ds.groups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, Available 1.52.0+) A list of Security Group IDs.
* `name_regex` - (Optional) A regex string to filter the resulting security groups by their names.
* `vpc_id` - (Optional) Used to retrieve security groups that belong to the specified VPC ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `tags` - (Optional) A map of tags assigned to the ECS instances. It must be in the format:
  ```
  data "apsarastack_security_groups" "taggedSecurityGroups" {
    tags = {
      tagKey1 = "tagValue1",
      tagKey2 = "tagValue2"
    }
  }
  ```

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Security Group IDs.
* `names` - A list of Security Group names.
* `groups` - A list of Security Groups. Each element contains the following attributes:
  * `id` - The ID of the security group.
  * `name` - The name of the security group.
  * `description` - The description of the security group.
  * `vpc_id` - The ID of the VPC that owns the security group.
  * `creation_time` - Creation time of the security group.
  * `tags` - A map of tags assigned to the ECS instance.

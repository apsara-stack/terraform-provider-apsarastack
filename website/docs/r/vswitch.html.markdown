---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_vswitch"
sidebar_current: "docs-apsarastack-resource-vswitch"
description: |-
  Provides a Apsarastack VPC switch resource.
---

# apsarastack\_vswitch

Provides a VPC switch resource.

## Example Usage

Basic Usage

```
resource "apsarastack_vpc" "vpc" {
  name       = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "vsw" {
  vpc_id            = "${apsarastack_vpc.vpc.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "cn-beijing-b"
}
```

## Module Support

You can use to the existing [vpc module](https://registry.terraform.io/modules/alibaba/vpc/apsarastack) 
to create a VPC and several VSwitches one-click.

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The AZ for the switch.
* `vpc_id` - (Required, ForceNew) The VPC ID.
* `cidr_block` - (Required, ForceNew) The CIDR block for the switch.
* `name` - (Optional) The name of the switch. Defaults to null.
* `description` - (Optional) The switch description. Defaults to null.
<!--* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.-->

### Timeouts

-> **NOTE:** Available in 1.79.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the vswitch (until it reaches the initial `Available` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the vswitch. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the switch.
* `availability_zone` The AZ for the switch.
* `cidr_block` - The CIDR block for the switch.
* `vpc_id` - The VPC ID.
* `name` - The name of the switch.
* `description` - The description of the switch.

## Import

Vswitch can be imported using the id, e.g.

```
$ terraform import apsarastack_vswitch.example vsw-abc123456
```

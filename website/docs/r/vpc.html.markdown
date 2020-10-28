---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_vpc"
sidebar_current: "docs-apsarastack-resource-vpc"
description: |-
  Provides a Apsarastack VPC resource.
---

# apsarastack\_vpc

Provides a VPC resource.

-> **NOTE:** Terraform will auto build a router and a route table while it uses `apsarastack_vpc` to build a vpc resource.

## Example Usage

Basic Usage

```
resource "apsarastack_vpc" "vpc" {
  name       = "tf_test_foo"
  cidr_block = "${var.cidr_block}"
}
```


## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, ForceNew) The CIDR block for the VPC.
* `name` - (Optional) The name of the VPC. Defaults to null.
* `description` - (Optional) The VPC description. Defaults to null.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the vpc (until it reaches the initial `Available` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the vpc. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPC.
* `cidr_block` - The CIDR block for the VPC.
* `name` - The name of the VPC.
* `description` - The description of the VPC.
* `router_id` - The ID of the router created by default on VPC creation.

---
subcategory: "Server Load Balancer (SLB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_slb"
sidebar_current: "docs-apsarastack-resource-slb"
description: |-
  Provides an Application Load Balancer resource.
---

# apsarastack\_slb

Provides an Application Load Balancer resource.

-> **NOTE:** At present, to avoid some unnecessary regulation confusion, SLB can not support apsarastack international account to create "paybybandwidth" instance.

-> **NOTE:** The supported specifications vary by region. Currently not all regions support guaranteed-performance instances.
For more details about guaranteed-performance instance, see [Guaranteed-performance instances](https://www.alibabacloud.com/help/doc-detail/27657.htm).

## Example Usage

```
variable "name" {
  default = "terraformtestslbconfig"
}
data "apsarastack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "apsarastack_slb" "default" {
  name          = "${var.name}"
  vswitch_id    = "${apsarastack_vswitch.default.id}"
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
    tag_i = 9
    tag_j = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the SLB. This name must be unique within your apsarastack account, can have a maximum of 80 characters,
must contain only alphanumeric characters or hyphens, such as "-","/",".","_", and must not begin or end with a hyphen. If not specified,
Terraform will autogenerate a name beginning with `tf-lb`.
* `address_type` - (Optional, ForceNew, Available in 1.55.3+) The network type of the SLB instance. Valid values: ["internet", "intranet"]. If load balancer launched in VPC, this value must be "intranet".
    - internet: After an Internet SLB instance is created, the system allocates a public IP address so that the instance can forward requests from the Internet.
    - intranet: After an intranet SLB instance is created, the system allocates an intranet IP address so that the instance can only forward intranet requests.
* `vswitch_id` - (Required for a VPC SLB, Forces New Resource) The VSwitch ID to launch in. If `address_type` is internet, it will be ignore.
* `tags` - (Optional) A mapping of tags to assign to the resource. The `tags` can have a maximum of 10 tag for every load balancer instance.
* `instance_charge_type` - (Optional, Available in v1.34.0+) The billing method of the load balancer. Valid values are "PrePaid" and "PostPaid". Default to "PostPaid".
* `period` - (Optional, Available in v1.34.0+) The duration that you will buy the resource, in month. It is valid when `instance_charge_type` is `PrePaid`. Default to 1. Valid values: [1-9, 12, 24, 36].

-> **NOTE:** A "Shared-Performance" instance can be changed to "Performance-guaranteed", but the change is irreversible.

-> **NOTE:** To change a "Shared-Performance" instance to a "Performance-guaranteed" instance, the SLB will have a short probability of business interruption (10 seconds-30 seconds). Advise to change it during the business downturn, or migrate business to other SLB Instances by using GSLB before changing.

-> **NOTE:** Currently, the apsarastack cloud international account does not support creating a PrePaid SLB instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the load balancer.
* `address` - The IP address of the load balancer.

## Import

Load balancer can be imported using the id, e.g.

```
$ terraform import apsarastack_slb.example lb-abc123456
```

---
subcategory: "ECS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_network_interface"
sidebar_current: "docs-apsarastack-resource-network-interface"
description: |-
  Provides an ECS Elastic Network Interface resource.
---

# apsarastack\_network\_interface

Provides an ECS Elastic Network Interface resource.

For information about Elastic Network Interface and how to use it, see [Elastic Network Interface](https://apsarastackdocument.oss-cn-hangzhou.aliyuncs.com/01_ApsaraStackEnterprise/V3.11.0-intl-en/Alibaba%20Cloud%20Apsara%20Stack%20Enterprise%202001%2C%20Internal_%20V3.11.0%20Developer%20Guide%20-%20Cloud%20Essentials%20and%20Security%2020200513.pdf?spm=a3c0i.214467.3807842930.7.61e76bdb1JWVyX&file=Alibaba%20Cloud%20Apsara%20Stack%20Enterprise%202001%2C%20Internal_%20V3.11.0%20Developer%20Guide%20-%20Cloud%20Essentials%20and%20Security%2020200513.pdf).

-> **NOTE** Only one of private_ips or private_ips_count can be specified when assign private IPs. 

## Example Usage

```
resource "apsarastack_security_group" "secgroup" {
  name        = "SurajG_security"
  description = "Hello Security Group"
  vpc_id      = apsarastack_vpc.vpc.id
}
resource "apsarastack_vpc" "vpc" {
  name       = "surajG_vpc"
  cidr_block = "10.0.0.0/16"
}

resource "apsarastack_vswitch" "vsw" {
  name       = "surajG_vsw"
  vpc_id            = apsarastack_vpc.vpc.id
  cidr_block        = apsarastack_vpc.vpc.cidr_block
  availability_zone = "cn-beijing-b"
}
resource "apsarastack_instance" "apsarainstance" {
  image_id              = "gj2j1g3-45h3nnc-454hj5g"
  instance_type        = "ecs.n4.large"
  system_disk_category = "cloud_efficiency"
  security_groups      = [apsarastack_security_group.secgroup.id]
  instance_name        = "apsarainstance"
  vswitch_id           = apsarastack_vswitch.vsw.id
}

resource "apsarastack_network_interface" "NetInterface" {
  name              = "ENI"
  vswitch_id        = apsarastack_vswitch.vsw.id
  security_groups   = apsarastack_security_group.secgroup.id
  private_ips_count = 1
  description = "Network Interface"
}
```

## Argument Reference

The following arguments are supported:

* `vswitch_id` - (Required, ForceNew) The VSwitch to create the ENI in.
* `security_groups` - (Required) A list of security group ids to associate with.
* `private_ip` - (Optional, ForceNew) The primary private IP of the ENI.
* `name` - (Optional) Name of the ENI. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-", ".", "_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
* `description` - (Optional) Description of the ENI. This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `private_ips`  - (Optional) List of secondary private IPs to assign to the ENI. Don't use both private_ips and private_ips_count in the same ENI resource block.
* `private_ips_count` - (Optional) Number of secondary private IPs to assign to the ENI. Don't use both private_ips and private_ips_count in the same ENI resource block.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `resource_group_id` - (ForceNew, ForceNew, Available in 1.57.0+) The Id of resource group which the network interface belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ENI ID.
* `mac` - (Available in 1.54.0+) The MAC address of an ENI.

## Import

ENI can be imported using the id, e.g.

```
$ terraform import apsarastack_network_interface.eni eni-abc1234567890000
```

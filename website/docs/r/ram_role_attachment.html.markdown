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
data "apsarastack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "apsarastack_instance_types" "default" {
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}

data "apsarastack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}


resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "apsarastack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${apsarastack_security_group.default.id}"
  cidr_ip           = "172.16.0.0/24"
}
variable "name" {
  default = "ecsInstanceVPCExample"
}
resource "apsarastack_instance" "foo" {
  vswitch_id = "${apsarastack_vswitch.default.id}"
  image_id   = "${data.apsarastack_images.default.images.0.id}"

  instance_type        = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  system_disk_category = "cloud_efficiency"

  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups            = ["${apsarastack_security_group.default.id}"]
  instance_name              = "${var.name}"
}
resource "apsarastack_ram_role" "role" {
  name     = "testrole"
  document = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "ecs.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
  description = "this is a test"
  force = true
}

resource "apsarastack_ram_role_attachment" "attach" {
  role_name = "${apsarastack_ram_role.role.name}"
  instance_ids = ["${apsarastack_instance.foo.*.id}"]
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
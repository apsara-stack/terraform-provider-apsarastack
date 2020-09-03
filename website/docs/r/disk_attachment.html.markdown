---
subcategory: "ECS"
layout: "apsarastack"
page_title: "ApsaraStack: apsarastack_disk_attachment"
sidebar_current: "docs-apsarastack-resource-disk-attachment"
description: |-
  Provides a ECS Disk Attachment resource.
---

# apsaratack\_disk\_attachment

Provides an apsarastack ECS Disk Attachment as a resource, to attach and detach disks from ECS Instances.

## Example Usage

Basic usage

```
# Create a new ECS disk-attachment and use it attach one disk to a new instance.

resource "apsarastack_security_group" "ecs_sg" {
  name        = "terraform-test-group"
  description = "New security group"
}

resource "apsarastack_disk" "ecs_disk" {
  availability_zone = "cn-beijing-a"
  size              = "50"

  tags = {
    Name = "TerraformTest-disk"
  }
}

resource "apsarastack_instance" "ecs_instance" {
  image_id              = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  instance_type         = "ecs.n4.small"
  availability_zone     = "cn-beijing-a"
  security_groups       = ["${apsarastack_security_group.ecs_sg.id}"]
  instance_name         = "Hello"
  instance_network_type = "classic"
  internet_charge_type  = "PayByBandwidth"

  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "apsarastack_disk_attachment" "ecs_disk_att" {
  disk_id     = "${apsarastack_disk.ecs_disk.id}"
  instance_id = "${apsarastack_instance.ecs_instance.id}"
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, Forces new resource) ID of the Instance to attach to.
* `disk_id` - (Required, Forces new resource) ID of the Disk to be attached.
* `device_name` - (Optional, Forces new resource) The device name which is used when attaching disk, it will be allocated automatically by system according to default order from /dev/xvdb to /dev/xvdz.
                                                          
## Attributes Reference

The following attributes are exported:

* `instance_id` - ID of the Instance.
* `disk_id` - ID of the Disk.
* `device_name` - The device name exposed to the instance.

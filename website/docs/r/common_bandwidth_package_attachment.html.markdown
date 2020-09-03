---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_common_bandwidth_package_attachment"
sidebar_current: "docs-apsarastack-resource-common-bandwidth-package-attachment"
description: |-
  Provides an apsarastack Common  Attachment resource.
---

# apsarastack\_common\_bandwidth\_package\_attachment

Provides an apsarastack Common Bandwidth Package Attachment resource for associating Common Bandwidth Package to EIP Instance.

-> **NOTE:** Terraform will auto build common bandwidth package attachment while it uses `apsarastack_common_bandwidth_package_attachment` to build a common bandwidth package attachment resource.

For information about common bandwidth package and how to use it, see [What is Common Bandwidth Package](https://www.alibabacloud.com/help/product/55092.htm).

## Example Usage

Basic Usage

```
resource "apsarastack_common_bandwidth_package" "foo" {
  bandwidth   = "2"
  name        = "test_common_bandwidth_package"
  description = "test_common_bandwidth_package"
}

resource "apsarastack_eip" "foo" {
  bandwidth            = "2"
  internet_charge_type = "PayByBandwidth"
}

resource "apsarastack_common_bandwidth_package_attachment" "foo" {
  bandwidth_package_id = "${apsarastack_common_bandwidth_package.foo.id}"
  instance_id          = "${apsarastack_eip.foo.id}"
}

```
## Argument Reference

The following arguments are supported:

* `bandwidth_package_id` - (Required, ForceNew) The bandwidth_package_id of the common bandwidth package attachment, the field can't be changed.
* `instance_id` - (Required, ForceNew) The instance_id of the common bandwidth package attachment, the field can't be changed.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the common bandwidth package attachment id and formates as `<bandwidth_package_id>:<instance_id>`.

## Import

The common bandwidth package attachemnt can be imported using the id, e.g.

```
$ terraform import apsarastack_common_bandwidth_package_attachment.foo cbwp-abc123456:eip-abc123456
```

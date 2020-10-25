---
subcategory: "RocketMQ"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ons_instance"
sidebar_current: "docs-apsarastack-resource-ons-instance"
description: |-
  Provides a apsarastack ONS Instance resource.
---

# apsarastack\_ons\_instance

Provides an ONS instance resource.

## Example Usage

Basic Usage

```
resource "apsarastack_ons_instance" "example" {
  name   = "tf-example-ons-instance"
  remark = "tf-example-ons-instance-remark"
}
```

## Argument Reference

The following arguments are supported:


* `name` - (Required)Two instances on a single account in the same region cannot have the same name. The length must be 3 to 64 characters. Chinese characters, English letters digits and hyphen are allowed.
* `remark` - (Optional)This attribute is a concise description of instance. The length cannot exceed 128.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above.
* `instance_type` - The edition of instance. 1 represents the postPaid edition, and 2 represents the platinum edition.
* `instance_status` - The status of instance. 1 represents the platinum edition instance is in deployment. 2 represents the postpaid edition instance are overdue. 5 represents the postpaid or platinum edition instance is in service. 7 represents the platinum version instance is in upgrade and the service is available.
* `release_time` - Platinum edition instance expiration time.



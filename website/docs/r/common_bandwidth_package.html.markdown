---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_common_bandwidth_package"
sidebar_current: "docs-apsarastack-resource-common-bandwidth-package"
description: |-
  Provides a Apsarastack Common Bandwidth Package resource.
---

# apsarastack\_common_bandwidth_package

Provides a common bandwidth package resource.

-> **NOTE:** Terraform will auto build common bandwidth package instance while it uses `apsarastack_common_bandwidth_package` to build a common bandwidth package resource.

## Example Usage

Basic Usage

```
resource "apsarastack_common_bandwidth_package" "foo" {
  bandwidth            = "200"
  internet_charge_type = "PayByBandwidth"
  name                 = "test-common-bandwidth-package"
  description          = "test-common-bandwidth-package"
}
```
## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) The bandwidth of the common bandwidth package, in Mbps.
* `internet_charge_type` - (Optional, ForceNew) The billing method of the common bandwidth package. Valid values are "PayByBandwidth" and "PayBy95" and "PayByTraffic". "PayBy95" is pay by classic 95th percentile pricing. International Account doesn't supports "PayByBandwidth" and "PayBy95". Default to "PayByTraffic".
* `ratio` - (Optional, ForceNew Available in 1.55.3+) Ratio of the common bandwidth package. It is valid when `internet_charge_type` is `PayBy95`. Default to 100. Valid values: [10-100].
* `name` - (Optional) The name of the common bandwidth package.
* `description` - (Optional) The description of the common bandwidth package instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the common bandwidth package instance id.



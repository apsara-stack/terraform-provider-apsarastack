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
  name                 = "test-common-bandwidth-package"
  description          = "test-common-bandwidth-package"
}
```
## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) The bandwidth of the common bandwidth package, in Mbps.
* `name` - (Optional) The name of the common bandwidth package.
* `description` - (Optional) The description of the common bandwidth package instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the common bandwidth package instance id.



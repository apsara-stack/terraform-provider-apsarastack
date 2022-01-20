---
subcategory: "RAM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ram_role"
sidebar_current: "docs-apsarastack-resource-ram-role"
description: |-
Provides a Ram role resource.
---

# apsarastack\_ram_role

Provides a Ram role resource.

## Example Usage

```

resource "apsarastack_ram_role" "ram_role" {
  name="testfoorole-001"
  role_policy_document="{\"Version\":\"1\",\"Statement\":[{\"Action\":\"sts:AssumeRole\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":[\"log.aliyuncs.com\"]}}]}"
  description="testing role"
}

```
## Argument Reference

The following arguments are supported:

* `name` - The name of the role.
* `role_policy_document` - The content of the role policy document.
* `description` - The description of the role.

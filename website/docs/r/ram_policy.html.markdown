---
subcategory: "RAM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ram_policy"
sidebar_current: "docs-apsarastack-resource-ram-policy"
description: |-
Provides a Ram policy resource.
---

# apsarastack\_ram_policy

Provides a Ram policy resource.

## Example Usage

```

resource "apsarastack_ram_policy" "ram_policy" {
  name="testfoopolicy-001"
  policy_document="{\"Version\":\"1\",\"Statement\":[{\"Condition\":{},\"Action\":[\"oss:PutObject\"],\"Resource\":[\"acs:oss:*:*:alibaba-cloud/log-collector\",\"acs:oss:*:*:alibaba-cloud/log-collector/*\"],\"Effect\":\"Allow\"}]}"
  description="testing policy"
}

```
## Argument Reference

The following arguments are supported:

* `name` - The name of the policy.
* `policy_document` - The content of the policy.
* `description` - The description of the policy.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `policy_type` - Policy type.
* `ram_id` - RAM ID